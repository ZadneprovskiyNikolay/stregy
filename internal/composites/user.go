package composites

import (
	"stregy/internal/adapters/api"
	user1 "stregy/internal/adapters/api/user"
	user2 "stregy/internal/adapters/pgorm/user"
	"stregy/internal/domain/user"
)

type UserComposite struct {
	Repository user.Repository
	Service    user.Service
	Handler    api.Handler
}

func NewUserComposite(composite *PostgresComposite) (*UserComposite, error) {
	repository := user2.NewRepository(composite.dbGORM)
	service := user.NewService(repository)
	handler := user1.NewHandler(service)
	return &UserComposite{
		Repository: repository,
		Service:    service,
		Handler:    handler,
	}, nil
}
