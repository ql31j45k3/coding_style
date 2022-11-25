package mongodb

import (
	"context"
	"errors"
	"fmt"
	"layout_2/internal/utils"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase interface {
	GetDatabase() string
}

type MongoDatabaseFunc func() string

func (f MongoDatabaseFunc) GetDatabase() string {
	return f()
}

type MongoCollection interface {
	GetCollectionName() string
}

type MongoCollectionFunc func() string

func (f MongoCollectionFunc) GetCollectionName() string {
	return f()
}

type MongoDatabaseAndCollection interface {
	MongoDatabase
	MongoCollection
}

func SetMongoDatabase(client *mongo.Client, model MongoDatabase, opts ...*options.DatabaseOptions) *mongo.Database {
	return client.Database(model.GetDatabase(), opts...)
}

func SetMongoCollection(db *mongo.Database, model MongoCollection, opts ...*options.CollectionOptions) *mongo.Collection {
	return db.Collection(model.GetCollectionName(), opts...)
}

func Disconnect(ctx context.Context, client *mongo.Client) error {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), utils.Time30S)
		defer cancelCtx()
	}

	if err := client.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func ExistsCollection(ctx context.Context, db *mongo.Database, filter interface{}) (bool, error) {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), utils.Time30S)
		defer cancelCtx()
	}

	collections, err := db.ListCollectionNames(ctx, filter)
	if err != nil {
		return false, err
	}

	if len(collections) == 0 {
		return false, nil
	}

	return true, nil
}

func CreateManyIndexMaxTime(ctx context.Context, maxTimeout time.Duration, client *mongo.Client,
	mongoDatabase, mongoCollection string, models []mongo.IndexModel) ([]string, error) {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), utils.Time30S)
		defer cancelCtx()
	}

	modelMongoDatabase := func() string {
		return mongoDatabase
	}

	modelMongoCollection := func() string {
		return mongoCollection
	}

	mongoDB := SetMongoDatabase(client, MongoDatabaseFunc(modelMongoDatabase))

	filter := bson.M{"name": MongoCollectionFunc(modelMongoCollection).GetCollectionName()}
	ok, err := ExistsCollection(ctx, mongoDB, filter)
	if err != nil {
		return []string{}, err
	}

	// 已存在，不建立 index
	if ok {
		return []string{}, nil
	}

	collection := SetMongoCollection(mongoDB, MongoCollectionFunc(modelMongoCollection))

	indexName, err := createManyIndexMaxTime(ctx, maxTimeout, collection, models)
	if err != nil {
		return []string{}, err
	}

	return indexName, nil
}

func createManyIndexMaxTime(ctx context.Context, maxTimeout time.Duration,
	collection *mongo.Collection, models []mongo.IndexModel) ([]string, error) {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), utils.Time30S)
		defer cancelCtx()
	}

	// 設定建立 index 上限時間
	opts := options.CreateIndexes().SetMaxTime(maxTimeout)
	return collection.Indexes().CreateMany(ctx, models, opts)
}

func ExistsData(ctx context.Context, collection *mongo.Collection, filter interface{}) (bool, error) {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), utils.Time30S)
		defer cancelCtx()
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func ReplaceOne(ctx context.Context, collection *mongo.Collection, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), utils.Time30S)
		defer cancelCtx()
	}

	if len(opts) == 0 {
		opts = append(opts, &options.ReplaceOptions{Upsert: &[]bool{true}[0]})
	}

	return collection.ReplaceOne(ctx, filter, replacement, opts...)
}

func InsertManyOrderedFalse(ctx context.Context, collection *mongo.Collection, documents []interface{}) (*mongo.InsertManyResult, error) {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), utils.Time30S)
		defer cancelCtx()
	}

	// 有錯誤發生持續寫入
	opts := options.InsertMany().SetOrdered(false)
	return collection.InsertMany(ctx, documents, opts)
}

func BulkWriteOrderedFalse(ctx context.Context, collection *mongo.Collection, models []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), utils.Time30S)
		defer cancelCtx()
	}

	// 有錯誤發生持續寫入
	opts := options.BulkWrite().SetOrdered(false)
	return collection.BulkWrite(ctx, models, opts)
}

func FindAndDecode(ctx context.Context, result interface{}, client *mongo.Client,
	model MongoDatabaseAndCollection, filter interface{}, opts ...*options.FindOptions) error {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), utils.Time30S)
		defer cancelCtx()
	}

	mongoDB := SetMongoDatabase(client, model)
	collection := SetMongoCollection(mongoDB, model)

	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return fmt.Errorf(ErrorFormatFind, model.GetCollectionName(), filter, err)
	}

	if err := CursorDecode(ctx, result, cursor); err != nil {
		return fmt.Errorf(ErrorFormatCursorDecode, model.GetCollectionName(), filter, err)
	}

	return nil
}

func AggregateAndDecode(ctx context.Context, result interface{}, client *mongo.Client,
	model MongoDatabaseAndCollection, pipeline interface{}, opts ...*options.AggregateOptions) error {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), utils.Time30S)
		defer cancelCtx()
	}

	mongoDB := SetMongoDatabase(client, model)
	collection := SetMongoCollection(mongoDB, model)

	if opts == nil {
		opts = append(opts, options.Aggregate().SetAllowDiskUse(true))
	}

	cursor, err := collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return fmt.Errorf(ErrorFormatAggregate, model.GetCollectionName(), pipeline, err)
	}

	if err := CursorDecode(ctx, result, cursor); err != nil {
		return fmt.Errorf(ErrorFormatAggregateCursorDecode, model.GetCollectionName(), pipeline, err)
	}

	return nil
}

type CursorConv interface {
	Conv(val interface{}) error
}

func CursorDecode(ctx context.Context, v interface{}, cursor *mongo.Cursor) error {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), utils.Time30S)
		defer cancelCtx()
	}

	vType := reflect.TypeOf(v)

	// 必須為 Ptr 型態，才可有效修改值
	if vType.Kind() != reflect.Ptr {
		return errors.New("CursorDecode need kind is Ptr")
	}

	if vType.Elem().Kind() == reflect.Slice {
		rspVale := reflect.MakeSlice(vType.Elem(), 0, 0)

		for cursor.Next(ctx) {
			cursorResult := reflect.New(rspVale.Type().Elem())
			if err := cursor.Decode(cursorResult.Interface()); err != nil {
				return err
			}

			v, ok := cursorResult.Interface().(CursorConv)
			if ok {
				if err := v.Conv(cursorResult.Interface()); err != nil {
					return errors.New("v.Conv(cursorResult.Interface) fail")
				}
			}

			rspVale = reflect.Append(rspVale, cursorResult.Elem())
		}

		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(rspVale.Interface()))
	} else if vType.Elem().Kind() == reflect.Struct {
		for cursor.Next(ctx) {
			cursorResult := reflect.New(vType.Elem())
			if err := cursor.Decode(cursorResult.Interface()); err != nil {
				return err
			}

			v2, ok := cursorResult.Interface().(CursorConv)
			if ok {
				if err := v2.Conv(cursorResult.Interface()); err != nil {
					return errors.New("v.Conv(cursorResult.Interface) fail")
				}
			}

			reflect.ValueOf(v).Elem().Set(cursorResult.Elem())
		}
	} else {
		return errors.New("only kind Struct or Slice")
	}

	return nil
}
