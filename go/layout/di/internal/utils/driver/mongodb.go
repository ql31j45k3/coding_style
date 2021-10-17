package driver

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/ql31j45k3/coding_style/go/layout/di/internal/utils/tools"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase interface {
	GetDatabase() string
}

type MongoDatabaseFunc func() string

func (f MongoDatabaseFunc) GetDatabase() string { return f() }

type MongoCollection interface {
	GetCollectionName(opCode string) string
}

type MongoCollectionFunc func(opCode string) string

func (f MongoCollectionFunc) GetCollectionName(opCode string) string { return f(opCode) }

type MongoDatabaseAndCollection interface {
	MongoDatabase
	MongoCollection
}

func SetMongoDatabase(client *mongo.Client, model MongoDatabase, opts ...*options.DatabaseOptions) *mongo.Database {
	return client.Database(model.GetDatabase(), opts...)
}

func SetMongoCollection(db *mongo.Database, model MongoCollection, opCode string, opts ...*options.CollectionOptions) *mongo.Collection {
	return db.Collection(model.GetCollectionName(opCode), opts...)
}

func Disconnect(ctx context.Context, client *mongo.Client) error {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelCtx()
	}

	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("client.Disconnect - %w", err)
	}

	return nil
}

func ExistsCollection(ctx context.Context, db *mongo.Database, filter bson.M) (bool, error) {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelCtx()
	}

	collections, err := db.ListCollectionNames(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("db.ListCollectionNames - %w", err)
	}

	if len(collections) == 0 {
		return false, nil
	}

	return true, nil
}

func CreateManyIndexMaxTime(ctx context.Context, collection *mongo.Collection, models []mongo.IndexModel) ([]string, error) {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelCtx()
	}

	// 設定建立 index 上限時間
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	return collection.Indexes().CreateMany(ctx, models, opts)
}

func ExistsCollectionAndCreateManyIndexMaxTime(existsCollectionTimeout time.Duration, createManyIndexTimeout time.Duration,
	mongoRS *mongo.Client, mongoDatabase, mongoCollection, opCode string, models []mongo.IndexModel) ([]string, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), existsCollectionTimeout)
	defer cancelCtx()

	modelMongoDatabase := func() string {
		return mongoDatabase
	}

	modelMongoCollection := func(opCode string) string {
		return opCode + mongoCollection
	}

	mongoRSDB := SetMongoDatabase(mongoRS, MongoDatabaseFunc(modelMongoDatabase))

	filter := bson.M{"name": MongoCollectionFunc(modelMongoCollection).GetCollectionName(opCode)}
	ok, err := ExistsCollection(ctx, mongoRSDB, filter)
	if err != nil {
		return []string{}, fmt.Errorf("ExistsCollectionSingle - %w", err)
	}

	// 已存在，不建立 index
	if ok {
		return []string{}, nil
	}

	ctxIndex, cancelCtxCursor := context.WithTimeout(context.Background(), createManyIndexTimeout)
	defer cancelCtxCursor()

	collection := SetMongoCollection(mongoRSDB, MongoCollectionFunc(modelMongoCollection), opCode)
	indexName, err := CreateManyIndexMaxTime(ctxIndex, collection, models)
	if err != nil {
		return []string{}, fmt.Errorf("CreateManyIndexMaxTime - %w", err)
	}

	return indexName, nil
}

func ExistsData(ctx context.Context, collection *mongo.Collection, condition bson.D) (bool, error) {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelCtx()
	}

	count, err := collection.CountDocuments(ctx, condition)
	if err != nil {
		return false, fmt.Errorf("collection.CountDocuments - %w", err)
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func ReplaceOne(ctx context.Context, collection *mongo.Collection, filter bson.M, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), 30*time.Second)
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
		ctx, cancelCtx = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelCtx()
	}

	// 有錯誤發生持續寫入
	opts := options.InsertMany().SetOrdered(false)
	return collection.InsertMany(ctx, documents, opts)
}

func BulkWriteOrderedFalse(ctx context.Context, collection *mongo.Collection, models []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancelCtx()
	}

	// 有錯誤發生持續寫入
	opts := options.BulkWrite().SetOrdered(false)
	return collection.BulkWrite(ctx, models, opts)
}

func FindAndDecode(findTimeout time.Duration, cursorTimeout time.Duration, result interface{},
	mongoRS *mongo.Client, model MongoDatabaseAndCollection, opCode string, filter interface{},
	opts ...*options.FindOptions) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), findTimeout)
	defer cancelCtx()

	mongoRSDB := SetMongoDatabase(mongoRS, model)
	collection := SetMongoCollection(mongoRSDB, model, opCode)

	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return fmt.Errorf(tools.ErrorFormatFind, model.GetCollectionName(opCode), filter, err)
	}

	ctxCursor, cancelCtxCursor := context.WithTimeout(context.Background(), cursorTimeout)
	defer cancelCtxCursor()

	if err := CursorDecode(ctxCursor, result, cursor); err != nil {
		return fmt.Errorf(tools.ErrorFormatCursorDecode, model.GetCollectionName(opCode), filter, err)
	}

	return nil
}

func AggregateAndDecode(aggregateTimeout time.Duration, cursorTimeout time.Duration, result interface{},
	mongoRS *mongo.Client, model MongoDatabaseAndCollection, opCode string, pipeline bson.A, opts ...*options.AggregateOptions) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), aggregateTimeout)
	defer cancelCtx()

	mongoRSDB := SetMongoDatabase(mongoRS, model)
	collection := SetMongoCollection(mongoRSDB, model, opCode)

	if opts == nil {
		opts = append(opts, options.Aggregate().SetAllowDiskUse(true))
	}

	cursor, err := collection.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return fmt.Errorf(tools.ErrorFormatAggregate, model.GetCollectionName(opCode), pipeline, err)
	}

	ctxCursor, cancelCtxCursor := context.WithTimeout(context.Background(), cursorTimeout)
	defer cancelCtxCursor()

	if err := CursorDecode(ctxCursor, result, cursor); err != nil {
		return fmt.Errorf(tools.ErrorFormatAggregateCursorDecode, model.GetCollectionName(opCode), pipeline, err)
	}

	return nil
}

type CursorConv interface {
	Conv(val interface{}) error
}

func CursorDecode(ctx context.Context, v interface{}, cursor *mongo.Cursor) error {
	var cancelCtx context.CancelFunc
	if ctx == nil {
		ctx, cancelCtx = context.WithTimeout(context.Background(), 30*time.Second)
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
				return fmt.Errorf("cursor.Decode - %w", err)
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
				return fmt.Errorf("cursor.Decode - %w", err)
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
