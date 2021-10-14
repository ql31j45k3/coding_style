package tools

const (
	MongoDBTestAPI = "test_api"
)

const (
	ErrorFormatFind         = "mongosql.Find [ Collection: %s, condition: %v, Error: %w ]"
	ErrorFormatFindOne      = "mongosql.FindOne [ Collection: %s, condition: %v, Error: %w ]"
	ErrorFormatCursorDecode = "cursor.Decode [ Collection: %s, condition: %v, Error: %w ]"

	ErrorFormatAggregate             = "mongo.Collection Aggregate [ Collection: %s, pipeline: %v, Error: %w ]"
	ErrorFormatAggregateCursorDecode = "mongo.Collection Aggregate cursor.Decode [ Collection: %s, pipeline: %v, Error: %w ]"

	ErrorFormatReplaceOne = "mongosql.ReplaceOne [ Collection: %s, condition: %v, Error: %w ]"
	ErrorFormatInsertOne  = "mongosql.InsertOne [ Collection: %s, Error: %w ]"
	ErrorFormatInsertMany = "mongosql.InsertMany [ Collection: %s, Error: %w ]"
	ErrorFormatBulkWrite  = "mongosql.BulkWrite [ Collection: %s, Error: %w ]"

	ErrorFormatUpdateOne  = "mongosql.UpdateOne [ Collection: %s, condition: %v, update: %v, Error: %w ]"
	ErrorFormatUpdateMany = "mongosql.UpdateMany [ Collection: %s, condition: %v, update: %v, Error: %w ]"

	ErrorFormatExistsData = "mongosql.ExistsData [ Collection: %s, condition: %v, Error: %w ]"

	ErrorFormatDeleteOne  = "mongosql.DeleteOne [ Collection: %s, condition: %v, Error: %w ]"
	ErrorFormatDeleteMany = "mongosql.DeleteMany [ Collection: %s, condition: %v, Error: %w ]"
)