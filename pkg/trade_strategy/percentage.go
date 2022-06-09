package trade_strategy

import (
	"github.com/nkorange/stock-memo/pkg/enter_strategy"
	"github.com/nkorange/stock-memo/pkg/stock"
)

type PercentageTradeStrategy struct {
	enterStrategy enter_strategy.EnterStrategy
}

func NewPercentageStrategy() *PercentageTradeStrategy {
	return &PercentageTradeStrategy{
		enterStrategy: enter_strategy.NewLastingDownEnterStrategy(),
	}
}

func (ps *PercentageTradeStrategy) Trade(history *stock.PriceHistory, startMoney float64, startDay int) float64 {
	var boughtCount int
	var moneyHold = startMoney
	if startDay >= len(history.Prices) {
		return startMoney
	}
	prices := history.Prices[startDay-1:]
	var buyPrice float64
	for ind := 0; ind < len(prices); {
		price := prices[ind]
		// Buy:
		if price.ClosePrice > 1.1*buyPrice {
			moneyToBuy := startMoney * 0.2
			buyCount := int(moneyToBuy / price.ClosePrice)
			buyMoney := float64(buyCount) * price.ClosePrice
			if buyMoney > moneyHold {
				// Money not enough:
				ind++
				continue
			}
			boughtCount = boughtCount + buyCount
			moneyHold = moneyHold - buyMoney
			buyPrice = price.ClosePrice
			ind++
			continue
		}

		// Sell:
		if price.ClosePrice < 0.9*buyPrice {
			moneyHold = moneyHold + float64(boughtCount)*price.ClosePrice
			boughtCount = 0
			buyPrice = 0
			nextBuyDay := ps.enterStrategy.Enter(prices[ind:])
			if nextBuyDay == -1 {
				break
			}
			ind = ind + nextBuyDay
			continue
		}

		// Neither buy nor sell:
		ind++
	}

	// Return current profit:
	return moneyHold + float64(boughtCount)*prices[len(prices)-1].ClosePrice - startMoney
}
