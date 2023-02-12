package timeseries

import (
	"fmt"
	"os"
	"time"
)

type TimeSeries []Value

type Value struct {
	Time  time.Time
	Value float64
}

func (t TimeSeries) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, v := range t {
		f.WriteString(fmt.Sprintf("%v,%v\n", v.Time.Format("2006-01-02 15:04"), v.Value))
	}

	return nil
}
