package acchistory

import (
	"fmt"
	"os"
	"strconv"
	"stregy/internal/domain/order"
	"stregy/internal/domain/symbol"
)

/*
resulting csv format:

	order id, position id, contingent type, diraction, type, size, submision time, submision price, execution time, execution price
*/
func SaveOrderHistory(orders []*order.Order, s symbol.Symbol, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, o := range orders {
		f.WriteString(fmt.Sprintf(
			"%d,%d,%s,%s,%s,%f,%s,%s,%s,%s\n",
			o.ID,
			o.Position.ID,
			getOrderContingentTypeString(o),
			o.Diraction.String(),
			o.Type.String(),
			o.Size,
			o.SubmissionTime.Format("2006-01-02 15:04:05"),
			formatPrice(o.Price, s.Precision),
			o.FCTime.Format("2006-01-02 15:04:05"),
			formatPrice(o.ExecutionPrice, s.Precision)))
	}

	return nil
}

func getOrderContingentTypeString(o *order.Order) string {
	res := "Cont"
	if o.Position.MainOrder.ID == o.ID {
		res = "Main"
	}

	return res
}

func formatPrice(f float64, precision int) string {
	return strconv.FormatFloat(f, 'f', precision, 64)
}
