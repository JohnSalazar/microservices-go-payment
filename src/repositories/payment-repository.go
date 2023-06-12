package repositories

import (
	"context"
	"time"

	"payment/src/models"

	"github.com/JohnSalazar/microservices-go-common/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type paymentRepository struct {
	database *mongo.Database
}

func NewPaymentRepository(
	database *mongo.Database,
) *paymentRepository {
	return &paymentRepository{
		database: database,
	}
}

func (r *paymentRepository) collectionName() string {
	return "payments"
}

func (r *paymentRepository) collection() *mongo.Collection {
	return r.database.Collection(r.collectionName())
}

func (r *paymentRepository) findOne(ctx context.Context, filter interface{}) (*models.Payment, error) {
	findOneOptions := options.FindOneOptions{}
	findOneOptions.SetSort(bson.M{"version": -1})

	newFilter := map[string]interface{}{
		"deleted": false,
	}
	mergeFilter := helpers.MergeFilters(newFilter, filter)

	payment := &models.Payment{}
	err := r.collection().FindOne(ctx, mergeFilter, &findOneOptions).Decode(payment)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *paymentRepository) findOneAndUpdate(ctx context.Context, filter interface{}, fields interface{}) *mongo.SingleResult {
	findOneAndUpdateOptions := options.FindOneAndUpdateOptions{}
	findOneAndUpdateOptions.SetReturnDocument(options.After)

	result := r.collection().FindOneAndUpdate(ctx, filter, bson.M{"$set": fields}, &findOneAndUpdateOptions)

	return result
}

func (r *paymentRepository) FindByOrderID(ctx context.Context, orderID primitive.ObjectID) (*models.Payment, error) {
	filter := bson.M{"order_id": orderID}

	return r.findOne(ctx, filter)
}

func (r *paymentRepository) FindByID(ctx context.Context, ID primitive.ObjectID) (*models.Payment, error) {
	filter := bson.M{"_id": ID}

	return r.findOne(ctx, filter)
}

func (r *paymentRepository) Create(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	payment.CreatedAt = time.Now().UTC()

	fields := bson.M{
		"_id":        payment.ID,
		"order_id":   payment.OrderID,
		"total":      payment.Total,
		"status":     payment.Status,
		"status_at":  payment.StatusAt,
		"created_at": payment.CreatedAt,
		"version":    0,
		"deleted":    false,
	}

	_, err := r.collection().InsertOne(ctx, fields)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	payment.Version++
	payment.UpdatedAt = time.Now().UTC()

	fields := bson.M{
		"total":      payment.Total,
		"status":     payment.Status,
		"status_at":  payment.StatusAt,
		"updated_at": payment.UpdatedAt,
		"version":    payment.Version,
	}

	filter := r.filterUpdate(payment)

	result := r.findOneAndUpdate(ctx, filter, fields)
	if result.Err() != nil {
		return nil, result.Err()
	}

	modelPayment := &models.Payment{}
	err := result.Decode(modelPayment)

	return modelPayment, err
}

func (r *paymentRepository) Delete(ctx context.Context, ID primitive.ObjectID) error {
	filter := bson.M{"_id": ID}

	fields := bson.M{"deleted": true}

	result := r.findOneAndUpdate(ctx, filter, fields)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (r *paymentRepository) filterUpdate(payment *models.Payment) interface{} {
	// objectId, _ := primitive.ObjectIDFromHex(payment.ID.String())
	filter := bson.M{
		"_id":     payment.ID,
		"version": payment.Version - 1,
	}

	return filter
}
