package core

import (
	"fmt"
	"path"
	"stregy/internal/domain/broker"
	"stregy/internal/domain/order"
	"stregy/internal/domain/quote"
	strategy1 "stregy/internal/domain/strategy"
	"stregy/pkg/logging"
	"time"
)

func (b *Backtest) Time() time.Time {
	return b.curTime
}
func (b *Backtest) Price() float64 {
	return b.lastPrice
}

func (b *Backtest) BacktestOnQuotes(s strategy1.Strategy, quotes <-chan quote.Quote, firstQuote quote.Quote) error {
	b.init(s)
	b.curTime = firstQuote.Time
	b.lastPrice = firstQuote.Open
	quoteGen, err := NewQuoteGenerator(s, b.TimeframeSec, firstQuote)
	if err != nil {
		return err
	}
	b.Printf("running backtest with strategy strat1 on period period [%s; %s]", b.StartTime.Format("2006-01-02 15:04:05"), b.EndTime.Format("2006-01-02 15:04:05"))

	b.runOnQuotes(quotes, quoteGen)

	return nil
}

func (b *Backtest) CreateReport(location string) {
	reportPath := ""
	if location == "" {
		reportPath = b.getDefaultReportPath()
	} else {
		reportPath = path.Join(location, b.ID+".csv")
	}

	err := b.AccountHistoryService.CreateReport(b.orderHistory, b.Symbol, reportPath)
	if err != nil {
		logging.GetLogger().Error(fmt.Sprintf("Error creating backtest report: %v", err))
	}
}

func (b *Backtest) init(s strategy1.Strategy) {
	b.initLogger()

	b.strategy = s
	b.orders = make(map[int64]*order.Order)
	b.positions = make(map[int64]*order.Position)
	b.termChan = make(chan bool)
}

func (b *Backtest) initLogger() {
	loggerCfg := broker.LoggingConfig{LogOrderStatusChange: false, PricePrecision: b.Symbol.Precision}
	b.logger = *broker.NewLogger(b.ID+".log", loggerCfg, b)
}

func (b *Backtest) runOnQuotes(quotes <-chan quote.Quote, quoteGen *QuoteGenerator) {
	run := true
	for run {
		select {
		case q, ok := <-quotes:
			if !ok {
				run = false
				break
			}

			b.curTime = q.Time
			b.lastPrice = q.Close

			b.strategy.OnTick(q.Close)

			for _, o := range b.orders {
				if o.Type == order.Limit {
					if o.Diraction == order.Long {
						if q.Low <= o.Price {
							b.executeOrder(o, b.lastPrice)
							continue
						}
					} else {
						if q.High >= o.Price {
							b.executeOrder(o, b.lastPrice)
							continue
						}
					}
				} else if o.Type == order.StopMarket {
					if o.Price >= q.Low && o.Price <= q.High {
						b.executeOrder(o, q.Close)
						continue
					}
				} else if o.Type == order.Market {
					b.executeOrder(o, q.Close)
					continue
				}
			}

			quoteGen.OnQuote(q)

		case <-b.termChan:
			run = false
		}
	}
}

func (b *Backtest) Terminate() {
	b.Status = Terminated
	b.termChan <- true
	b.logger.Print("Terminated")
}
