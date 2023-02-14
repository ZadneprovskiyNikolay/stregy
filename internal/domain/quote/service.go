package quote

import (
	"fmt"
	"time"
)

type Service interface {
	Get(symbol string, start, end time.Time, timeframeSec int) (<-chan Quote, Quote)
	Load(symbol, filePath, delimiter string, timeframeSec int) error
}

type service struct {
	repository Repository

	queryRowsLimit int
}

func NewService(repository Repository) Service {
	return &service{repository: repository, queryRowsLimit: 262144}
}

func (s *service) Get(symbol string, start, end time.Time, timeframeSec int) (<-chan Quote, Quote) {
	ch := make(chan Quote, s.queryRowsLimit)
	go quoteGenerator(ch, s, symbol, start, end, timeframeSec)
	return ch, s.firstQuote(symbol, start, end, timeframeSec)
}

func quoteGenerator(ch chan<- Quote, s *service, symbol string, start, end time.Time, timeframeSec int) error {
	batchStart := start
	if err := CheckIsValidTimeframe(timeframeSec); err != nil {
		return err
	}

	var quotes []Quote
	if timeframeSec != 0 {
		quotes = make([]Quote, 0, 262144)
	}

	for batchStart.Before(end) {
		if timeframeSec == 0 {
			err, lastQuoteTime := s.repository.GetAndPushToChan(ch, symbol, batchStart, end, s.queryRowsLimit, timeframeSec)
			if err != nil {
				return err
			}

			batchStart = lastQuoteTime.Add(time.Millisecond * 1)
		} else {
			quotes, err := s.repository.Get(quotes, symbol, batchStart, end, s.queryRowsLimit, timeframeSec)
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

			batchStart = quotes[len(quotes)-1].Time.Add(time.Millisecond * 1)
		}

	}
	close(ch)

	return nil
}

func (s *service) Load(symbol, filePath, delimiter string, timeframeSec int) error {
	return s.repository.Load(symbol, filePath, delimiter, timeframeSec)
}

func (s *service) firstQuote(symbol string, start, end time.Time, timeframeSec int) Quote {
	quotes, _ := s.repository.Get(make([]Quote, 0, 1), symbol, start, end, 1, timeframeSec)
	if len(quotes) == 0 {
		return Quote{}
	}

	return quotes[0]
}
