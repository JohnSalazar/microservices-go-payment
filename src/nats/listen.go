package nats

import (
	"payment/src/application/commands"
	"payment/src/nats/listeners"

	"github.com/nats-io/nats.go"
	"github.com/oceano-dev/microservices-go-common/config"

	common_nats "github.com/oceano-dev/microservices-go-common/nats"
	common_service "github.com/oceano-dev/microservices-go-common/services"
)

type listen struct {
	js nats.JetStreamContext
}

const queueGroupName string = "payments-service"

var (
	subscribe          common_nats.Listener
	commandErrorHelper *common_nats.CommandErrorHelper

	paymentCreateCommand        *listeners.PaymentCreateCommandListener
	paymentUpdateStatusCommand  *listeners.UpdateStatusPaymentCommandListener
	cancelPaymentByOrderCommand *listeners.CancelPaymentByOrderCommandListener
)

func NewListen(
	config *config.Config,
	js nats.JetStreamContext,
	paymentCommandHandler *commands.PaymentCommandHandler,
	email common_service.EmailService,
) *listen {
	subscribe = common_nats.NewListener(js)
	commandErrorHelper = common_nats.NewCommandErrorHelper(config, email)

	paymentCreateCommand = listeners.NewPaymentCreateCommandListener(paymentCommandHandler, email, commandErrorHelper)
	paymentUpdateStatusCommand = listeners.NewUpdateStatusPaymentCommandListener(paymentCommandHandler, email, commandErrorHelper)
	cancelPaymentByOrderCommand = listeners.NewCancelPaymentByOrderCommandListener(paymentCommandHandler, email, commandErrorHelper)
	return &listen{
		js: js,
	}
}

func (l *listen) Listen() {
	go subscribe.Listener(string(common_nats.PaymentCreate), queueGroupName, queueGroupName+"_0", paymentCreateCommand.ProcessPaymentCreateCommand())
	go subscribe.Listener(string(common_nats.PaymentUpdate), queueGroupName, queueGroupName+"_1", paymentUpdateStatusCommand.ProcessUpdateStatusPaymentCommand())
	go subscribe.Listener(string(common_nats.PaymentCancel), queueGroupName, queueGroupName+"_2", cancelPaymentByOrderCommand.ProcessCancelPaymentByOrderCommand())
}
