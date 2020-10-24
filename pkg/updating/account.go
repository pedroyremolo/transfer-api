package updating

type Account struct {
	ID      string  `bson:"_id"`
	Balance float64 `bson:"balance"`
}
