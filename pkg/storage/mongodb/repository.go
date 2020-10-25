package mongodb

import (
	"context"
	"errors"
	"fmt"
	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"github.com/pedroyremolo/transfer-api/pkg/listing"
	"github.com/pedroyremolo/transfer-api/pkg/updating"
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

func (s *Storage) AddTransfer(ctx context.Context, transfer adding.Transfer) (string, error) {
	collection := s.client.Database(databaseName).Collection(transfersCollection)
	insertionCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	dbTransfer := Transfer{
		ID:                   primitive.NewObjectID(),
		OriginAccountID:      transfer.OriginAccountID,
		DestinationAccountID: transfer.DestinationAccountID,
		Amount:               transfer.Amount,
		CreatedAt:            transfer.CreatedAt,
	}

	oid, err := collection.InsertOne(insertionCtx, dbTransfer)
	if err != nil {
		// TODO Err logging
		return "", err
	}
	return oid.InsertedID.(primitive.ObjectID).Hex(), err
}

func (s *Storage) GetAccountByID(ctx context.Context, id string) (listing.Account, error) {
	collection := s.client.Database(databaseName).Collection(accountsCollection)
	queryContext, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var account Account
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return listing.Account{}, ErrNoAccountWasFound
	}
	result := collection.FindOne(queryContext, bson.D{{"_id", oid}})
	if err := result.Decode(&account); err != nil {
		if err == mongo.ErrNoDocuments {
			err = ErrNoAccountWasFound
		}
		return listing.Account{}, err
	}
	return listing.Account{
		ID:        account.ID.Hex(),
		Name:      account.Name,
		CPF:       account.CPF,
		Secret:    account.Secret,
		Balance:   account.Balance,
		CreatedAt: &account.CreatedAt,
	}, nil
}

func (s *Storage) GetAccountByCPF(ctx context.Context, cpf string) (listing.Account, error) {
	collection := s.client.Database(databaseName).Collection(accountsCollection)
	queryContext, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var account Account
	result := collection.FindOne(queryContext, bson.D{{"cpf", cpf}})
	if err := result.Decode(&account); err != nil {
		if err == mongo.ErrNoDocuments {
			err = ErrNoAccountWasFound
		}
		return listing.Account{}, err
	}
	return listing.Account{
		ID:        account.ID.Hex(),
		Name:      account.Name,
		CPF:       account.CPF,
		Secret:    account.Secret,
		Balance:   account.Balance,
		CreatedAt: &account.CreatedAt,
	}, nil
}

func (s *Storage) GetAccounts(ctx context.Context) ([]listing.Account, error) {
	accounts := make([]listing.Account, 0)

	queryContext, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	cursor, err := s.client.Database(databaseName).Collection(accountsCollection).Find(queryContext, bson.D{})
	if err != nil {
		return accounts, err
	}
	defer func() {
		err = cursor.Close(queryContext)
		if err != nil {
			panic(err)
		}
	}()

	for cursor.Next(queryContext) {
		var a Account
		if err = cursor.Decode(&a); err != nil {
			// TODO Log err
			continue
		}

		accounts = append(accounts, listing.Account{
			ID:        a.ID.Hex(),
			Name:      a.Name,
			CPF:       a.CPF,
			Secret:    a.Secret,
			Balance:   a.Balance,
			CreatedAt: &a.CreatedAt,
		})
	}

	return accounts, nil
}

func (s *Storage) AddToken(ctx context.Context, token authenticating.Token) error {
	collection := s.client.Database(databaseName).Collection(tokensCollection)
	insertionCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	_, err := collection.InsertOne(insertionCtx, token)

	return err
}

func (s *Storage) GetTokenByID(ctx context.Context, id primitive.ObjectID) (authenticating.Token, error) {
	collection := s.client.Database(databaseName).Collection(tokensCollection)
	queryCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	var token authenticating.Token

	result := collection.FindOne(queryCtx, bson.D{{Key: "_id", Value: id}})
	if err := result.Decode(&token); err != nil {
		if err == mongo.ErrNoDocuments {
			return authenticating.Token{}, ErrNoTokenWasFound
		}
		return authenticating.Token{}, err
	}
	return token, nil
}

func (s *Storage) UpdateAccounts(ctx context.Context, accounts []updating.Account) error {
	collection := s.client.Database(databaseName).Collection(accountsCollection)

	updatesSessCtx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()
	for _, account := range accounts {
		id, decodeErr := primitive.ObjectIDFromHex(account.ID)
		if decodeErr != nil {
			return decodeErr
		}
		result, updtErr := collection.UpdateOne(
			updatesSessCtx,
			bson.D{{"_id", id}},
			bson.D{{"$set", bson.D{{"balance", account.Balance}}}},
		)
		if updtErr != nil || result.ModifiedCount == 0 {
			return fmt.Errorf("failed to update account %s", id)
		}
	}
	// TODO Realize how to make it work (atomic transaction)
	//updatesCb := func(sessCtx mongo.SessionContext) (interface{}, error) {
	//	for _, account := range accounts {
	//		id, decodeErr := primitive.ObjectIDFromHex(account.ID)
	//		if decodeErr != nil {
	//			_ = sessCtx.AbortTransaction(sessCtx)
	//			return nil, decodeErr
	//		}
	//		result, updtErr := collection.UpdateOne(
	//			sessCtx,
	//			bson.D{{"_id", id}},
	//			bson.D{{"$set", bson.D{{"balance", account.Balance}}}},
	//		)
	//		if updtErr != nil || result.ModifiedCount == 0 {
	//			_ = sessCtx.AbortTransaction(sessCtx)
	//			return nil, updtErr
	//		}
	//	}
	//
	//	return nil, nil
	//}
	//
	//session, err := s.client.StartSession()
	//if err != nil {
	//	return err
	//}
	//defer session.EndSession(updatesSessCtx)
	//
	//_, err = session.WithTransaction(updatesSessCtx, updatesCb)
	//if err != nil {
	//
	//	return err
	//}

	return nil
}
