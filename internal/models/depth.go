package models

type Depth struct {
	Timestamp int64 `json:"timestamp"`
	Asks      []Ask `json:"asks"`
	Bids      []Ask `json:"bids"`
}

type Ask struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}
