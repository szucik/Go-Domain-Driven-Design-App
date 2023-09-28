package mongo

import (
	"context"
	"fmt"
	"github.com/szucik/trade-helper/user/document"

	"github.com/szucik/trade-helper/transaction"
	"github.com/szucik/trade-helper/user"
	"go.mongodb.org/mongo-driver/bson" // Add this line
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Replace the placeholder with your Atlas connection string
const uri = "mongodb://127.0.0.1:27017"

type Repository struct {
	user         *mongo.Collection
	transactions *mongo.Collection
}

func NewDatabase(ctx context.Context) (Repository, error) {
	db, err := connectWithMongoDB(ctx)
	if err != nil {
		return Repository{}, fmt.Errorf("mongo NewDatabase failed: %w", err)
	}

	return Repository{
		user:         db.Collection("user"),
		transactions: db.Collection("transactions"),
	}, nil
}

func (r Repository) SignUp(ctx context.Context, aggregate user.Aggregate) (string, error) {
	document := document.NewDocument(aggregate)

	one, err := r.user.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}

	fmt.Println(one.InsertedID)
	return aggregate.User().Email, nil
}

func (r Repository) GetUserByEmail(ctx context.Context, email string) (user.Aggregate, error) {
	var result user.Aggregate
	filter := bson.D{{Key: "email", Value: email}}
	err := r.user.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return user.Aggregate{}, fmt.Errorf("GetUserByEmail failed: %w", err)
	}
	return result, nil
}

func (r Repository) GetUsers(ctx context.Context) ([]user.Aggregate, error) {
	cursor, err := r.user.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("GetUsers failed: %w", err)
	}
	var results []user.Aggregate
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("GetUsers failed: %w", err)
	}
	return results, nil
}

func (r Repository) SaveAggregate(ctx context.Context, aggregate user.Aggregate) error {
	_, err := r.user.InsertOne(ctx, aggregate)
	if err != nil {
		return fmt.Errorf("SaveAggregate failed: %w", err)
	}
	return nil
}

func (r Repository) AddTransaction(ctx context.Context, transaction transaction.ValueObject) (string, error) {
	res, err := r.transactions.InsertOne(ctx, transaction)
	if err != nil {
		return "", fmt.Errorf("AddTransaction failed: %w", err)
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// func (r Repository) GetUserByEmail(ctx context.Context, userName string) (user.Aggregate, error) {
// 	//TODO implement me
// 	panic("implement me")
// }

// func (r Repository) GetUsers(ctx context.Context) ([]user.Aggregate, error) {
// 	//TODO implement me
// 	panic("implement me")
// }

// 	fmt.Println(one)
// 	return aggregate.User().Email, nil
// }

// func (r Repository) SaveAggregate(ctx context.Context, aggregate user.Aggregate) error {

// 	panic("implement me")
// }

// func (r Repository) AddTransaction(ctx context.Context, transaction transaction.ValueObject) (string, error) {
// 	//TODO implement me
// 	panic("implement me")
// }

func connectWithMongoDB(ctx context.Context) (*mongo.Database, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("connection with Mongo failed: %w", err)
	}

	// Ping the database to confirm a successful connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("ping to Mongo failed: %w", err)
	}

	fmt.Println("Successfully connected to MongoDB!")

	db := client.Database("database_name") // replace with your database name

	return db, nil
}
