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
	users        *mongo.Collection
	transactions *mongo.Collection
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
	doc := document.NewDocument(aggregate)
	filter := bson.M{"email": doc.User.Email}
	update := bson.M{"$setOnInsert": doc}

	updateOptions := options.Update().SetUpsert(true)
	_, err := r.users.UpdateOne(ctx, filter, update, updateOptions)
	if err != nil {
		return "", err
	}

	return aggregate.User().Email, nil
}

func (r Repository) GetUserByName(ctx context.Context, userName string) (user.Aggregate, error) {
	var doc document.Document
	filter := bson.M{"username": userName}
	err := r.users.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return user.Aggregate{}, fmt.Errorf("GetUserByName failed: %w", err)
	}
	aggregate, err := doc.NewAggregate()
	if err != nil {
		return user.Aggregate{}, err
	}
	//Todo - create aggregate from bson
	return aggregate, nil
}

func (r Repository) GetUserByEmail(ctx context.Context, email string) (user.Aggregate, error) {
	var result user.User
	filter := bson.M{"email": email}
	err := r.users.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return user.Aggregate{}, fmt.Errorf("GetUserByEmail failed: %w", err)
	}
	//Todo - create aggregate  from bson
	return user.Aggregate{}, nil
}

func (r Repository) GetUsers(ctx context.Context) ([]user.Aggregate, error) {
	cursor, err := r.users.Find(ctx, bson.D{})
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

	db := client.Database("tradehelper") // replace with your database name

	return db, nil
}
