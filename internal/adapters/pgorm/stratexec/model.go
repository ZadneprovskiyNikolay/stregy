package stratexec

import (
	"database/sql/driver"
	"stregy/internal/adapters/pgorm/exgaccount"
	btcore "stregy/internal/domain/backtest/core"
	"stregy/internal/domain/symbol"
	"time"

	"github.com/google/uuid"
)

type StrategyExecutionStatus string

func (se *StrategyExecutionStatus) Scan(value interface{}) error {
	*se = StrategyExecutionStatus(value.(string))
	return nil
}

func (se StrategyExecutionStatus) Value() (driver.Value, error) {
	return string(se), nil
}

type StrategyExecution struct {
	StrategyExecutionId uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	StrategyName        string    `gorm:"type:varchar(100);not null"`
	ExchangeAccount     exgaccount.ExchangeAccount
	ExchangeAccountID   *uuid.UUID
	SymbolName          string                  `gorm:"not null"`
	StartTime           time.Time               `gorm:"type:timestamp;not null"`
	EndTime             time.Time               `gorm:"type:timestamp"`
	Status              StrategyExecutionStatus `gorm:"type:strategy_execution_status;not null"`
}

func (s *StrategyExecution) ToBacktest() *btcore.Backtest {
	return &btcore.Backtest{
		ID:           s.StrategyExecutionId.String(),
		StrategyName: s.StrategyName,
		StartTime:    s.StartTime,
		EndTime:      s.EndTime,
		Symbol:       symbol.Symbol{Name: s.SymbolName},
		Status:       btcore.StrategyExecutionStatus(string(s.Status)),
	}
}
