package commands

import (
	"context"
	"errors"
	"payment/src/application/events"
	"payment/src/dtos"
	"payment/src/models"
	"payment/src/repositories/interfaces"
	"payment/src/validators"
	"strings"
	"time"

	common_models "github.com/JohnSalazar/microservices-go-common/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentCommandHandler struct {
	paymentRepository   interfaces.PaymentRepository
	paymentEventHandler *events.PaymentEventHandler
}

func NewPaymentCommandHandler(
	paymentRepository interfaces.PaymentRepository,
	paymentEventHandler *events.PaymentEventHandler,
) *PaymentCommandHandler {
	return &PaymentCommandHandler{
		paymentRepository:   paymentRepository,
		paymentEventHandler: paymentEventHandler,
	}
}

func (payment *PaymentCommandHandler) CreatePaymentCommandHandler(ctx context.Context, command *CreatePaymentCommand) error {

	paymentDto := &dtos.AddPayment{
		OrderID: command.OrderID,
		Total:   command.Total,
		Status:  uint(common_models.SentForPaymentConfirmation),
	}

	result := validators.ValidateAddPayment(paymentDto)
	if result != nil {
		return errors.New(strings.Join(result.([]string), ""))
	}

	paymentModel := &models.Payment{
		ID:        primitive.NewObjectID(),
		OrderID:   paymentDto.OrderID,
		Total:     paymentDto.Total,
		Status:    paymentDto.Status,
		StatusAt:  time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}

	paymentExists, _ := payment.paymentRepository.FindByOrderID(ctx, paymentDto.OrderID)
	if paymentExists != nil {
		return errors.New("already a payment for this order")
	}

	paymentModel, err := payment.paymentRepository.Create(ctx, paymentModel)
	if err != nil {
		return err
	}

	paymentEvent := &events.PaymentCreatedEvent{
		ID:         paymentModel.ID,
		OrderID:    paymentModel.OrderID,
		Total:      paymentModel.Total,
		CardNumber: command.CardNumber,
		Kid:        command.Kid,
		Status:     paymentModel.Status,
		StatusAt:   paymentModel.StatusAt,
		CreatedAt:  paymentModel.CreatedAt,
		Version:    paymentModel.Version,
	}

	go payment.paymentEventHandler.PaymentCreatedEventHandler(ctx, paymentEvent)

	return nil
}

func (payment *PaymentCommandHandler) UpdateStatusPaymentCommandHandler(ctx context.Context, command *UpdateStatusPaymentCommand) error {
	paymentDto := &dtos.UpdateStatusPayment{
		ID:       command.ID,
		Status:   command.Status,
		StatusAt: command.StatusAt,
	}

	result := validators.ValidateUpdateStatusPayment(paymentDto)
	if result != nil {
		return errors.New(strings.Join(result.([]string), ""))
	}

	paymentExists, err := payment.paymentRepository.FindByID(ctx, paymentDto.ID)
	if err != nil {
		return err
	}

	paymentModel := &models.Payment{
		ID:        paymentDto.ID,
		OrderID:   paymentExists.OrderID,
		Total:     paymentExists.Total,
		Status:    paymentDto.Status,
		StatusAt:  paymentDto.StatusAt,
		UpdatedAt: time.Now().UTC(),
		Version:   paymentExists.Version,
	}

	paymentModel, err = payment.paymentRepository.Update(ctx, paymentModel)
	if err != nil {
		return err
	}

	paymentEvent := &events.PaymentStatusUpdatedEvent{
		ID:        paymentModel.ID,
		OrderID:   paymentModel.OrderID,
		Status:    paymentModel.Status,
		StatusAt:  paymentModel.StatusAt,
		UpdatedAt: paymentModel.UpdatedAt,
		Version:   paymentModel.Version,
	}

	go payment.paymentEventHandler.PaymentStatusUpdatedEventHandler(ctx, paymentEvent)

	return nil
}

func (payment *PaymentCommandHandler) CancelPaymentByOrderCommandHandler(ctx context.Context, command *CancelPaymentByOrderCommand) error {
	paymentDto := &dtos.UpdateStatusPaymentByOrder{
		OrderID:  command.OrderID,
		Status:   command.Status,
		StatusAt: command.StatusAt,
	}

	result := validators.ValidateUpdateStatusPaymentByOrder(paymentDto)
	if result != nil {
		return errors.New(strings.Join(result.([]string), ""))
	}

	paymentExists, err := payment.paymentRepository.FindByOrderID(ctx, paymentDto.OrderID)
	if err != nil {
		return err
	}

	paymentModel := &models.Payment{
		ID:        paymentExists.ID,
		OrderID:   paymentDto.OrderID,
		Total:     paymentExists.Total,
		Status:    paymentDto.Status,
		StatusAt:  paymentDto.StatusAt,
		UpdatedAt: time.Now().UTC(),
		Version:   paymentExists.Version,
	}

	paymentModel, err = payment.paymentRepository.Update(ctx, paymentModel)
	if err != nil {
		return err
	}

	paymentEvent := &events.PaymentStatusUpdatedEvent{
		ID:        paymentModel.ID,
		OrderID:   paymentModel.OrderID,
		Status:    paymentModel.Status,
		StatusAt:  paymentModel.StatusAt,
		UpdatedAt: paymentModel.UpdatedAt,
		Version:   paymentModel.Version,
	}

	go payment.paymentEventHandler.PaymentStatusUpdatedByOrderEventHandler(ctx, paymentEvent)

	return nil
}
