package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"payment/src/application/commands"
	"payment/src/application/events"
	"payment/src/controllers"
	payment_nats "payment/src/nats"
	"payment/src/repositories"
	"payment/src/routers"
	"payment/src/security"
	"payment/src/tasks"
	payment_tasks "payment/src/tasks/interfaces"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/helpers"
	"github.com/oceano-dev/microservices-go-common/httputil"

	common_grpc_client "github.com/oceano-dev/microservices-go-common/grpc/email/client"
	common_nats "github.com/oceano-dev/microservices-go-common/nats"
	common_repositories "github.com/oceano-dev/microservices-go-common/repositories"
	common_security "github.com/oceano-dev/microservices-go-common/security"
	common_services "github.com/oceano-dev/microservices-go-common/services"
	common_tasks "github.com/oceano-dev/microservices-go-common/tasks"
	common_validator "github.com/oceano-dev/microservices-go-common/validators"

	provider "github.com/oceano-dev/microservices-go-common/trace/otel/jaeger"

	"go.mongodb.org/mongo-driver/mongo"

	consul "github.com/hashicorp/consul/api"
	common_consul "github.com/oceano-dev/microservices-go-common/consul"
)

type Main struct {
	config                 *config.Config
	client                 *mongo.Client
	natsConn               *nats.Conn
	managerSecurityRSAKeys security.ManagerSecurityRSAKeys
	managerCertificates    common_security.ManagerCertificates
	adminMongoDbService    *common_services.AdminMongoDbService
	paymentTasks           payment_tasks.VerifyPaymentTask
	httpServer             httputil.HttpServer
	consulClient           *consul.Client
	serviceID              string
}

func NewMain(
	config *config.Config,
	client *mongo.Client,
	natsConn *nats.Conn,
	managerSecurityRSAKeys security.ManagerSecurityRSAKeys,
	managerCertificates common_security.ManagerCertificates,
	adminMongoDbService *common_services.AdminMongoDbService,
	paymentTasks payment_tasks.VerifyPaymentTask,
	httpServer httputil.HttpServer,
	consulClient *consul.Client,
	serviceID string,
) *Main {
	return &Main{
		config:                 config,
		client:                 client,
		natsConn:               natsConn,
		managerSecurityRSAKeys: managerSecurityRSAKeys,
		managerCertificates:    managerCertificates,
		adminMongoDbService:    adminMongoDbService,
		paymentTasks:           paymentTasks,
		httpServer:             httpServer,
		consulClient:           consulClient,
		serviceID:              serviceID,
	}
}

var production *bool
var disableTrace *bool

func main() {
	production = flag.Bool("prod", false, "use -prod=true to run in production mode")
	disableTrace = flag.Bool("disable-trace", false, "use disable-trace=true if you want to disable tracing completly")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	app, err := startup(ctx)
	if err != nil {
		panic(err)
	}

	err = app.client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer app.client.Disconnect(ctx)

	err = app.client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	defer app.natsConn.Close()

	providerTracer, err := provider.NewProvider(provider.ProviderConfig{
		JaegerEndpoint: app.config.Jaeger.JaegerEndpoint,
		ServiceName:    app.config.Jaeger.ServiceName,
		ServiceVersion: app.config.Jaeger.ServiceVersion,
		Production:     *production,
		Disabled:       *disableTrace,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer providerTracer.Close(ctx)
	log.Println("Connected to Jaegger")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	userMongoExporter, err := app.adminMongoDbService.VerifyMongoDBExporterUser()
	if err != nil {
		log.Fatal(err)
	}

	if !userMongoExporter {
		log.Fatal("MongoDB Exporter user not found!")
	}

	app.managerSecurityRSAKeys.GetAllRSAPrivateKeys()
	app.paymentTasks.Run()
	app.httpServer.RunTLSServer()

	<-done
	err = app.consulClient.Agent().ServiceDeregister(app.serviceID)
	if err != nil {
		log.Printf("consul deregister error: %s", err)
	}

	log.Print("Server Stopped")
	os.Exit(0)
}

func startup(ctx context.Context) (*Main, error) {
	config := config.LoadConfig(*production, "./config/")
	helpers.CreateFolder(config.Folders)
	common_validator.NewValidator("en")

	consulClient, serviceID, err := common_consul.NewConsulClient(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	checkServiceName := common_tasks.NewCheckServiceNameTask()

	certificateServiceNameDone := make(chan bool)
	go checkServiceName.ReloadServiceName(
		ctx,
		config,
		consulClient,
		config.Certificates.ServiceName,
		common_consul.CertificatesAndSecurityKeys,
		certificateServiceNameDone)
	<-certificateServiceNameDone

	emailsServiceNameDone := make(chan bool)
	go checkServiceName.ReloadServiceName(
		ctx,
		config,
		consulClient,
		config.EmailService.ServiceName,
		common_consul.EmailService,
		emailsServiceNameDone)
	<-emailsServiceNameDone

	metricService, err := common_services.NewMetricsService(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	client, err := repositories.NewMongoClient(config)
	if err != nil {
		return nil, err
	}

	certificatesService := common_services.NewCertificatesService(config)
	managerCertificates := common_security.NewManagerCertificates(config, certificatesService)
	emailService := common_grpc_client.NewEmailServiceClientGrpc(config, certificatesService)

	checkCertificates := common_tasks.NewCheckCertificatesTask(config, managerCertificates, emailService)
	certsDone := make(chan bool)
	go checkCertificates.Start(ctx, certsDone)
	<-certsDone

	nc, err := common_nats.NewNats(config, certificatesService)
	if err != nil {
		log.Fatalf("Nats connect error: %+v", err)
	}
	log.Printf("Nats Connected Status: %+v	", nc.Status().String())

	subjects := common_nats.GetPaymentSubjects()
	js, err := common_nats.NewJetStream(nc, "payment", subjects)
	if err != nil {
		log.Fatalf("Nats JetStream create error: %+v", err)
	}

	natsPublisher := common_nats.NewPublisher(js)

	database := repositories.NewMongoDatabase(config, client)
	adminMongoDbRepository := common_repositories.NewAdminMongoDbRepository(database)
	adminMongoDbService := common_services.NewAdminMongoDbService(config, adminMongoDbRepository)

	paymentRepository := repositories.NewPaymentRepository(database)
	managerSecurityRSAKeys := security.NewManagerSecurityRSAKeys(config)

	common_SecurityRSAKeysService := common_services.NewSecurityRSAKeysService(config, certificatesService)
	common_ManagerSecurityRSAKeys := common_security.NewManagerSecurityRSAKeys(config, common_SecurityRSAKeysService)

	paymentTasks := tasks.NewVerifyPaymentTask(managerSecurityRSAKeys, common_ManagerSecurityRSAKeys, paymentRepository, emailService, natsPublisher)

	paymentEventHandler := events.NewPaymentEventHandler(paymentTasks, natsPublisher)
	paymentCommandHandler := commands.NewPaymentCommandHandler(paymentRepository, paymentEventHandler)

	listens := payment_nats.NewListen(
		config,
		js,
		paymentCommandHandler,
		emailService)

	listens.Listen()

	paymentController := controllers.NewPaymentController(managerSecurityRSAKeys)
	router := routers.NewRouter(config, metricService, paymentController)
	httpserver := httputil.NewHttpServer(config, router.RouterSetup(), certificatesService)
	app := NewMain(
		config,
		client,
		nc,
		managerSecurityRSAKeys,
		managerCertificates,
		adminMongoDbService,
		paymentTasks,
		httpserver,
		consulClient,
		serviceID,
	)

	return app, nil
}
