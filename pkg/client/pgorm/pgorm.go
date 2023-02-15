package pgorm

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second * 5, // Slow SQL threshold
			LogLevel:                  logger.Info,     // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,           // Disable color
		},
	)
}

func NewGormClient(host, port, username, password, database string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", username, password, host, port, database)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
}
