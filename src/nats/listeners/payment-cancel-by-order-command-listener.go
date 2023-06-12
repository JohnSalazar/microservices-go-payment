package listeners

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"payment/src/application/commands"

	common_nats "github.com/JohnSalazar/microservices-go-common/nats"
	common_service "github.com/JohnSalazar/microservices-go-common/services"
	trace "github.com/JohnSalazar/microservices-go-common/trace/otel"
	"github.com/nats-io/nats.go"
)

type CancelPaymentByOrderCommandListener struct {
	commandHandler *commands.PaymentCommandHandler
	email          common_service.EmailService
	errorHelper    *common_nats.CommandErrorHelper
}

func NewCancelPaymentByOrderCommandListener(
	commandHandler *commands.PaymentCommandHandler,
	email common_service.EmailService,
	errorHelper *common_nats.CommandErrorHelper,
) *CancelPaymentByOrderCommandListener {
	return &CancelPaymentByOrderCommandListener{
		commandHandler: commandHandler,
		email:          email,
		errorHelper:    errorHelper,
	}
}

func (c *CancelPaymentByOrderCommandListener) ProcessCancelPaymentByOrderCommand() nats.MsgHandler {
	return func(msg *nats.Msg) {
		ctx := context.Background()
		_, span := trace.NewSpan(ctx, fmt.Sprintf("publish.%s\n", msg.Subject))
		defer span.End()

		paymentCommand := &commands.CancelPaymentByOrderCommand{}
		err := json.Unmarshal(msg.Data, paymentCommand)
		if c.errorHelper.CheckUnmarshal(msg, err) == nil {
			err = c.commandHandler.CancelPaymentByOrderCommandHandler(ctx, paymentCommand)
			c.errorHelper.CheckCommandError(span, msg, err)
		}

		err = msg.Ack()
		if err != nil {
			log.Printf("stan msg.Ack error: %v\n", err)
		}
	}
}
