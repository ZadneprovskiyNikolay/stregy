package quote

import (
	"time"
)

type Quote struct {
	Time  time.Time `gorm:"primaryKey;type:timestamp"`
	OHLCV []byte    `gorm:"type:bytea"`
}

type OHLCV struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int32
}

type period struct {
	StartTime time.Time
	EndTime   time.Time
}

func NewPeriod(start, limit time.Time, duration time.Duration) period {
	var p period
	p.StartTime = start
	p.EndTime = start.Add(duration)
	if p.EndTime.After(limit) {
		p.EndTime = limit
	}

	return p
}
