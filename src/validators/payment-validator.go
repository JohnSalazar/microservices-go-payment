package validators

import (
	"payment/src/dtos"
	"time"

	common_validator "github.com/oceano-dev/microservices-go-common/validators"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type addPayment struct {
	OrderID primitive.ObjectID `from:"orderId" json:"orderId" validate:"required"`
	Total   float32            `from:"total" json:"total" validate:"required"`
}

type updateStatusPayment struct {
	ID       primitive.ObjectID `from:"id" json:"id" validate:"required"`
	Status   uint               `from:"status" json:"status" validate:"required"`
	StatusAt time.Time          `from:"status_at" json:"status_at" validate:"required"`
}

type updateStatusPaymentByOrder struct {
	OrderID  primitive.ObjectID `from:"orderId" json:"orderId" validate:"required"`
	Status   uint               `from:"status" json:"status" validate:"required"`
	StatusAt time.Time          `from:"status_at" json:"status_at" validate:"required"`
}

func ValidateAddPayment(fields *dtos.AddPayment) interface{} {
	addPayment := addPayment{
		OrderID: fields.OrderID,
		Total:   fields.Total,
	}

	err := common_validator.Validate(addPayment)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUpdateStatusPayment(fields *dtos.UpdateStatusPayment) interface{} {
	updateStatusPayment := updateStatusPayment{
		ID:       fields.ID,
		Status:   fields.Status,
		StatusAt: fields.StatusAt,
	}

	err := common_validator.Validate(updateStatusPayment)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUpdateStatusPaymentByOrder(fields *dtos.UpdateStatusPaymentByOrder) interface{} {
	updateStatusPaymentByOrder := updateStatusPaymentByOrder{
		OrderID:  fields.OrderID,
		Status:   fields.Status,
		StatusAt: fields.StatusAt,
	}

	err := common_validator.Validate(updateStatusPaymentByOrder)
	if err != nil {
		return err
	}

	return nil
}
