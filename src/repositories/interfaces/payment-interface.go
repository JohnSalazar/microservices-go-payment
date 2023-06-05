package interfaces

import (
	"context"
	"payment/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentRepository interface {
	FindByOrderID(ctx context.Context, orderID primitive.ObjectID) (*models.Payment, error)
	FindByID(ctx context.Context, ID primitive.ObjectID) (*models.Payment, error)
	Create(ctx context.Context, payment *models.Payment) (*models.Payment, error)
	Update(ctx context.Context, payment *models.Payment) (*models.Payment, error)
	Delete(ctx context.Context, ID primitive.ObjectID) error
}
