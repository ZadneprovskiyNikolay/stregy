package acchistory

import "time"

type Drawdown struct {
	TimeSeries       TimeSeries
	MaxEquity        float64
	CalcFreq         time.Duration
	NextCalcTime     time.Time
	EquityCalculator EquityCalculator
}

type EquityCalculator interface {
	GetEquity() float64
}

func (m *Drawdown) Init(t time.Time, calcFreq time.Duration, equityCalculator EquityCalculator) {
	m.EquityCalculator = equityCalculator

	if calcFreq == time.Second*0 {
		m.NextCalcTime = time.Unix(1<<63-62135596801, 999999999) // max time
	} else {
		m.NextCalcTime = t
	}
}

func (m *Drawdown) Update(t time.Time) {
	if t.Before(m.NextCalcTime) {
		return
	}

	equity := m.EquityCalculator.GetEquity()
	if equity > m.MaxEquity {
		m.MaxEquity = equity
	}

	m.TimeSeries = append(m.TimeSeries, TSValue{t, (m.MaxEquity - equity) / m.MaxEquity * 100})
	m.NextCalcTime = m.NextCalcTime.Add(m.CalcFreq)
}

func (m *Drawdown) Save(path string) error {
	return m.TimeSeries.Save(path)
}
