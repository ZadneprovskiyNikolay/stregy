package broker

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"stregy/internal/domain/order"
	"stregy/pkg/utils"
	"strings"
	"time"
)

type Logger struct {
	Config LoggingConfig
	file   *os.File

	clock Clock
}

type Clock interface {
	Time() time.Time
}

type LoggingConfig struct {
	LogOrderStatusChange bool
	PricePrecision       int
	OmmitTimePrefix      bool
}

func NewLogger(logName string, cfg LoggingConfig, clock Clock) *Logger {
	wd, _ := os.Getwd()
	dir := path.Join(wd, "logs", "backtest")
	os.Mkdir(dir, os.ModePerm)
	fpath := path.Join(dir, logName)
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(fmt.Errorf("could not create log file: %s", err.Error()))
	}

	return &Logger{file: f, Config: cfg, clock: clock}
}

func (l *Logger) Print(s string) {
	if !l.Config.OmmitTimePrefix {
		s = l.timePrefix() + s
	}
	fmt.Println(s)
	l.file.WriteString(s + "\n")
}

func (l *Logger) Printf(format string, v ...interface{}) {
	for i := 0; i < len(v); i++ {
		switch v[i].(type) {
		case order.Order:
			o := v[i].(order.Order)
			v[i] = l.FormatOrder(&o)
		case *order.Order:
			o := v[i].(*order.Order)
			v[i] = l.FormatOrder(o)
		}
	}

	l.Print(fmt.Sprintf(format, v...))
}

func (l *Logger) LogOrderStatusChange(o *order.Order) {
	if l.Config.LogOrderStatusChange {
		l.PrintOrderStatus(o)
	}
}

func (l *Logger) FormatOrder(o *order.Order) string {
	executionPrice := strconv.FormatFloat(o.ExecutionPrice, 'f', l.Config.PricePrecision, 64)
	if o.Status != order.FilledOrder {
		executionPrice = ""
	}
	price := strconv.FormatFloat(o.Price, 'f', l.Config.PricePrecision, 64)
	if o.Price == 0 {
		price = strings.Repeat(" ", len(executionPrice))
	}
	return fmt.Sprintf("Order #%d %s %s %s : %s %s", o.ID, formatOrderDiraction(o.Diraction), formatOrderType(o.Type), price, o.Status.String(), executionPrice)
}

func (l *Logger) PrintOrderStatus(o *order.Order) {
	l.Printf("Order #%d: %s", o.ID, o.Status.String())
}

func (l *Logger) timePrefix() string {
	return l.clock.Time().Format("2006-01-02 15:04:05") + ": "
}

func formatOrderDiraction(d order.OrderDiraction) string {
	const maxLength = 10
	return utils.AddTrailingWhitespaces(d.String(), maxLength)
}

func formatOrderType(t order.OrderType) string {
	const maxLength = 12
	return utils.AddTrailingWhitespaces(t.String(), maxLength)
}
