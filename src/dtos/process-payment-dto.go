package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProcessPayment struct {
	ID         primitive.ObjectID `json:"id"`
	Total      float32            `json:"total"`
	CardNumber []byte             `json:"cardNumber"`
	Kid        string             `json:"kid"`
	VerifiedAt time.Time          `json:"verified_at"`
}
