package mongodb

import (
	"context"
	"fmt"
	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

type Storage struct {
	client *mongo.Client
}

var (
	databaseName = os.Getenv("APP_DOCUMENT_DB_NAME")
	username     = os.Getenv("APP_DOCUMENT_DB_USERNAME")
	password     = os.Getenv("APP_DOCUMENT_DB_SECRET")
	host         = os.Getenv("APP_DOCUMENT_DB_HOST")
	port         = os.Getenv("APP_DOCUMENT_DB_PORT")
)

const (
	accountsCollection = "accounts"
)

func NewStorageFromEnv() (*Storage, error) {
	var err error

	s := new(Storage)

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	clientOptions := options.Client().ApplyURI(uri)
	s.client, err = mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Storage) Connect(ctx context.Context) {
	mongoConnCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	err := s.client.Connect(mongoConnCtx)
	if err != nil {
		panic(err)
	}
}

func (s *Storage) Disconnect(ctx context.Context) {
	if err := s.client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func (s *Storage) AddAccount(ctx context.Context, account adding.Account) (string, error) {
	collection := s.client.Database(databaseName).Collection(accountsCollection)
	insertionCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	dbAccount := Account{
		ID:        primitive.NewObjectID(),
		Name:      account.Name,
		CPF:       account.CPF,
		Secret:    string(account.Secret),
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt,
	}

	oid, err := collection.InsertOne(insertionCtx, dbAccount)
	if err != nil {
		return "", err
	}
	return oid.InsertedID.(primitive.ObjectID).Hex(), err
}
