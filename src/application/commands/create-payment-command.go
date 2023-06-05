package commands

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatePaymentCommand struct {
	ID         primitive.ObjectID `json:"id"`
	OrderID    primitive.ObjectID `json:"orderId"`
	Total      float32            `json:"total"`
	CardNumber []byte             `json:"cardNumber"`
	Kid        string             `json:"kid"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty"`
	Version    uint               `json:"version"`
	Deleted    bool               `json:"deleted,omitempty"`
}
