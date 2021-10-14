package system

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func RegisterRouterHealth(condAPI APIHealthCond) {
	h := health{
		dbM:     condAPI.DBM,
		mongoRS: condAPI.MongoRS,
	}

	condAPI.R.GET("/health", h.check)
}

type health struct {
	_ struct{}

	dbM *gorm.DB

	mongoRS *mongo.Client
}

func (h *health) check(c *gin.Context) {
	if err := h.checkPG(); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := h.checkMongo(); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, "")
}

func (h *health) checkPG() error {
	db, err := h.dbM.DB()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("checkPG - dbM.DB")

		return fmt.Errorf("postgreSQL Connection failed - dbM.DB() - %w", err)
	}

	if err := db.Ping(); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("checkPG - dbM.Ping")

		return fmt.Errorf("postgreSQL Connection failed - db.Ping() - %w", err)
	}

	return nil
}

func (h *health) checkMongo() error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelCtx()

	if err := h.mongoRS.Ping(ctx, readpref.Primary()); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("checkMongo - Ping")

		return fmt.Errorf("mongoRecord Connection failed - mongoRS.Ping - %w", err)
	}

	return nil
}
