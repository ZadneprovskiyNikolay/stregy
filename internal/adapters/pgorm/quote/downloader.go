package quote

import (
	"stregy/internal/domain/quote"
	"stregy/pkg/logging"
)

type Downloader struct {
	counter int

	dest        chan<- quote.Quote
	symbol      string
	StartBufCap int

	repository quote.Repository

	Period       chan period
	StartPushing chan bool
	Terminate    chan bool
	QuotesPushed chan int
}

func NewDownloader(dest chan<- quote.Quote, symbol string, repository quote.Repository, startBufCap int) Downloader {
	return Downloader{
		dest:         dest,
		symbol:       symbol,
		StartBufCap:  startBufCap,
		repository:   repository,
		Period:       make(chan period, 1),
		StartPushing: make(chan bool, 1),
		Terminate:    make(chan bool, 1),
		QuotesPushed: make(chan int, 1),
	}
}

func (d *Downloader) Start() {
	go func() {
		buf := make([]quote.Quote, 0, d.StartBufCap)

		for {
			select {
			case <-d.Terminate:
				return
			case period := <-d.Period:
				buf, err := d.repository.Get(buf, d.symbol, period.StartTime, period.EndTime)
				if err != nil {
					logging.GetLogger().Error(err)
				}
				if err != nil || len(buf) == 0 {
					d.QuotesPushed <- 0
					return
				}

				<-d.StartPushing
				for _, q := range buf {
					d.dest <- q
				}
				d.QuotesPushed <- len(buf)
				buf = buf[:0]
			}
		}
	}()
}
