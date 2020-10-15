package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Account struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	CPF       string             `bson:"cpf"`
	Secret    string             `bson:"secret"`
	Balance   float64            `bson:"balance"`
	CreatedAt time.Time          `bson:"created_at"`
}
