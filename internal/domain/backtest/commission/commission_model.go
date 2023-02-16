package commission

import "stregy/internal/domain/order"

type CommissionModel interface {
	GetCommission(o *order.Order) float64
}
