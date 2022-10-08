package index

import (
	"github.com/ql31j45k3/coding_style/go/layout/configs"
	"github.com/ql31j45k3/coding_style/go/layout/internal/modules/order"
	utilsDriver "github.com/ql31j45k3/coding_style/go/layout/internal/utils/driver"
	"github.com/ql31j45k3/coding_style/go/layout/internal/utils/tools"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	MongoOrder = order.MongoOrder
)

func newRepositoryIndex() repositoryIndex {
	return &indexMongo{}
}

type repositoryIndex interface {
	CreateIndexOrder(mongoRS *mongo.Client)
}

type indexMongo struct {
}

func (im *indexMongo) CreateIndexOrder(mongoRS *mongo.Client) {
	indexModel := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{Key: "order_id", Value: bsonx.Int32(1)}},
		},
	}

	_, err := utilsDriver.ExistsCollectionAndCreateManyIndexMaxTime(configs.Mongo.GetTimeout(), configs.Mongo.GetTimeout(),
		mongoRS, tools.MongoDBTestAPI, MongoOrder, "", indexModel)
	if err != nil {
		log.WithFields(log.Fields{
			"err":             err,
			"mongoDatabase":   tools.MongoDBTestAPI,
			"mongoCollection": MongoOrder,
		}).Error("CreateIndexOrder - utilsDriver.ExistsCollectionAndCreateManyIndexMaxTime")
		return
	}
}
