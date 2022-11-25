package mongodb

const (
	ErrorFormatFind         = "mongo.Collection Find [ collection: %s, condition: %v, err: %w ]"
	ErrorFormatCursorDecode = "cursor.Decode [ Collection: %s, condition: %v, Error: %w ]"

	ErrorFormatAggregate             = "mongo.Collection Aggregate [ collection: %s, pipeline: %v, err: %w ]"
	ErrorFormatAggregateCursorDecode = "mongo.Collection Aggregate cursor.Decode [ collection: %s, pipeline: %v, err: %w ]"
)
