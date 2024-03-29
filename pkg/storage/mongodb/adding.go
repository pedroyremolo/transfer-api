package mongodb

import (
	"context"
	"time"

	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Storage) AddAccount(ctx context.Context, account adding.Account) (string, error) {
	collection := s.client.Database(databaseName).Collection(accountsCollection)
	insertionCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	s.log.Infof("Adding account %v to mongodb repo coll %s", account, collection.Name())
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
		s.log.Errorf("CPF %s already exists in our repo", dbAccount.CPF)
		return "", ErrCPFAlreadyExists
	}
	return oid.InsertedID.(primitive.ObjectID).Hex(), err
}

func (s *Storage) AddTransfer(ctx context.Context, transfer adding.Transfer) (string, error) {
	collection := s.client.Database(databaseName).Collection(transfersCollection)
	insertionCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	s.log.Infof("Adding transfer %v to mongodb repo coll %s", transfer, collection.Name())
	originOID, _ := primitive.ObjectIDFromHex(transfer.OriginAccountID)
	transferOID, _ := primitive.ObjectIDFromHex(transfer.DestinationAccountID)
	dbTransfer := Transfer{
		ID:                   primitive.NewObjectID(),
		OriginAccountID:      originOID,
		DestinationAccountID: transferOID,
		Amount:               transfer.Amount,
		CreatedAt:            transfer.CreatedAt,
	}

	oid, err := collection.InsertOne(insertionCtx, dbTransfer)
	if err != nil {
		s.log.Errorf("Unexpected err when adding transfer %s of origin account %s", dbTransfer.ID, dbTransfer.OriginAccountID)
		return "", err
	}
	return oid.InsertedID.(primitive.ObjectID).Hex(), err
}
