package composites

import (
	symbol1 "stregy/internal/adapters/pgorm/symbol"
	"stregy/internal/domain/symbol"
)

type SymbolComposite struct {
	Service symbol.Service
}

func NewSymbolComposite(composite *PostgresComposite) (*SymbolComposite, error) {
	repository := symbol1.NewRepository(composite.dbGORM)

	return &SymbolComposite{
		Service: symbol.NewService(repository),
	}, nil
}
