package strat1

import (
	"stregy/internal/domain/broker"
	"stregy/internal/domain/order"
	"stregy/internal/domain/strategy"
	"stregy/local/strategies/posmerge"
	"time"
)

type strat1 struct {
	broker         broker.Broker
	positionMerger posmerge.PositionMerger

	newMainOrder *order.Order
	newSlPrice   float64
	newTpPrice   float64

	totalAverages      int
	nextAveragingPrice float64
	mergedPosition     *posmerge.MergedPosition

	prevClose      float64
	prevClosesUp   int
	prevClosesDown int
}

func NewStrategy(broker broker.Broker) strategy.Strategy {
	positionMerger := posmerge.NewPositionMerger(broker)
	strat := strat1{
		broker:         broker,
		positionMerger: positionMerger}
	positionMerger.SetCallbackObject(&strat)

	return &strat
}

func (*strat1) Name() string {
	return "strat1"
}

func (*strat1) PrimaryTimeframeSec() int {
	return 1
}

func (*strat1) QuoteTimeframesNeeded() []int {
	return []int{5}
}

func (*strat1) TimeBeforeCallbacks() time.Duration {
	return time.Minute * 0
}

var _ strategy.Strategy = (*strat1)(nil)
