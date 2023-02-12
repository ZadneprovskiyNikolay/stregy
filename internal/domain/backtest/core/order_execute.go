package core

import (
	"stregy/internal/domain/order"
	"stregy/pkg/utils"
)

func (b *Backtest) executeOrder(o *order.Order, price float64) {
	isContingent := (o.ID != o.Position.MainOrder.ID)
	if isContingent {
		o.Size = utils.Min(o.Size, o.Position.Size)
		o.Position.Size -= o.Size
	} else {
		o.Position.Size = o.Size
	}

	o.Status = order.FilledOrder
	o.ExecutionPrice = price
	o.FCTime = b.time
	delete(b.orders, o.ID)

	b.logger.LogOrderStatusChange(o)
	b.strategy.OnOrder(*o)

	if !isContingent {
		b.activateContingentOrders(o)
	} else if o.Position.Size == 0 {
		b.cancelContingentOrders(o)
	}

	b.updateBalance(o)
}

func (b *Backtest) activateContingentOrders(o *order.Order) {
	for _, oCtg := range o.Position.CtgOrders {
		b.submitOrder(oCtg)
	}
}
