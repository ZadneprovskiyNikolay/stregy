package draw

import (
	"net/http"
	"os/exec"
	"stregy/pkg/timeseries"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type LineChartPlotter struct {
	PageName string
	Charts   []LineChart

	shutDown chan bool
}

type LineChart struct {
	Name string
	X    []interface{}
	Y    []opts.LineData
}

func DrawLineCharts(pageName string, charts ...LineChart) {
	l := LineChartPlotter{PageName: pageName, Charts: charts, shutDown: make(chan bool)}
	http.HandleFunc("/"+pageName, l.draw)

	go func() {
		http.ListenAndServe(":8081", nil)
	}()

	go func() {
		cmd := exec.Command("explorer", "http://localhost:8081/"+pageName)
		cmd.Run()
	}()

	<-l.shutDown
}

func FromTimeSeries(name string, timeSeries timeseries.TimeSeries) LineChart {
	x := make([]interface{}, 0, len(timeSeries))
	y := make([]opts.LineData, 0, len(timeSeries))

	for _, tsValue := range timeSeries {
		x = append(x, tsValue.Time)
		y = append(y, opts.LineData{Value: tsValue.Value})
	}

	return LineChart{Name: name, X: x, Y: y}
}

func (l *LineChartPlotter) draw(w http.ResponseWriter, _ *http.Request) {
	// lines := make([]charts.Line, len(l.Charts))
	defer func() {
		time.Sleep(time.Second * 1)
		l.shutDown <- true
	}()

	for _, chart := range l.Charts {
		line := charts.NewLine()

		line.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{
				PageTitle: l.PageName,
				Width:     "1920px",
				Height:    "950px",
				Theme:     "dark"}),
			charts.WithTitleOpts(opts.Title{
				Title: chart.Name,
			}))

		line.SetXAxis(chart.X).
			AddSeries("y", chart.Y)
		line.Render(w)
	}
}
