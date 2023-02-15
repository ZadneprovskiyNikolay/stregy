package quote

import (
	"time"
)

type Repository interface {
	Get(dest []Quote, symbol string, startTime, endTime time.Time) ([]Quote, error)
	GetFirst(symbol string, startTime, endTime time.Time) (Quote, error)
	GetAndPushToChan(dest chan<- Quote, symbol string, startTime, endTime time.Time) error
	Upload(symbol, filePath, delimiter string, timeframeSec int) error
}
