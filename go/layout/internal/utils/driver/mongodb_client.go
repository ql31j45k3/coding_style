package driver

import (
	"context"
	"fmt"
	"time"

	"github.com/ql31j45k3/coding_style/go/layout/internal/utils/tools"

	"go.mongodb.org/mongo-driver/event"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoOption func(opt *options.ClientOptions) *options.ClientOptions

func NewMongoDBConnect(ctx context.Context, uri string, mongoOptions ...mongoOption) (*mongo.Client, error) {
	clientOpt := options.Client()
	if tools.IsNotEmpty(uri) {
		clientOpt.ApplyURI(uri)
	}

	for _, option := range mongoOptions {
		option(clientOpt)
	}

	return mongo.Connect(ctx, clientOpt)
}

func WithMongoHosts(hosts []string) mongoOption {
	return func(opt *options.ClientOptions) *options.ClientOptions {
		if len(hosts) > 0 {
			opt.SetHosts(hosts)
		}

		return opt
	}
}

func WithMongoAuth(authMechanism, username, password string) mongoOption {
	return func(opt *options.ClientOptions) *options.ClientOptions {
		if authMechanism == "Direct" {
			return opt
		}

		if authMechanism == "PLAIN" {
			opt.SetAuth(options.Credential{
				AuthMechanism: authMechanism,
				Username:      username,
				Password:      password,
			})

			return opt
		}

		if authMechanism == "SCRAM" {
			opt.SetAuth(options.Credential{
				Username: username,
				Password: password,
			})

			return opt
		}

		return opt
	}
}

func WithMongoReplicaSet(replicaSet string) mongoOption {
	return func(opt *options.ClientOptions) *options.ClientOptions {
		if tools.IsNotEmpty(replicaSet) {
			opt.SetReplicaSet(replicaSet)
		}

		return opt
	}
}

func WithMongoPool(minPoolSize, maxPoolSize uint64, maxConnIdleTime time.Duration) mongoOption {
	return func(opt *options.ClientOptions) *options.ClientOptions {
		opt.SetMinPoolSize(minPoolSize)
		opt.SetMaxPoolSize(maxPoolSize)
		opt.SetMaxConnIdleTime(maxConnIdleTime)

		return opt
	}
}

func WithMongoPoolMonitor() mongoOption {
	return func(opt *options.ClientOptions) *options.ClientOptions {
		var po *event.MonitorPoolOptions
		m := &event.PoolMonitor{
			Event: func(poolEvent *event.PoolEvent) {
				log.WithFields(log.Fields{
					"poolEvent val": fmt.Sprintf("%+v", poolEvent),
				}).Debug("poolEvent")

				if poolEvent.Type == event.PoolCreated {
					po = poolEvent.PoolOptions
				}

				log.WithFields(log.Fields{
					"poolEvent.PoolOptions": fmt.Sprintf("%+v", po),
				}).Debug("PoolOptions")
			},
		}

		opt.SetPoolMonitor(m)

		return opt
	}
}
