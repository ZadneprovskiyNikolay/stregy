package backtest

import "time"

type BacktestDTO struct {
	StrategyName string
	SymbolName   string
	StartDate    time.Time
	EndDate      time.Time
}
