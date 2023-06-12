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

type UpdateStatusPaymentCommandListener struct {
	commandHandler *commands.PaymentCommandHandler
	email          common_service.EmailService
	errorHelper    *common_nats.CommandErrorHelper
}

func NewUpdateStatusPaymentCommandListener(
	commandHandler *commands.PaymentCommandHandler,
	email common_service.EmailService,
	errorHelper *common_nats.CommandErrorHelper,
) *UpdateStatusPaymentCommandListener {
	return &UpdateStatusPaymentCommandListener{
		commandHandler: commandHandler,
		email:          email,
		errorHelper:    errorHelper,
	}
}

func (c *UpdateStatusPaymentCommandListener) ProcessUpdateStatusPaymentCommand() nats.MsgHandler {
	return func(msg *nats.Msg) {
		ctx := context.Background()
		_, span := trace.NewSpan(ctx, fmt.Sprintf("publish.%s\n", msg.Subject))
		defer span.End()

		paymentCommand := &commands.UpdateStatusPaymentCommand{}
		err := json.Unmarshal(msg.Data, paymentCommand)
		if c.errorHelper.CheckUnmarshal(msg, err) == nil {
			err = c.commandHandler.UpdateStatusPaymentCommandHandler(ctx, paymentCommand)
			c.errorHelper.CheckCommandError(span, msg, err)
		}

		err = msg.Ack()
		if err != nil {
			log.Printf("stan msg.Ack error: %v\n", err)
		}
	}
}
