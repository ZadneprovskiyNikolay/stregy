package composites

import (
	"stregy/internal/adapters/api"
	btapi "stregy/internal/adapters/api/bt"
	"stregy/internal/adapters/pgorm/stratexec"
	"stregy/internal/domain/backtest"
	"stregy/internal/domain/exgaccount"
	"stregy/internal/domain/quote"
	"stregy/internal/domain/symbol"
	"stregy/internal/domain/tick"
	"stregy/internal/domain/user"
)

type BacktestComposite struct {
	Service backtest.Service
	Handler api.Handler
}

func NewBacktestComposite(
	pgormComposite *PostgresComposite,
	exgAccService exgaccount.Service,
	userService user.Service,
	tickService tick.Service,
	quoteService quote.Service,
	symbolService symbol.Service,
) (*BacktestComposite, error) {
	repository := stratexec.NewRepository(pgormComposite.dbGORM)
	service := backtest.NewService(
		repository,
		tickService,
		quoteService,
		exgAccService,
		symbolService)
	handler := btapi.NewHandler(service, userService)
	return &BacktestComposite{
		Service: service,
		Handler: handler,
	}, nil
}
