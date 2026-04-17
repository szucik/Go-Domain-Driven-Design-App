package mongo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/szucik/trade-helper/apperrors"
	"github.com/szucik/trade-helper/transaction"
	"github.com/szucik/trade-helper/user"
	"github.com/szucik/trade-helper/user/document"
)

type Repository struct {
	users        *mongo.Collection
	transactions *mongo.Collection
}

type Document struct {
	User       user.User                    `bson:"user"`
	Portfolios []document.PortfolioDocument `bson:"portfolios"`
}

type transactionDocument struct {
	ID            uuid.UUID                   `bson:"id"`
	UserName      string                      `bson:"username"`
	PortfolioName string                      `bson:"portfolio_name"`
	Symbol        string                      `bson:"symbol"`
	Type          transaction.TransactionType `bson:"type"`
	Quantity      string                      `bson:"quantity"`
	Price         string                      `bson:"price"`
	Created       time.Time                   `bson:"created"`
}

func toTransactionDocument(t transaction.Transaction) transactionDocument {
	return transactionDocument{
		ID:            t.ID,
		UserName:      t.UserName,
		PortfolioName: t.PortfolioName,
		Symbol:        t.Symbol,
		Type:          t.Type,
		Quantity:      t.Quantity.String(),
		Price:         t.Price.String(),
		Created:       t.Created,
	}
}

func (d transactionDocument) toValueObject() (transaction.ValueObject, error) {
	qty, err := decimal.NewFromString(d.Quantity)
	if err != nil {
		return transaction.ValueObject{}, fmt.Errorf("invalid quantity %q: %w", d.Quantity, err)
	}
	price, err := decimal.NewFromString(d.Price)
	if err != nil {
		return transaction.ValueObject{}, fmt.Errorf("invalid price %q: %w", d.Price, err)
	}
	return transaction.Transaction{
		ID:            d.ID,
		UserName:      d.UserName,
		PortfolioName: d.PortfolioName,
		Symbol:        d.Symbol,
		Type:          d.Type,
		Quantity:      qty,
		Price:         price,
		Created:       d.Created,
	}.NewTransaction()
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
		if mongo.IsDuplicateKeyError(err) {
			return "", apperrors.Error("registration failed", "Conflict", http.StatusConflict)
		}
		return "", err
	}
	return aggregate.User().Email, nil
}

func (r Repository) GetUserByName(ctx context.Context, userName string) (user.Aggregate, error) {
	var result Document
	filter := bson.D{{Key: "user.username", Value: userName}}
	if err := r.users.FindOne(ctx, filter).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user.Aggregate{}, apperrors.Error("user not found", "NotFound", http.StatusNotFound)
		}
		return user.Aggregate{}, err
	}
	return result.toAggregate()
}

func (r Repository) GetUserByEmail(ctx context.Context, email string) (user.Aggregate, error) {
	var result Document
	filter := bson.D{{Key: "user.email", Value: email}}
	if err := r.users.FindOne(ctx, filter).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user.Aggregate{}, apperrors.Error("user not found", "NotFound", http.StatusNotFound)
		}
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
		return apperrors.Error("user not found", "NotFound", http.StatusNotFound)
	}
	return nil
}

func (r Repository) AddTransaction(ctx context.Context, vo transaction.ValueObject) (string, error) {
	doc := toTransactionDocument(vo.Transaction())
	res, err := r.transactions.InsertOne(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("AddTransaction failed: %w", err)
	}
	return res.InsertedID.(bson.ObjectID).Hex(), nil
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

	var docs []transactionDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("GetTransactions decode failed: %w", err)
	}

	result := make([]transaction.ValueObject, 0, len(docs))
	for _, d := range docs {
		vo, err := d.toValueObject()
		if err != nil {
			return nil, err
		}
		result = append(result, vo)
	}
	return result, nil
}

func (d Document) toAggregate() (user.Aggregate, error) {
	doc := document.Document{
		User:       d.User,
		Portfolios: d.Portfolios,
	}
	return doc.NewAggregate()
}

func connectWithMongoDB(ctx context.Context) (*mongo.Database, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://127.0.0.1:27017/"
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)
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
