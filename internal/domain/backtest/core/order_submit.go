package core

import (
	"stregy/internal/domain/order"
)

func (b *Backtest) SubmitOrder(o order.Order, ctgOrders ...order.Order) (*order.Order, error) {
	err := checkIsValidOrder(&o)
	if err != nil {
		return nil, err
	}

	for _, oCtg := range ctgOrders {
		err = checkIsValidCtgOrder(&oCtg, &o)
		if err != nil {
			o.Status = order.RejectedOrder
			return nil, err
		}
	}

	o.ID = b.newOrderID()
	p := &order.Position{ID: b.newPositionID(), MainOrder: &o}
	o.Position = p
	b.submitOrder(&o)

	for _, oCtg := range ctgOrders {
		oCtg.ID = b.newOrderID()
		oCtg.Status = order.InactiveOrder
		oCtg.Position = p
		o.Position.CtgOrders = append(o.Position.CtgOrders, oCtg.Copy())
	}

	return (&o).Copy(), nil
}

func (b *Backtest) AddCtgOrder(posID int64, o order.Order) (*order.Order, error) {
	p, ok := b.positions[posID]
	if !ok {
		return nil, &PositionNotFoundError{posID}
	}
	if err := checkIsValidCtgOrder(&o, p.MainOrder); err != nil {
		return nil, err
	}

	p.CtgOrders = append(p.CtgOrders, o.Copy())
	b.submitOrder(&o)

	return (&o).Copy(), nil
}

func (b *Backtest) submitOrder(o *order.Order) {
	o.Status = order.SubmittedOrder
	o.SubmissionTime = b.time

	b.orders[o.ID] = o
	b.OrderHistory = append(b.OrderHistory, o)

	b.logger.LogOrderStatusChange(o)
}

func (b *Backtest) newOrderID() int64 {
	id := b.orderCount
	b.orderCount += 1
	return id
}

func (b *Backtest) newPositionID() int64 {
	id := b.positionCount
	b.positionCount += 1
	return id
}
