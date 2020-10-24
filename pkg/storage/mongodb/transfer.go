package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Transfer struct {
	ID                   primitive.ObjectID `bson:"_id"`
	OriginAccountID      string             `bson:"account_origin_id"`
	DestinationAccountID string             `bson:"destination_origin_id"`
	Amount               float64            `bson:"amount"`
	CreatedAt            time.Time          `bson:"created_at"`
}
