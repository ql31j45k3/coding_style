package index

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func StartCreateIndex(mongoRS *mongo.Client) {
	index := newUseCaseIndex(newRepositoryIndex())
	index.DoCreateIndex(mongoRS)
}
