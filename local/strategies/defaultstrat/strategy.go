package defaultstrat

import (
	"stregy/internal/domain/order"
	"stregy/internal/domain/quote"
	"stregy/internal/domain/strategy"
	"time"
)

type Strategy struct {
}

func NewStrategy() *Strategy {
	return &Strategy{}
}

func (*Strategy) Name() string {
	return "defaultstrat"
}

func (*Strategy) OnOrder(o order.Order) {
}

func (*Strategy) OnQuote(q quote.Quote, timeframe int) {
}

func (*Strategy) OnTick(price float64) {
}

func (*Strategy) PrimaryTimeframeSec() int {
	return 60
}

func (*Strategy) QuoteTimeframesNeeded() []int {
	return []int{}
}

func (*Strategy) TimeBeforeCallbacks() time.Duration {
	return time.Minute * 5 * 14
}

var _ strategy.Strategy = (*Strategy)(nil)
