package driver

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresM(host, user, password, dbname, port string,
	maxIdle, maxOpen int, maxLifetime time.Duration, logMode logger.LogLevel) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, "disable")

	dbM, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.New(log.StandardLogger(), logger.Config{Colorful: false}).LogMode(logMode),
	})
	if err != nil {
		return nil, err
	}

	db, err := dbM.DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdle)
	db.SetMaxOpenConns(maxOpen)
	db.SetConnMaxLifetime(maxLifetime)

	return dbM, nil
}
