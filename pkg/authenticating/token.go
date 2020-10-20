package authenticating

import "go.mongodb.org/mongo-driver/bson/primitive"

type Token struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	ClientID string             `json:"client_id,omitempty" bson:"client_id"`
	Digest   string             `json:"token,omitempty" bson:"digest"`
}
