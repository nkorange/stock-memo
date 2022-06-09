package enter_strategy

import "github.com/nkorange/stock-memo/pkg/stock"

type LastingDownEnterStrategy struct {
	lastingDays int
}

func NewLastingDownEnterStrategy() *LastingDownEnterStrategy {
	return &LastingDownEnterStrategy{
		lastingDays: 3,
	}
}

func (ldes *LastingDownEnterStrategy) Enter(prices []*stock.Price) int {
	cur := 0
	for ind := 0; ind < len(prices)-1; ind++ {
		if prices[ind].ClosePrice > prices[ind+1].ClosePrice {
			cur++
		} else {
			cur = 0
		}
		if cur >= ldes.lastingDays {
			return ind + 1
		}
	}
	return -1
}
