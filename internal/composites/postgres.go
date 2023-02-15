package composites

import (
	"database/sql"
	"stregy/internal/adapters/pgorm/migration"
	"stregy/pkg/client/pgorm"
	"stregy/pkg/client/pq"

	"gorm.io/gorm"
)

type PostgresComposite struct {
	dbGORM *gorm.DB
	dbPQ   *sql.DB
}

func NewPGormComposite(host, port, username, password, database string) (*PostgresComposite, error) {
	dbGorm, err := pgorm.NewGormClient(host, port, username, password, database)
	if err != nil {
		return nil, err
	}
	err = migration.Migrate(dbGorm)
	if err != nil {
		return nil, err
	}
	dbPQ, err := pq.NewPqClient(username, password, database)
	if err != nil {
		return nil, err
	}
	return &PostgresComposite{dbGORM: dbGorm, dbPQ: dbPQ}, nil
}
