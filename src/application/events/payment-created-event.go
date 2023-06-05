package events

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentCreatedEvent struct {
	ID         primitive.ObjectID `json:"id"`
	OrderID    primitive.ObjectID `json:"orderId"`
	Total      float32            `json:"total"`
	CardNumber []byte             `json:"cardNumber"`
	Kid        string             `json:"kid"`
	Status     uint               `json:"status"`
	StatusAt   time.Time          `json:"status_at"`
	CreatedAt  time.Time          `json:"created_at"`
	Version    uint               `json:"version"`
}
