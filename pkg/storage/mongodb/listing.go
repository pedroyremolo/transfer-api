package mongodb

import (
	"context"
	"github.com/pedroyremolo/transfer-api/pkg/listing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (s *Storage) GetAccountByID(ctx context.Context, id string) (listing.Account, error) {
	collection := s.client.Database(databaseName).Collection(accountsCollection)
	queryContext, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	s.log.Infof("Retrieving account %v to mongodb repo coll %s", id, collection.Name())
	var account Account
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.log.Errorf("Err when serializing id %s to ObjectID", id)
		return listing.Account{}, ErrNoAccountWasFound
	}
	result := collection.FindOne(queryContext, bson.D{{"_id", oid}})
	if err := result.Decode(&account); err != nil {
		if err == mongo.ErrNoDocuments {
			s.log.Errorf("No account was found with id %s", id)
			return listing.Account{}, ErrNoAccountWasFound
		}
		s.log.Errorf("Unexpected err %v when retrieving account %s", err, id)
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
	s.log.Infof("Retrieving account of cpf %v to mongodb repo coll %s", cpf, collection.Name())
	result := collection.FindOne(queryContext, bson.D{{"cpf", cpf}})
	if err := result.Decode(&account); err != nil {
		if err == mongo.ErrNoDocuments {
			s.log.Errorf("No account was found with cpf %s", cpf)
			return listing.Account{}, ErrNoAccountWasFound
		}
		s.log.Errorf("Unexpected err %v when retrieving account of cpf %s", err, cpf)
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

	s.log.Infof("Retrieving all accounts of mongodb repo coll %s", accountsCollection)
	cursor, err := s.client.Database(databaseName).Collection(accountsCollection).Find(queryContext, bson.D{})
	if err != nil {
		s.log.Errorf("Unexpected err %v when retrieving accounts", err)
		return accounts, err
	}
	defer func() {
		err = cursor.Close(queryContext)
		if err != nil {
			s.log.Panicf("Err %v occurred when closing cursor", err)
			panic(err)
		}
	}()

	for cursor.Next(queryContext) {
		var a Account
		if err = cursor.Decode(&a); err != nil {
			s.log.Errorf("Err %v occurred when decoding account from mongo repo", err)
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

func (s *Storage) GetTransfersByKey(ctx context.Context, transferKey string, transferValue string) ([]listing.Transfer, error) {
	queryContext, cancel := context.WithTimeout(ctx, time.Second*100)
	defer cancel()

	s.log.Infof("Retrieving transfers by %s with transferValue %s of mongodb repo coll %s", transferKey, transferValue, transfersCollection)
	transfers := make([]listing.Transfer, 0)
	oid, _ := primitive.ObjectIDFromHex(transferValue)
	cur, err := s.client.Database(databaseName).Collection(transfersCollection).Find(ctx, bson.D{{Key: transferKey, Value: oid}})
	func() {
		err = cur.Close(queryContext)
		if err != nil {
			s.log.Panicf("Err %v occurred when closing cursor", err)
			panic(err)
		}
	}()
	if err != nil {
		s.log.Errorf("Err %v occurred when retrieving transfers by %s with transferValue %s", err, transferKey, transferValue)
		return transfers, err
	}
	for cur.Next(queryContext) {
		var t Transfer
		if err = cur.Decode(&t); err != nil {
			s.log.Errorf("Err %v occurred when decoding transfer from mongo repo", err)
			continue
		}

		transfers = append(transfers, listing.Transfer{
			ID:                   t.ID.Hex(),
			OriginAccountID:      t.OriginAccountID.Hex(),
			DestinationAccountID: t.DestinationAccountID.Hex(),
			Amount:               t.Amount,
			CreatedAt:            t.CreatedAt,
		})
	}
	return transfers, nil
}
