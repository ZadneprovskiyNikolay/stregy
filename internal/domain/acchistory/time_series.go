package acchistory

import (
	"fmt"
	"os"
	"stregy/pkg/utils"
	"time"
)

type TimeSeries []TSValue

type TSValue struct {
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
		f.WriteString(fmt.Sprintf("%v,%v\n", utils.FormatTime(v.Time), v.Value))
	}

	return nil
}
