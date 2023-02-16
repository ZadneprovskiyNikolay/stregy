package core

import (
	"stregy/internal/domain/acchistory"
	"stregy/internal/domain/backtest/commission"
	"stregy/internal/domain/broker"
	"stregy/internal/domain/order"
	"stregy/internal/domain/strategy"
	"stregy/internal/domain/symbol"
	"time"
)

type StrategyExecutionStatus string

const (
	Created    StrategyExecutionStatus = "Created"
	Running    StrategyExecutionStatus = "Running"
	Finished   StrategyExecutionStatus = "Finished"
	Terminated StrategyExecutionStatus = "Terminated"
	Crashed    StrategyExecutionStatus = "Crashed"
)

type Backtest struct {
	ID           string
	StrategyName string
	StartTime    time.Time
	EndTime      time.Time
	Symbol       symbol.Symbol
	Status       StrategyExecutionStatus

	strategy   strategy.Strategy
	Commission commission.CommissionModel

	time  time.Time
	price float64

	OrderHistory        []*order.Order
	Balance             acchistory.Balance
	Drawdown            acchistory.Drawdown
	TotalCommission     float64
	TotalVolumeTraded   float64
	TotalOrdersExecuted int

	orders        map[int64]*order.Order
	positions     map[int64]*order.Position
	orderCount    int64
	positionCount int64

	termChan chan bool

	logger broker.Logger
}

var _ broker.Broker = (*Backtest)(nil)
