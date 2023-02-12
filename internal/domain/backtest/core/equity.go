package core

import "stregy/internal/domain/order"

func (b *Backtest) GetEquity() float64 {
	equity := b.Balance.GetLast()

	for _, p := range b.positions {
		if p.MainOrder.Diraction == order.Long {
			equity += b.price - p.MainOrder.ExecutionPrice
		} else {
			equity += p.MainOrder.ExecutionPrice - b.price
		}
	}

	return equity
}
