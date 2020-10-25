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
