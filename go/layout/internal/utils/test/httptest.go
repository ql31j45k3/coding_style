package test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http/httptest"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ql31j45k3/coding_style/go/layout/configs"

	"github.com/gin-gonic/gin"
	utilsDriver "github.com/ql31j45k3/coding_style/go/layout/internal/utils/driver"
	"gorm.io/gorm"
)

func Start() (*gin.Engine, *gorm.DB, *mongo.Client, error) {
	if err := configs.Start(); err != nil {
		return nil, nil, nil, err
	}

	dbM, err := utilsDriver.NewPostgresM(configs.Gorm.GetHost(), configs.Gorm.GetUser(), configs.Gorm.GetPassword(),
		configs.Gorm.GetDBName(), configs.Gorm.GetPort(),
		configs.Gorm.GetMaxIdle(), configs.Gorm.GetMaxOpen(), configs.Gorm.GetMaxLifetime(),
		configs.Gorm.GetLogMode())
	if err != nil {
		return nil, nil, nil, err
	}

	mongoRS, err := newMongoRS()
	if err != nil {
		return nil, nil, nil, err
	}

	r := gin.Default()

	return r, dbM, mongoRS, nil
}

func newMongoRS() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), configs.Mongo.GetTimeout())
	defer cancel()

	return utilsDriver.NewMongoDBConnect(ctx, "",
		utilsDriver.WithMongoHosts(configs.Mongo.GetHosts()),
		utilsDriver.WithMongoAuth(configs.Mongo.GetAuthMechanism(), configs.Mongo.GetUsername(), configs.Mongo.GetPassword()),
		utilsDriver.WithMongoReplicaSet(configs.Mongo.GetReplicaSet()),
		utilsDriver.WithMongoPool(configs.Mongo.GetMinPoolSize(), configs.Mongo.GetMaxPoolSize(), configs.Mongo.GetMaxConnIdleTime()),
		utilsDriver.WithMongoPoolMonitor(),
	)
}

// HttptestRequest 根據特定請求 URL 和參數 param
func HttptestRequest(r *gin.Engine, method, uri string, reader io.Reader) (int, []byte, error) {
	req := httptest.NewRequest(method, uri, reader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	defer func() {
		// resp 故意不用參數帶入方式，是因 golangci-lint bodyclose 檢驗不通過，需讓 defer 自行取外面參數才可通過檢驗
		if err := resp.Body.Close(); err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("HttptestRequest - result.Body.Close()")
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return w.Code, body, nil
}
