package commission

import (
	"flag"
	"fmt"
	"stregy/internal/domain/order"
)

type BinanceCommission struct {
	Maker float64
	Taker float64
}

var BinanceMakerCommission = flag.Float64("maker_commission", 0, "")
var BinanceTakerCommission = flag.Float64("taker_commission", 0, "")

func NewBinanceCommissionModel() (*BinanceCommission, error) {
	err := checkIsValidCommissions(*BinanceMakerCommission, *BinanceTakerCommission)
	if err != nil {
		return nil, err
	}

	return &BinanceCommission{Maker: *BinanceMakerCommission, Taker: *BinanceTakerCommission}, nil
}

func (v *BinanceCommission) GetCommission(o *order.Order) float64 {
	if o.Type == order.Market || o.Type == order.StopMarket {
		return o.ExecutionPrice * o.Size * v.Taker
	}

	return o.ExecutionPrice * o.Size * v.Maker
}

func checkIsValidCommissions(values ...float64) error {
	for _, v := range values {
		if v < 0 || v >= 1 {
			return fmt.Errorf("binance commission should be between 0 and 1")
		}
	}

	return nil
}
