package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateStatusPaymentByOrder struct {
	OrderID  primitive.ObjectID `json:"orderId"`
	Status   uint               `json:"status"`
	StatusAt time.Time          `json:"status_at"`
}
