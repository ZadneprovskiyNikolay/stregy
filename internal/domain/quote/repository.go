package quote

import (
	"time"
)

type Repository interface {
	Get(dest []Quote, symbol string, startTime, endTime time.Time, limit, timeframeSec int) ([]Quote, error)
	GetAndPushToChan(dest chan<- Quote, symbol string, startTime, endTime time.Time, limit, timeframeSec int) (error, time.Time)
	Load(symbol, filePath, delimiter string, timeframeSec int) error
}
