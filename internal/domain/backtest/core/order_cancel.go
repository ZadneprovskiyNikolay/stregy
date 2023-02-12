package core

import (
	"stregy/internal/domain/order"
)

func (b *Backtest) CancelOrder(id int64) error {
	o, ok := b.orders[id]
	if !ok {
		return &OrderNotFoundError{id}
	}

	delete(b.orders, id)

	o.Status = order.CancelledOrder
	o.FCTime = b.time
	if o.ID != o.Position.MainOrder.ID {
		o.Position.RemoveCgtOrder(id)
	}

	b.logger.LogOrderStatusChange(o)
	b.strategy.OnOrder(*o)

	return nil
}

func (b *Backtest) cancelContingentOrders(o *order.Order) {
	for _, oCtg := range o.Position.CtgOrders {
		b.CancelOrder(oCtg.ID)
	}
}
