package quote

import (
	"database/sql"
	"stregy/internal/domain/quote"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type repository struct {
	dbPQ   *sql.DB
	dbGORM *gorm.DB
}

func NewRepository(dbGorm *gorm.DB, dbPQ *sql.DB) quote.Repository {
	return &repository{dbPQ: dbPQ, dbGORM: dbGorm}
}
