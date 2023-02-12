package acchistory

import (
	"stregy/pkg/timeseries"
	"time"
)

type Drawdown struct {
	TimeSeries timeseries.TimeSeries
	CalcFreq   time.Duration

	maxEquity        float64
	nextCalcTime     time.Time
	equityCalculator EquityCalculator
}

type EquityCalculator interface {
	GetEquity() float64
}

func (m *Drawdown) Init(t time.Time, calcFreq time.Duration, equityCalculator EquityCalculator) {
	m.equityCalculator = equityCalculator

	if calcFreq == time.Second*0 {
		m.nextCalcTime = time.Unix(1<<63-62135596801, 999999999) // max time
	} else {
		m.nextCalcTime = t
	}
}

func (m *Drawdown) Update(t time.Time) {
	if t.Before(m.nextCalcTime) {
		return
	}

	equity := m.equityCalculator.GetEquity()
	if equity > m.maxEquity {
		m.maxEquity = equity
	}

	m.TimeSeries = append(m.TimeSeries, timeseries.Value{t, (m.maxEquity - equity) / m.maxEquity * 100})
	m.nextCalcTime = m.nextCalcTime.Add(m.CalcFreq)
}

func (m *Drawdown) Save(path string) error {
	return m.TimeSeries.Save(path)
}
