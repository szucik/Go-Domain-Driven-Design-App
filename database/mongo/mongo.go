package mongo

import (
	"context"
	"fmt"
	"github.com/szucik/trade-helper/portfolio"
	"github.com/szucik/trade-helper/transaction"
	"github.com/szucik/trade-helper/user"
	"github.com/szucik/trade-helper/user/document"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

const uri = "mongodb://127.0.0.1:27017/"

type Repository struct {
	users        *mongo.Collection
	transactions *mongo.Collection
}

type Document struct {
	User       user.User          `bson:"user"`
	Portfolios []portfolio.Entity `bson:"portfolios"`
}

func NewDatabase(ctx context.Context) (Repository, error) {
	db, err := connectWithMongoDB(ctx)
	if err != nil {
		return Repository{}, fmt.Errorf("mongo NewDatabase failed: %w", err)
	}

	return Repository{
		users:        db.Collection("users"),
		transactions: db.Collection("transactions"),
	}, nil
}

func (r Repository) SignUp(ctx context.Context, aggregate user.Aggregate) (string, error) {
	document := document.NewDocument(aggregate)

	_, err := r.users.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}

	return aggregate.User().Email, nil
}

func (r Repository) GetUserByName(ctx context.Context, userName string) (user.Aggregate, error) {
	var result Document
	filter := bson.D{{"user.username", userName}}
	err := r.users.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return user.Aggregate{}, err
	}

	aggregate, err := result.User.NewAggregate()
	if err != nil {
		return user.Aggregate{}, err
	}
	return aggregate, nil
}

func (r Repository) GetUserByEmail(ctx context.Context, email string) (user.Aggregate, error) {
	var result Document
	filter := bson.M{"email": email}
	err := r.users.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return user.Aggregate{}, err
	}

	aggregate, err := result.User.NewAggregate()
	if err != nil {
		return user.Aggregate{}, err
	}
	return aggregate, nil
}

func (r Repository) GetUsers(ctx context.Context) ([]user.Aggregate, error) {
	filter := bson.D{{"users", ""}}
	cursor, err := r.users.Find(ctx, filter)
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
	_, err := r.users.InsertOne(ctx, aggregate)
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
	dbName := os.Getenv("MONGO_DB_NAME")
	db := client.Database(dbName)

	return db, nil
}
