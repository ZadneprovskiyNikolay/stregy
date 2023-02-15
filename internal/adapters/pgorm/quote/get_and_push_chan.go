package quote

import (
	"stregy/internal/domain/quote"
	"stregy/pkg/logging"
	"time"
)

func (r *repository) GetAndPushToChan(
	dest chan<- quote.Quote,
	symbol string,
	startTime time.Time,
	endTime time.Time,
) error {
	logger = logging.GetLogger()

	dTotal := 5
	downloaders := make([]Downloader, 0, dTotal)

	periodDuration := time.Hour * 24 * 1
	nextPeriod := NewPeriod(startTime, endTime, periodDuration)
	for i := 0; i < dTotal; i++ {
		d := NewDownloader(dest, symbol, r, int(periodDuration.Seconds()))
		d.Period <- nextPeriod
		if nextPeriod.StartTime.After(endTime) {
			dTotal = i
			break
		}

		d.Start()
		downloaders = append(downloaders, d)

		nextPeriod = NewPeriod(nextPeriod.EndTime, endTime, periodDuration)
	}

	dIdx := 0
	downloaders[dIdx].StartPushing <- true
	for {
		n := <-downloaders[dIdx].QuotesPushed

		if n == 0 {
			for i := 0; i < dTotal; i++ {
				downloaders[dIdx].Terminate <- true
				dIdx = nextCircularIndex(dIdx, dTotal)
			}
			break
		}

		downloaders[dIdx].Period <- nextPeriod

		dIdx = nextCircularIndex(dIdx, dTotal)
		downloaders[dIdx].StartPushing <- true
		nextPeriod = NewPeriod(nextPeriod.EndTime, endTime, periodDuration)
	}

	return nil
}
