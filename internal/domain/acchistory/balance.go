package acchistory

import (
	"time"
)

type Balance struct {
	TimeSeries TimeSeries
}

func (b *Balance) Update(balance float64, t time.Time) {
	(*(b)).TimeSeries = append(b.TimeSeries, TSValue{t, balance})
}

func (b Balance) GetLast() float64 {
	return b.TimeSeries[len(b.TimeSeries)-1].Value
}

func (b Balance) Save(path string) error {
	return b.TimeSeries.Save(path)
}
