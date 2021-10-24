package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transfer struct {
	ID                   primitive.ObjectID `bson:"_id"`
	OriginAccountID      primitive.ObjectID `bson:"account_origin_id"`
	DestinationAccountID primitive.ObjectID `bson:"account_destination_id"`
	Amount               float64            `bson:"amount"`
	CreatedAt            time.Time          `bson:"created_at"`
}
