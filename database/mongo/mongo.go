package mongo

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/szucik/trade-helper/portfolio"
	"github.com/szucik/trade-helper/transaction"
	"github.com/szucik/trade-helper/user"
	"github.com/szucik/trade-helper/user/document"
)

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

	r := Repository{
		users:        db.Collection("users"),
		transactions: db.Collection("transactions"),
	}

	if err := createIndexes(ctx, r); err != nil {
		return Repository{}, fmt.Errorf("mongo NewDatabase createIndexes failed: %w", err)
	}

	return r, nil
}

func createIndexes(ctx context.Context, r Repository) error {
	_, err := r.users.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user.email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("unique_email"),
		},
		{
			Keys:    bson.D{{Key: "user.username", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("unique_username"),
		},
	})
	return err
}

func (r Repository) SignUp(ctx context.Context, aggregate user.Aggregate) (string, error) {
	doc := document.NewDocument(aggregate)
	_, err := r.users.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}
	return aggregate.User().Email, nil
}

func (r Repository) GetUserByName(ctx context.Context, userName string) (user.Aggregate, error) {
	var result Document
	filter := bson.D{{Key: "user.username", Value: userName}}
	if err := r.users.FindOne(ctx, filter).Decode(&result); err != nil {
		return user.Aggregate{}, err
	}
	return result.toAggregate()
}

func (r Repository) GetUserByEmail(ctx context.Context, email string) (user.Aggregate, error) {
	var result Document
	filter := bson.D{{Key: "user.email", Value: email}}
	if err := r.users.FindOne(ctx, filter).Decode(&result); err != nil {
		return user.Aggregate{}, err
	}
	return result.toAggregate()
}

func (r Repository) GetUsers(ctx context.Context, p user.PaginationIn) ([]user.Aggregate, error) {
	opts := options.Find()
	if p.Limit > 0 {
		opts.SetSkip(int64(p.Page * p.Limit))
		opts.SetLimit(int64(p.Limit))
	}

	cursor, err := r.users.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("GetUsers failed: %w", err)
	}

	var docs []Document
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("GetUsers failed: %w", err)
	}

	aggregates := make([]user.Aggregate, 0, len(docs))
	for _, doc := range docs {
		a, err := doc.toAggregate()
		if err != nil {
			return nil, fmt.Errorf("GetUsers failed: %w", err)
		}
		aggregates = append(aggregates, a)
	}
	return aggregates, nil
}

func (r Repository) SaveAggregate(ctx context.Context, aggregate user.Aggregate) error {
	doc := document.NewDocument(aggregate)
	filter := bson.D{{Key: "user.email", Value: aggregate.User().Email}}
	_, err := r.users.ReplaceOne(ctx, filter, doc)
	if err != nil {
		return fmt.Errorf("SaveAggregate failed: %w", err)
	}
	return nil
}

func (r Repository) UpdateUser(ctx context.Context, currentUsername string, aggregate user.Aggregate) error {
	doc := document.NewDocument(aggregate)
	filter := bson.D{{Key: "user.username", Value: currentUsername}}
	_, err := r.users.ReplaceOne(ctx, filter, doc)
	if err != nil {
		return fmt.Errorf("UpdateUser failed: %w", err)
	}
	return nil
}

func (r Repository) DeleteUser(ctx context.Context, username string) error {
	filter := bson.D{{Key: "user.username", Value: username}}
	res, err := r.users.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("DeleteUser failed: %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("DeleteUser: user %q not found", username)
	}
	return nil
}

func (r Repository) AddTransaction(ctx context.Context, vo transaction.ValueObject) (string, error) {
	t := vo.Transaction()
	res, err := r.transactions.InsertOne(ctx, t)
	if err != nil {
		return "", fmt.Errorf("AddTransaction failed: %w", err)
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r Repository) GetTransactions(ctx context.Context, username, portfolioName string) ([]transaction.ValueObject, error) {
	filter := bson.D{
		{Key: "username", Value: username},
		{Key: "portfolio_name", Value: portfolioName},
	}
	cursor, err := r.transactions.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("GetTransactions failed: %w", err)
	}

	var txns []transaction.Transaction
	if err = cursor.All(ctx, &txns); err != nil {
		return nil, fmt.Errorf("GetTransactions decode failed: %w", err)
	}

	result := make([]transaction.ValueObject, 0, len(txns))
	for _, t := range txns {
		vo, err := t.NewTransaction()
		if err != nil {
			return nil, err
		}
		result = append(result, vo)
	}
	return result, nil
}

func (d Document) toAggregate() (user.Aggregate, error) {
	aggregate, err := d.User.NewAggregate()
	if err != nil {
		return user.Aggregate{}, err
	}
	for _, p := range d.Portfolios {
		if err := aggregate.AddPortfolio(p); err != nil {
			return user.Aggregate{}, err
		}
	}
	return aggregate, nil
}

func connectWithMongoDB(ctx context.Context) (*mongo.Database, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://127.0.0.1:27017/"
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("connection with Mongo failed: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping to Mongo failed: %w", err)
	}

	fmt.Println("Successfully connected to MongoDB!")
	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "tradehelper"
	}
	return client.Database(dbName), nil
}
