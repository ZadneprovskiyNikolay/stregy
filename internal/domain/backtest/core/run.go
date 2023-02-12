package core

import (
	"stregy/internal/domain/broker"
	"stregy/internal/domain/order"
	"stregy/internal/domain/quote"
	strategy1 "stregy/internal/domain/strategy"
	"time"
)

func (b *Backtest) Time() time.Time {
	return b.time
}
func (b *Backtest) Price() float64 {
	return b.price
}

func (b *Backtest) BacktestOnQuotes(
	s strategy1.Strategy,
	quotes <-chan quote.Quote,
	firstQuote quote.Quote,
	balance float64,
) error {
	b.init(s, firstQuote, balance)

	quoteGen, err := NewQuoteGenerator(s, b.TimeframeSec, firstQuote)
	if err != nil {
		return err
	}
	b.Printf("running backtest with strategy strat1 on period period [%s; %s]", b.StartTime.Format("2006-01-02 15:04:05"), b.EndTime.Format("2006-01-02 15:04:05"))

	b.runOnQuotes(quotes, quoteGen)

	return nil
}

func (b *Backtest) init(s strategy1.Strategy, q quote.Quote, balance float64) {
	b.time = q.Time
	b.price = q.Open

	b.initLogger()

	b.strategy = s
	b.orders = make(map[int64]*order.Order)
	b.positions = make(map[int64]*order.Position)
	b.termChan = make(chan bool)
	b.Balance.Update(balance, b.time)
	b.Drawdown.Init(b.time, time.Minute*5, b)
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

			b.time = q.Time
			b.price = q.Close

			b.Drawdown.Update(b.time)

			b.strategy.OnTick(q.Close)

			for _, o := range b.orders {
				if o.Type == order.Limit {
					if o.Diraction == order.Long {
						if q.Low <= o.Price {
							b.executeOrder(o, b.price)
							continue
						}
					} else {
						if q.High >= o.Price {
							b.executeOrder(o, b.price)
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
