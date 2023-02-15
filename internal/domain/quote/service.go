package quote

import (
	"time"
)

type Service interface {
	Get(symbol string, start, end time.Time) (<-chan Quote, Quote)
	Upload(symbol, filePath, delimiter string, timeframeSec int) error
}

type service struct {
	repository Repository

	queryRowsLimit int
}

func NewService(repository Repository) Service {
	return &service{repository: repository, queryRowsLimit: 262144}
}

func (s *service) Get(symbol string, start, end time.Time) (<-chan Quote, Quote) {
	ch := make(chan Quote, s.queryRowsLimit)
	go quoteGenerator(ch, s, symbol, start, end)
	return ch, s.firstQuote(symbol, start, end)
}

func quoteGenerator(ch chan<- Quote, s *service, symbol string, start, end time.Time) error {
	err := s.repository.GetAndPushToChan(ch, symbol, start, end)
	if err != nil {
		return err
	}

	close(ch)

	return nil
}

func (s *service) Upload(symbol, filePath, delimiter string, timeframeSec int) error {
	return s.repository.Upload(symbol, filePath, delimiter, timeframeSec)
}

func (s *service) firstQuote(symbol string, start, end time.Time) Quote {
	q, err := s.repository.GetFirst(symbol, start, end)
	if err != nil {
		return Quote{}
	}

	return q
}
