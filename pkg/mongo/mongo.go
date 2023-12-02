package mongo

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	serviceName = "mongo_%s"
)

type Client interface {
	GetConnection(ctx context.Context) *mongo.Client
}

type Mongo struct {
	address    string
	username   string
	password   string
	ssl        bool
	dbConnOnce sync.Once
	db         *mongo.Client
}

func New(
	address string,
	username string,
	password string,
	ssl bool,
) Client {
	return &Mongo{
		address:  address,
		username: username,
		password: password,
		ssl:      ssl,
	}
}

func (m *Mongo) GetConnection(ctx context.Context) *mongo.Client {
	if m.db == nil {
		m.dbConnOnce.Do(func() {
			var err error

			url := fmt.Sprintf("mongodb://%s/?", m.address)
			if m.ssl == false {
				url += "&ssl=false"
			}

			clientOptions := options.Client().ApplyURI(url)
			clientOptions.SetAuth(options.Credential{
				Username: m.username,
				Password: m.password,
			})

			m.db, err = mongo.Connect(ctx, clientOptions)
			if err != nil {
				panic(err)
			}

		})
	}
	return m.db
}

func (m *Mongo) ServiceName() string {
	return fmt.Sprintf(serviceName, m.address)
}

func (m *Mongo) IsHealthy(ctx context.Context) bool {
	err := m.GetConnection(ctx).Ping(ctx, nil)
	return err == nil
}
