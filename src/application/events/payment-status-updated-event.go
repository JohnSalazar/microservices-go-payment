package events

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentStatusUpdatedEvent struct {
	ID        primitive.ObjectID `json:"id"`
	OrderID   primitive.ObjectID `json:"orderId"`
	Status    uint               `json:"status"`
	StatusAt  time.Time          `json:"status_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Version   uint               `json:"version"`
}
