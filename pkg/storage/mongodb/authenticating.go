package mongodb

import (
	"context"
	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (s *Storage) AddToken(ctx context.Context, token authenticating.Token) error {
	collection := s.client.Database(databaseName).Collection(tokensCollection)
	insertionCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	s.log.Infof("Adding token %v to mongodb repo coll %s", token, collection.Name())
	_, err := collection.InsertOne(insertionCtx, token)
	if err != nil {
		s.log.Errorf("Unexpected err %v occurred when adding token %s", err, token.ID.Hex())
	}
	return err
}

func (s *Storage) GetTokenByID(ctx context.Context, id primitive.ObjectID) (authenticating.Token, error) {
	collection := s.client.Database(databaseName).Collection(tokensCollection)
	queryCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	var token authenticating.Token

	s.log.Infof("Retrieving token %v to mongodb repo coll %s", id, collection.Name())
	result := collection.FindOne(queryCtx, bson.D{{Key: "_id", Value: id}})
	if err := result.Decode(&token); err != nil {
		if err == mongo.ErrNoDocuments {
			s.log.Errorf("No token was found for id %s", id)
			return authenticating.Token{}, ErrNoTokenWasFound
		}
		s.log.Errorf("Unexpected err %v when retrieving token %s", err, id)
		return authenticating.Token{}, err
	}
	return token, nil
}
