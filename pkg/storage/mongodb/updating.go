package mongodb

import (
	"context"
	"fmt"
	"github.com/pedroyremolo/transfer-api/pkg/updating"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (s *Storage) UpdateAccounts(ctx context.Context, accounts []updating.Account) error {
	collection := s.client.Database(databaseName).Collection(accountsCollection)

	s.log.Infof("Updating %v accounts of mongo repo coll %s", len(accounts), collection.Name())
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
			s.log.Errorf("Failed to update account %s", account.ID)
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
