package backtester

import (
	"context"
)

type Repository interface {
	CreateBacktester(ctx context.Context, backtester Backtester, exchangeAccountID string) (*Backtester, error)
}
