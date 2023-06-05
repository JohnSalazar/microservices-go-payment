package routers

import (
	"fmt"
	"payment/src/controllers"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/middlewares"
	common_service "github.com/oceano-dev/microservices-go-common/services"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Router struct {
	config            *config.Config
	serviceMetrics    common_service.Metrics
	paymentController *controllers.PaymentController
}

func NewRouter(
	config *config.Config,
	serviceMetrics common_service.Metrics,
	paymentController *controllers.PaymentController,

) *Router {
	return &Router{
		config:            config,
		serviceMetrics:    serviceMetrics,
		paymentController: paymentController,
	}
}

func (r *Router) RouterSetup() *gin.Engine {
	router := r.initRoute()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CORS())
	router.Use(location.Default())
	router.Use(otelgin.Middleware(r.config.Jaeger.ServiceName))
	router.Use(middlewares.Metrics(r.serviceMetrics))

	router.GET("/healthy", middlewares.Healthy())
	router.GET("/metrics", middlewares.MetricsHandler())

	v1 := router.Group(fmt.Sprintf("/api/%s", r.config.ApiVersion))

	v1.GET("/keys", r.paymentController.RSAPublicKey)

	return router
}

func (r *Router) initRoute() *gin.Engine {
	if r.config.Production {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	return gin.New()
}
