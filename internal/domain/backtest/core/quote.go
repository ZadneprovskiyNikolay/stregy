package core

import (
	"stregy/internal/domain/quote"
	"time"
)

type quoteFeed struct {
	open              float64
	high              float64
	low               float64
	close             float64
	openTime          time.Time
	closeTime         time.Time
	timeframe         int
	timeframeDuration time.Duration
}

type QuoteGenerator struct {
	feeds          []*quoteFeed
	callbackObject callbackObject
}

type callbackObject interface {
	OnQuote(quote quote.Quote, timeframe int)
	QuoteTimeframesNeeded() []int
}

func NewQuoteGenerator(cb callbackObject, firstQuote quote.Quote) (*QuoteGenerator, error) {
	timeframes := cb.QuoteTimeframesNeeded()

	feeds := make([]*quoteFeed, 0, len(timeframes))
	for _, timeframe := range timeframes {
		if err := quote.CheckIsValidTimeframe(timeframe); err != nil {
			return nil, err
		}

		timeframeDuration := time.Duration(timeframe) * time.Minute
		openTime := StartTimeOfBar(firstQuote.Time, timeframe)
		closeTime := openTime.Add(timeframeDuration)
		feeds = append(feeds, &quoteFeed{
			open:              firstQuote.Open,
			low:               firstQuote.Open,
			openTime:          openTime,
			closeTime:         closeTime,
			timeframe:         timeframe,
			timeframeDuration: timeframeDuration,
		})
	}

	qg := QuoteGenerator{
		feeds:          feeds,
		callbackObject: cb,
	}

	return &qg, nil
}

func (qg *QuoteGenerator) OnQuote(q quote.Quote) {
	for _, feed := range qg.feeds {
		newQuote := feed.feed(q)
		if newQuote != nil {
			qg.callbackObject.OnQuote(*newQuote, feed.timeframe)
		}
	}
}

func (f *quoteFeed) feed(q quote.Quote) *quote.Quote {
	var feedQuote *quote.Quote
	if !q.Time.Before(f.closeTime) {
		feedQuote = &quote.Quote{
			Time:  f.openTime,
			Open:  f.open,
			High:  f.high,
			Low:   f.low,
			Close: f.close}

		f.open = q.Open
		f.high = q.High
		f.low = q.Low
		f.close = q.Close
		f.openTime = StartTimeOfBar(q.Time, f.timeframe)
		f.closeTime = f.openTime.Add(f.timeframeDuration)
	} else {
		if q.High > f.high {
			f.high = q.High
		}
		if q.Low < f.low {
			f.low = q.Low
		}
		f.close = q.Close
	}

	return feedQuote
}

func StartTimeOfBar(t time.Time, timeframe int) time.Time {
	dayMinutes := (t.Hour()*60 + t.Minute())
	durFromStart := time.Duration(dayMinutes%timeframe) * time.Minute
	return t.Add(-durFromStart)
}
