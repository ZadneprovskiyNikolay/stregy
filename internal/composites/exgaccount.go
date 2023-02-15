package composites

import (
	"stregy/internal/adapters/api"
	exgaccount1 "stregy/internal/adapters/api/exgaccount"
	exgaccount2 "stregy/internal/adapters/pgorm/exgaccount"
	"stregy/internal/domain/exgaccount"
	"stregy/internal/domain/user"
)

type ExchangeAccountComposite struct {
	Repository exgaccount.Repository
	Service    exgaccount.Service
	Handler    api.Handler
}

func NewExchangeAccountComposite(composite *PostgresComposite, userService user.Service) (*ExchangeAccountComposite, error) {
	repository := exgaccount2.NewRepository(composite.dbGORM)
	service := exgaccount.NewService(repository)
	handler := exgaccount1.NewHandler(service, userService)
	return &ExchangeAccountComposite{
		Repository: repository,
		Service:    service,
		Handler:    handler,
	}, nil
}
