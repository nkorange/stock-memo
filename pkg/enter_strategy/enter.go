package enter_strategy

import "github.com/nkorange/stock-memo/pkg/stock"

type EnterStrategy interface {
	Enter(prices []*stock.Price) int
}
