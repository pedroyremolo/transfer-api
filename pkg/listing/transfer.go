package listing

import "time"

type Transfer struct {
	ID                   string    `json:"id"`
	OriginAccountID      string    `json:"account_origin_id"`
	DestinationAccountID string    `json:"account_destination_id"`
	Amount               float64   `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}
