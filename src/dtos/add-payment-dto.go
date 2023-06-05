package dtos

import "go.mongodb.org/mongo-driver/bson/primitive"

type AddPayment struct {
	ID      primitive.ObjectID `json:"id"`
	OrderID primitive.ObjectID `json:"orderId"`
	Total   float32            `json:"total"`
	Status  uint               `json:"status"`
}
