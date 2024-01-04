package checks

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCheck struct {
	client  *mongo.Client
	Timeout int
}

func NewMongoCheck(timeout int, client *mongo.Client) *MongoCheck {
	return &MongoCheck{
		client:  client,
		Timeout: timeout,
	}
}

func (m *MongoCheck) Pass() bool {
	if m.client == nil {
		return false
	}

	timeout := time.Second * time.Duration(m.Timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := m.client.Ping(ctx, nil)
	if err != nil {
		return false
	}

	return true
}

func (m *MongoCheck) Name() string {
	return "mongodb"
}
