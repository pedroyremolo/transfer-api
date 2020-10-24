package adding

import "time"

type Transfer struct {
	OriginAccountID      string
	DestinationAccountID string
	Amount               float64
	CreatedAt            time.Time
}
