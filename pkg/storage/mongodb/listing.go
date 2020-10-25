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

func (s *Storage) GetTransfersByKey(ctx context.Context, key string, value string) ([]listing.Transfer, error) {
	queryContext, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	transfers := make([]listing.Transfer, 0)
	cur, err := s.getDocumentsByKey(queryContext, transfersCollection, key, value)
	func() {
		err = cur.Close(queryContext)
		if err != nil {
			panic(err)
		}
	}()
	if err != nil {
		return transfers, err
	}
	for cur.Next(queryContext) {
		var t Transfer
		if err = cur.Decode(&t); err != nil {
			//TODO Log err
			continue
		}

		transfers = append(transfers, listing.Transfer{
			ID:                   t.ID.Hex(),
			OriginAccountID:      t.OriginAccountID,
			DestinationAccountID: t.DestinationAccountID,
			Amount:               t.Amount,
			CreatedAt:            t.CreatedAt,
		})
	}
	return transfers, nil
}

func (s *Storage) getDocumentsByKey(ctx context.Context, collName string, key string, value string) (*mongo.Cursor, error) {
	cursor, err := s.client.Database(databaseName).Collection(collName).Find(ctx, bson.D{{key, value}})
	if err != nil {
		curErr := cursor.Close(ctx)
		if curErr != nil {
			panic(curErr)
		}
		return nil, err
	}

	return cursor, nil
}
