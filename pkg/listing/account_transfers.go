package listing

type AccountTransfers struct {
	Sent     []Transfer `json:"sent"`
	Received []Transfer `json:"received"`
}
