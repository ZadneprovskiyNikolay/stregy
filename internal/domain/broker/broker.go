package broker

import (
	"stregy/internal/domain/order"
	"time"
)

type Broker interface {
	SubmitOrder(o order.Order, ctgOrders ...order.Order) (*order.Order, error)
	AddCtgOrder(posID int64, o order.Order) (*order.Order, error)
	ChangeOrderPrice(id int64, price float64) error
	CancelOrder(id int64) error
	GetBalance() float64

	Time() time.Time
	Price() float64

	Print(s string)
	Printf(format string, v ...interface{})

	Terminate()
}
