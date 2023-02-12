package acchistory

import (
	"stregy/pkg/timeseries"
	"time"
)

type Balance struct {
	TimeSeries timeseries.TimeSeries
}

func (b *Balance) Update(balance float64, t time.Time) {
	(*(b)).TimeSeries = append(b.TimeSeries, timeseries.Value{t, balance})
}

func (b Balance) GetLast() float64 {
	return b.TimeSeries[len(b.TimeSeries)-1].Value
}

func (b Balance) Save(path string) error {
	return b.TimeSeries.Save(path)
}
