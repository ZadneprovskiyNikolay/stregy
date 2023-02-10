package backtest

import (
	"stregy/internal/domain/order"
	"stregy/internal/domain/symbol"
)

type AccountHistoryReporter interface {
	CreateReport(orders []*order.Order, s symbol.Symbol, filePath string) error
}
