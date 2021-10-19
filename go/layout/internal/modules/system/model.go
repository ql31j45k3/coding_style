package system

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type APIHealthCond struct {
	dig.In

	R *gin.Engine

	DBM     *gorm.DB      `name:"dbM"`
	MongoRS *mongo.Client `name:"mongoRS"`
}

type APIDocCond struct {
	dig.In

	R *gin.Engine
}
