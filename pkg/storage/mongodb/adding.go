package mongodb

import (
	"context"
	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

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
