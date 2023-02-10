package composites

import (
	"stregy/internal/adapters/api"
	strategy1 "stregy/internal/adapters/api/strategy"
	strategy3 "stregy/internal/adapters/fs/strategy"
	"stregy/internal/domain/strategy"
	"stregy/internal/domain/user"
)

type StrategyComposite struct {
	Service     strategy.Service
	userService user.Service
	Handler     api.Handler
}

func NewStrategyComposite(userService user.Service) (*StrategyComposite, error) {
	storage := strategy3.NewStorage()
	service := strategy.NewService(storage)
	handler := strategy1.NewHandler(service, userService)
	return &StrategyComposite{
		Service: service,
		Handler: handler,
	}, nil
}
