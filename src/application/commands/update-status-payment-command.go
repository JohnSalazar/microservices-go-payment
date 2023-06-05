package commands

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateStatusPaymentCommand struct {
	ID       primitive.ObjectID `json:"id"`
	Status   uint               `json:"status"`
	StatusAt time.Time          `json:"status_at"`
}
