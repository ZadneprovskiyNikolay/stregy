package composites

import (
	quote1 "stregy/internal/adapters/pgorm/quote"
	"stregy/internal/domain/quote"
)

type QuoteComposite struct {
	Repository quote.Repository
	Service    quote.Service
}

func NewQuoteComposite(composite *PostgresComposite) (*QuoteComposite, error) {
	repository := quote1.NewRepository(composite.dbGORM, composite.dbPQ)
	service := quote.NewService(repository)

	return &QuoteComposite{
		Repository: repository,
		Service:    service,
	}, nil
}
