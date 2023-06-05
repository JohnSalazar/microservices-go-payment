package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	OrderID   primitive.ObjectID `bson:"order_id" json:"orderId"`
	Total     float32            `bson:"total" json:"total"`
	Status    uint               `bson:"status" json:"status"`
	StatusAt  time.Time          `bson:"status_at" json:"status_at"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
	Version   uint               `bson:"version" json:"version"`
	Deleted   bool               `bson:"deleted" json:"deleted,omitempty"`
}
