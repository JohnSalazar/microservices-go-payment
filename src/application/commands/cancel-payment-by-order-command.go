package commands

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CancelPaymentByOrderCommand struct {
	OrderID  primitive.ObjectID `json:"orderId"`
	Status   uint               `json:"status"`
	StatusAt time.Time          `json:"status_at"`
}
