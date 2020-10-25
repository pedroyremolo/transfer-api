package adding

import "time"

type Transfer struct {
	OriginAccountID      string  `json:"account_origin_id"`
	DestinationAccountID string  `json:"account_destination_id"`
	Amount               float64 `json:"amount"`
	CreatedAt            time.Time
}
