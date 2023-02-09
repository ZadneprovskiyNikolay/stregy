package bt

import (
	"stregy/internal/domain/order"
)

func (b *Backtester) SubmitOrder(o order.Order, ctgOrders ...order.Order) (*order.Order, error) {
	err := checkIsValidOrder(&o)
	if err != nil {
		return nil, err
	}

	for _, oCtg := range ctgOrders {
		err = checkIsValidCtgOrder(&oCtg, &o)
		if err != nil {
			o.Status = order.Rejected
			return nil, err
		}
	}

	o.ID = b.newOrderID()
	p := &order.Position{ID: b.newPositionID(), MainOrder: &o}
	o.Position = p
	b.submitOrder(&o)

	for _, oCtg := range ctgOrders {
		oCtg.ID = b.newOrderID()
		oCtg.Status = order.Inactive
		oCtg.Position = p
		o.Position.CtgOrders = append(o.Position.CtgOrders, oCtg.Copy())
	}

	return (&o).Copy(), nil
}

func (b *Backtester) AddCtgOrder(posID int64, o order.Order) (*order.Order, error) {
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

func (b *Backtester) submitOrder(o *order.Order) {
	o.Status = order.Submitted
	o.SubmissionTime = b.curTime

	b.orders[o.ID] = o
	b.orderHistory = append(b.orderHistory, o)

	b.logger.LogOrderStatusChange(o)
}

func (b *Backtester) newOrderID() int64 {
	id := b.orderCount
	b.orderCount += 1
	return id
}

func (b *Backtester) newPositionID() int64 {
	id := b.positionCount
	b.positionCount += 1
	return id
}
