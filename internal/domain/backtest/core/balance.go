package core

import (
	"stregy/internal/domain/order"
)

func (b *Backtest) GetBalance() float64 {
	return b.Balance.GetLast()
}

func (b *Backtest) updateBalance(o *order.Order) {
	if o == nil || o.Status != order.FilledOrder || o.Position.MainOrder.ID == o.ID {
		return
	}

	balance := b.Balance.GetLast()
	commission := b.Commission.GetCommission(o)
	balance -= commission
	b.TotalCommission += commission

	p := o.Position
	if p.MainOrder.Diraction == order.Long {
		balance += o.ExecutionPrice - p.MainOrder.ExecutionPrice
	} else {
		balance += p.MainOrder.ExecutionPrice - o.ExecutionPrice
	}
	b.Balance.Update(balance, b.time)
}
