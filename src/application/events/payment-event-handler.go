package events

import (
	"context"
	"encoding/json"
	"payment/src/dtos"
	"payment/src/tasks/interfaces"
	"time"

	common_models "github.com/oceano-dev/microservices-go-common/models"
	common_nats "github.com/oceano-dev/microservices-go-common/nats"
)

type PaymentEventHandler struct {
	paymentTask interfaces.VerifyPaymentTask
	publisher   common_nats.Publisher
}

func NewPaymentEventHandler(
	paymentTask interfaces.VerifyPaymentTask,
	publisher common_nats.Publisher,
) *PaymentEventHandler {
	return &PaymentEventHandler{
		paymentTask: paymentTask,
		publisher:   publisher,
	}
}

func (payment *PaymentEventHandler) PaymentCreatedEventHandler(ctx context.Context, event *PaymentCreatedEvent) error {

	updateStatusOrder := &dtos.UpdateStatusOrder{
		ID:       event.OrderID,
		Status:   uint(common_models.AwaitingPaymentConfirmation),
		StatusAt: event.StatusAt,
	}

	data, _ := json.Marshal(updateStatusOrder)
	err := payment.publisher.Publish(string(common_nats.OrderStatus), data)
	if err != nil {
		return err
	}

	// processar pagamento
	processPayment := &dtos.ProcessPayment{
		ID:         event.ID,
		Total:      event.Total,
		CardNumber: event.CardNumber,
		Kid:        event.Kid,
		VerifiedAt: time.Now().UTC(),
	}
	payment.paymentTask.AddPayment(processPayment)

	return nil
}

func (payment *PaymentEventHandler) PaymentStatusUpdatedEventHandler(ctx context.Context, event *PaymentStatusUpdatedEvent) error {

	updateStatusOrder := &dtos.UpdateStatusOrder{
		ID:       event.OrderID,
		Status:   event.Status,
		StatusAt: event.StatusAt,
	}

	data, _ := json.Marshal(updateStatusOrder)
	err := payment.publisher.Publish(string(common_nats.OrderStatus), data)
	if err != nil {
		return err
	}

	return nil
}

func (payment *PaymentEventHandler) PaymentStatusUpdatedByOrderEventHandler(ctx context.Context, event *PaymentStatusUpdatedEvent) error {

	return nil
}
