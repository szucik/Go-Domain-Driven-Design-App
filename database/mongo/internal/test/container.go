package test

import (
	"context"
	"fmt"

	"github.com/ory/dockertest/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const containerName = "mongo_integration_test"

type Mongo struct {
	Client  *mongo.Client
	Cleanup func() error
}

func RunMongoDB(ctx context.Context) (*Mongo, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	if err := pool.RemoveContainerByName(containerName); err != nil {
		return nil, err
	}

	mongoContainer, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "6.0",
		Name:       containerName,
	})

	if err != nil {
		return nil, err
	}

	if err = mongoContainer.Expire(600); err != nil {
		return nil, err
	}

	network, ok := mongoContainer.Container.NetworkSettings.Networks["bridge"]
	if !ok {
		return nil, fmt.Errorf("bridge network does not exist for consul test")
	}

	uri := fmt.Sprintf("mongodb://%s:27017/", network.IPAddress)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err = checkConnection(ctx, pool, client); err != nil {
		mongoContainer.Close()

		return nil, err
	}

	return &Mongo{
		Client:  client,
		Cleanup: mongoContainer.Close,
	}, err
}

func checkConnection(ctx context.Context, pool *dockertest.Pool, client *mongo.Client) error {
	return pool.Retry(func() error {
		return client.Ping(ctx, nil)
	})
}
