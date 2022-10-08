package index

import (
	"go.mongodb.org/mongo-driver/mongo"

	log "github.com/sirupsen/logrus"
)

func newUseCaseIndex(repositoryIndex repositoryIndex) useCaseIndex {
	return &index{
		repositoryIndex: repositoryIndex,
	}
}

type useCaseIndex interface {
	DoCreateIndex(mongoRS *mongo.Client)
}

type index struct {
	repositoryIndex
}

func (i *index) DoCreateIndex(mongoRS *mongo.Client) {
	i.CreateIndexOrder(mongoRS)
	log.Info("DoCreateIndex - CreateIndexOrder")
}
