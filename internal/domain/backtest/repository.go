package backtest

import btcore "stregy/internal/domain/backtest/core"

type Repository interface {
	Create(backtest *btcore.Backtest) (*btcore.Backtest, error)
	Save(backtest *btcore.Backtest) (*btcore.Backtest, error)
	GetBacktest(id string) (*btcore.Backtest, error)
}
