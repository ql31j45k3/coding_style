package order

import (
	"context"
	"fmt"

	"github.com/ql31j45k3/coding_style/go/layout/configs"
	utilsDriver "github.com/ql31j45k3/coding_style/go/layout/internal/utils/driver"
	"github.com/ql31j45k3/coding_style/go/layout/internal/utils/tools"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

func newRepositoryOrder() repositoryOrder {
	return &orderMongo{}
}

// repositoryOrder 示範不同操作類型的寫法
type repositoryOrder interface {
	GetOrderID(mongoRS *mongo.Client, account string) (orderData, error)
	GetOrders(mongoRS *mongo.Client, account string) ([]orderData, error)
	GetOrderAggregate(mongoRS *mongo.Client, account string) ([]orderData, error)

	ReplaceOne(mongoRS *mongo.Client, account string, oldData orderData) error

	ExistsOrder(mongoRS *mongo.Client, account string) (bool, error)
	UpdateOrder(mongoRS *mongo.Client, account string) error
	UpdateManyOrder(mongoRS *mongo.Client, account string) error

	DeleteOneOrder(mongoRS *mongo.Client, account string) error

	BulkInsertManyOrder(mongoRS *mongo.Client, account2Order map[string]order) error

	InsertOneOrder(mongoRS *mongo.Client, data order) error
	InsertManyOpAgentSettlementMember(mongoRS *mongo.Client, data []order) error
}

type orderMongo struct {
	_ struct{}
}

func (om *orderMongo) GetOrderID(mongoRS *mongo.Client, account string) (orderData, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), configs.Mongo.GetTimeout())
	defer cancelCtx()

	result := orderData{}
	mongoRSDB := utilsDriver.SetMongoDatabase(mongoRS, result)
	collection := utilsDriver.SetMongoCollection(mongoRSDB, result, "")

	condition := bson.M{
		"account": account,
	}

	if err := collection.FindOne(ctx, condition).Decode(&result); err != nil {
		return result, fmt.Errorf(tools.ErrorFormatFindOne, result.GetCollectionName(""), condition, err)
	}

	return result, nil
}

func (om *orderMongo) GetOrders(mongoRS *mongo.Client, account string) ([]orderData, error) {
	condition := bson.M{
		"account": account,
	}

	model := orderData{}
	var result []orderData
	if err := utilsDriver.FindAndDecode(configs.Mongo.GetTimeout(), configs.Mongo.GetTimeout(), &result, mongoRS, model, "", condition); err != nil {
		return result, fmt.Errorf("utilsDriver.FindAndDecode - %w", err)
	}

	return result, nil
}

func (om *orderMongo) GetOrderAggregate(mongoRS *mongo.Client, account string) ([]orderData, error) {
	match := bson.M{
		"account": account,
	}

	group := bson.M{
		"_id":   "$vendor_code",
		"count": bson.M{"$sum": 1},

		"amount": bson.M{"$sum": "$amount"},
	}

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: match}},
		bson.D{{Key: "$group", Value: group}},
	}

	model := orderData{}
	var result []orderData
	if err := utilsDriver.AggregateAndDecode(configs.Mongo.GetTimeout(), configs.Mongo.GetTimeout(), &result, mongoRS, model, "", pipeline); err != nil {
		return result, fmt.Errorf("utilsDriver.AggregateAndDecode - %w", err)
	}

	return result, nil
}

func (om *orderMongo) ReplaceOne(mongoRS *mongo.Client, account string, oldData orderData) error {
	condition := bson.M{
		"account": account,
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), configs.Mongo.GetTimeout())
	defer cancelCtx()

	model := orderData{}
	mongoRSDB := utilsDriver.SetMongoDatabase(mongoRS, model)
	collection := utilsDriver.SetMongoCollection(mongoRSDB, model, "")

	if _, err := utilsDriver.ReplaceOne(ctx, collection, condition, oldData); err != nil {
		return fmt.Errorf(tools.ErrorFormatReplaceOne, model.GetCollectionName(""), condition, err)
	}

	return nil
}

func (om *orderMongo) ExistsOrder(mongoRS *mongo.Client, account string) (bool, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), configs.Mongo.GetTimeout())
	defer cancelCtx()

	model := orderData{}
	mongoRSDB := utilsDriver.SetMongoDatabase(mongoRS, model)
	collection := utilsDriver.SetMongoCollection(mongoRSDB, model, "")

	condition := bson.D{
		{Key: "account", Value: account},
	}

	exist, err := utilsDriver.ExistsData(ctx, collection, condition)
	if err != nil {
		return false, fmt.Errorf(tools.ErrorFormatExistsData, model.GetCollectionName(""), condition, err)
	}

	return exist, nil
}

func (om *orderMongo) UpdateOrder(mongoRS *mongo.Client, account string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), configs.Mongo.GetTimeout())
	defer cancelCtx()

	model := orderData{}
	mongoRSDB := utilsDriver.SetMongoDatabase(mongoRS, model)
	collection := utilsDriver.SetMongoCollection(mongoRSDB, model, "")

	condition := bson.M{
		"account": account,
	}

	update := bson.M{
		"$set": bson.M{
			"order_id": "abc",
		},
	}

	if _, err := collection.UpdateOne(ctx, condition, update); err != nil {
		return fmt.Errorf(tools.ErrorFormatUpdateOne, model.GetCollectionName(""), condition, update, err)
	}

	return nil
}

func (om *orderMongo) UpdateManyOrder(mongoRS *mongo.Client, account string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), configs.Mongo.GetTimeout())
	defer cancelCtx()

	model := orderData{}
	mongoRSDB := utilsDriver.SetMongoDatabase(mongoRS, model)
	collection := utilsDriver.SetMongoCollection(mongoRSDB, model, "")

	condition := bson.M{
		"account": account,
	}

	update := bson.M{
		"$set": bson.M{
			"order_id": "abc",
		},
	}

	if _, err := collection.UpdateMany(ctx, condition, update); err != nil {
		return fmt.Errorf(tools.ErrorFormatUpdateMany, model.GetCollectionName(""), condition, update, err)
	}

	return nil
}

func (om *orderMongo) DeleteOneOrder(mongoRS *mongo.Client, account string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), configs.Mongo.GetTimeout())
	defer cancelCtx()

	model := orderData{}
	mongoRSDB := utilsDriver.SetMongoDatabase(mongoRS, model)
	collection := utilsDriver.SetMongoCollection(mongoRSDB, model, "")

	condition := bson.M{
		"account": account,
	}

	if _, err := collection.DeleteOne(ctx, condition); err != nil {
		return fmt.Errorf(tools.ErrorFormatDeleteOne, model.GetCollectionName(""), condition, err)
	}

	return nil
}

func (om *orderMongo) BulkInsertManyOrder(mongoRS *mongo.Client, account2Order map[string]order) error {
	if len(account2Order) == 0 {
		return nil
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), configs.Mongo.GetTimeout())
	defer cancelCtx()

	model := orderData{}
	mongoRSDB := utilsDriver.SetMongoDatabase(mongoRS, model)
	collection := utilsDriver.SetMongoCollection(mongoRSDB, model, "")

	writeModels := make([]mongo.WriteModel, 0)

	for _, value := range account2Order {
		writeModel := mongo.NewInsertOneModel()
		writeModels = append(writeModels, writeModel.SetDocument(value))
	}

	if _, err := utilsDriver.BulkWriteOrderedFalse(ctx, collection, writeModels); err != nil {
		return fmt.Errorf(tools.ErrorFormatBulkWrite, model.GetCollectionName(""), err)
	}

	return nil
}

func (om *orderMongo) InsertOneOrder(mongoRS *mongo.Client, data order) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), configs.Mongo.GetTimeout())
	defer cancelCtx()

	model := orderData{}
	mongoRSDB := utilsDriver.SetMongoDatabase(mongoRS, model)
	collection := utilsDriver.SetMongoCollection(mongoRSDB, model, "")

	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return fmt.Errorf(tools.ErrorFormatInsertOne, model.GetCollectionName(""), err)
	}

	return nil
}

func (om *orderMongo) InsertManyOpAgentSettlementMember(mongoRS *mongo.Client, data []order) error {
	if len(data) == 0 {
		return nil
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), configs.Mongo.GetTimeout())
	defer cancelCtx()

	model := orderData{}
	mongoRSDB := utilsDriver.SetMongoDatabase(mongoRS, model)
	collection := utilsDriver.SetMongoCollection(mongoRSDB, model, "")

	var tempData []interface{}

	for i := range data {
		tempData = append(tempData, data[i])
	}

	if _, err := utilsDriver.InsertManyOrderedFalse(ctx, collection, tempData); err != nil {
		return fmt.Errorf(tools.ErrorFormatInsertMany, model.GetCollectionName(""), err)
	}

	return nil
}
