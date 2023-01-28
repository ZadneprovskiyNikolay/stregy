package quote

import (
	"fmt"
	"time"
)

type Service interface {
	Get(symbol string, start, end time.Time, timeframe int) chan Quote
	Load(symbol, filePath, delimiter string, timeframe string) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Get(symbol string, start, end time.Time, timeframeSec int) chan Quote {
	ch := make(chan Quote, 256)
	go quoteGenerator(ch, s, symbol, start, end, timeframeSec)
	return ch
}

func quoteGenerator(ch chan<- Quote, s *service, symbol string, start, end time.Time, timeframeSec int) error {
	batchStart := start
	batchEnd := batchStart.AddDate(0, 0, 1)
	if batchEnd.After(end) {
		batchEnd = end
	}
	if 86400%timeframeSec != 0 {
		return fmt.Errorf("one day is not a multiple of requested timeframe")
	}

	for {
		quotes, err := s.repository.GetByInterval(symbol, batchStart, batchEnd)
		if err != nil {
			return err
		}
		if len(quotes) == 0 {
			break
		}

		quotesAgg, err := AggregateQuotes(quotes, timeframeSec)
		if err != nil {
			panic(fmt.Sprintf("error aggregating quotes: %v\n", err))
		}

		for _, quote := range quotesAgg {
			ch <- quote
		}

		batchStart = batchEnd
		batchEnd = batchStart.AddDate(0, 0, 1)
		if batchEnd.After(end) {
			batchEnd = end
		}
	}
	close(ch)

	return nil
}

func (s *service) Load(symbol, filePath, delimiter string, timeframe string) error {
	return s.repository.Load(symbol, filePath, delimiter, timeframe)
}
