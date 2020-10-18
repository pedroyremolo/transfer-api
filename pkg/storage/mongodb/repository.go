package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

type Storage struct {
	client *mongo.Client
}

const (
	accountsCollection = "accounts"
)

var ErrCPFAlreadyExists = errors.New("this cpf could not be inserted in our DB")

var (
	databaseName = os.Getenv("APP_DOCUMENT_DB_NAME")
	username     = os.Getenv("APP_DOCUMENT_DB_USERNAME")
	password     = os.Getenv("APP_DOCUMENT_DB_SECRET")
	host         = os.Getenv("APP_DOCUMENT_DB_HOST")
	port         = os.Getenv("APP_DOCUMENT_DB_PORT")
	indexMap     = map[string][]mongo.IndexModel{
		accountsCollection: {
			{
				Keys:    bson.M{"cpf": 1},
				Options: options.Index().SetUnique(true),
			},
		},
	}
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
	disconnectCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()
	if err := s.client.Disconnect(disconnectCtx); err != nil {
		panic(err)
	}
}

func (s *Storage) CreateIndexes(ctx context.Context) {
	db := s.client.Database(databaseName)
	indexCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()
	for collName, indexes := range indexMap {
		_, err := db.Collection(collName).Indexes().CreateMany(indexCtx, indexes)
		if err != nil {
			panic(err)
		}
	}
}

func (s *Storage) AddAccount(ctx context.Context, account adding.Account) (string, error) {
	collection := s.client.Database(databaseName).Collection(accountsCollection)
	insertionCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	dbAccount := Account{
		ID:        primitive.NewObjectID(),
		Name:      string(account.Name),
		CPF:       string(account.CPF),
		Secret:    string(account.Secret),
		Balance:   float64(account.Balance),
		CreatedAt: account.CreatedAt,
	}

	oid, err := collection.InsertOne(insertionCtx, dbAccount)
	if err != nil {
		// TODO Err logging
		return "", ErrCPFAlreadyExists
	}
	return oid.InsertedID.(primitive.ObjectID).Hex(), err
}
