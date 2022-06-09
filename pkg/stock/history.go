package stock

type PriceHistory struct {
	Prices []*Price
}

type Price struct {
	Date         string  `json:"Date"`
	ClosePrice   float64 `json:"Price"`
	OpenPrice    float64 `json:"Open"`
	HighestPrice float64 `json:"High"`
	LowestPrice  float64 `json:"Low"`
	TradeAmount  string  `json:"Vol."`
	ChangeRatio  string  `json:"Change %"`
}
