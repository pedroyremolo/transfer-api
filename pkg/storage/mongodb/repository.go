package mongodb

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/pedroyremolo/transfer-api/pkg/log/lgr"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	client *mongo.Client
	log    *logrus.Logger
}

const (
	accountsCollection  = "accounts"
	tokensCollection    = "tokens"
	transfersCollection = "transfers"
)

var ErrCPFAlreadyExists = errors.New("this cpf could not be inserted in our DB")
var ErrNoAccountWasFound = errors.New("no account was found with the given filter parameters")
var ErrNoTokenWasFound = errors.New("no token was found with the given filter parameters")

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
	s.log = lgr.NewDefaultLogger()
	if err != nil {
		s.log.Errorf("Err %v occurred when getting an instance of MongoClient", err)
		return nil, err
	}
	return s, nil
}

func (s *Storage) Connect(ctx context.Context) {
	mongoConnCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	err := s.client.Connect(mongoConnCtx)
	if err != nil {
		s.log.Panicf("Err %v occurred when connecting to mongodb", err)
		panic(err)
	}
}

func (s *Storage) Disconnect(ctx context.Context) {
	disconnectCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()
	if err := s.client.Disconnect(disconnectCtx); err != nil {
		s.log.Panicf("Err %v occurred when disconnecting to mongodb", err)
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
			s.log.Panicf("Err %v occurred when setting indexes to mongodb", err)
			panic(err)
		}
	}
}
