package trade_strategy

import "github.com/nkorange/stock-memo/pkg/stock"

type TradeStrategy interface {
	Trade(history *stock.PriceHistory, startMoney float64, startDay int) float64
}
