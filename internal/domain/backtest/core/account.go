package core

import (
	"stregy/internal/domain/order"
)

func (b *Backtest) GetBalance() float64 {
	return b.balance
}

func (b *Backtest) updateBalance(o *order.Order) {
	if o == nil || o.Status != order.FilledOrder || o.Position.MainOrder.ID == o.ID {
		return
	}

	p := o.Position
	if p.MainOrder.Diraction == order.Long {
		b.balance += o.ExecutionPrice - p.MainOrder.ExecutionPrice
	} else {
		b.balance += p.MainOrder.ExecutionPrice - o.ExecutionPrice
	}
}
