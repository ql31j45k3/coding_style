package tools

const (
	MongoDBTestAPI = "test_api"
)

const (
	DefaultNotAssignInt = -1
)

const (
	MongoDBDuplicateKey = "E11000"

	TimezoneUTC    = "UTC"
	TimezoneTaipei = "Asia/Taipei"

	// YYYY-MM-DD
	TimeFormatDay = "2006-01-02"
	// YYYY-MM-DD hh:mm:ss
	TimeFormatSecond = "2006-01-02 15:04:05"
	// YYYY-MM-DD hh
	TimeFormatHour = "2006-01-02 15:00:00"

	TimeFormatMonth = "2006-01"
	TimeFormatHour2 = "2006-01-02 15"
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
