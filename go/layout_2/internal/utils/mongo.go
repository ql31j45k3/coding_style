package utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoPagination 增加 mongodb limit, skip 條件
func MongoPagination(opts options.FindOptions, rowCount, offset int) options.FindOptions {
	opts.SetLimit(int64(rowCount)).SetSkip(int64(offset))
	return opts
}

// MongoBsonDAppend 依照 condition 判斷是否拼湊 Where 條件
func MongoBsonDAppend(filter bson.D, condition bool, query bson.E) bson.D {
	if !condition {
		return filter
	}

	filter = append(filter, query)

	return filter
}
