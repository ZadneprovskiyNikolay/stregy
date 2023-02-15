package quote

import (
	"encoding/hex"
	"fmt"
	"stregy/internal/domain/quote"
	"stregy/pkg/logging"
	"stregy/pkg/utils"
	"strings"
	"time"
	"unsafe"
)

func (r repository) Get(
	dest []quote.Quote,
	symbol string,
	startTime, endTime time.Time,
) ([]quote.Quote, error) {

	tableName := getTableName(symbol)
	startTimeStr := utils.FormatTime(startTime)
	endTimeStr := utils.FormatTime(endTime)

	query := fmt.Sprintf(
		"SELECT * FROM \"%s\" WHERE time >= '%s' AND time < '%s' ORDER BY time",
		tableName,
		startTimeStr,
		endTimeStr)
	rows, _ := r.dbPQ.Query(query)
	if rows == nil {
		return []quote.Quote{}, nil
	}
	defer rows.Close()

	for rows.Next() {
		var t time.Time
		var ohlcvHexString string
		rows.Scan(&t, &ohlcvHexString)

		ohlcvBytesSlice, _ := hex.DecodeString(ohlcvHexString)

		ohlcv := (*OHLCV)(unsafe.Pointer((*[36]byte)(ohlcvBytesSlice)))
		dest = append(dest, quote.Quote{
			Time:   t,
			Open:   ohlcv.Open,
			High:   ohlcv.High,
			Low:    ohlcv.Low,
			Close:  ohlcv.Close,
			Volume: ohlcv.Volume})
	}

	return dest, nil
}

var logger logging.Logger

func (r *repository) GetFirst(symbol string, startTime, endTime time.Time) (quote.Quote, error) {
	tableName := getTableName(symbol)
	startTimeStr := utils.FormatTime(startTime)
	endTimeStr := utils.FormatTime(endTime)

	query := fmt.Sprintf("SELECT * FROM \"%s\" WHERE time >= '%s' AND time <= '%s' ORDER BY time LIMIT 1", tableName, startTimeStr, endTimeStr)
	rows, err := r.dbPQ.Query(
		query,
		tableName,
		startTimeStr,
		endTimeStr)
	if err != nil {
		return quote.Quote{}, err
	}
	if rows == nil {
		return quote.Quote{}, fmt.Errorf("no entries found")
	}
	defer rows.Close()

	var t time.Time
	var ohlcvHexString string
	rows.Scan(&t, &ohlcvHexString)
	ohlcvBytesSlice, _ := hex.DecodeString(ohlcvHexString)
	ohlcv := *(*OHLCV)(unsafe.Pointer((*[36]byte)(ohlcvBytesSlice)))
	return quote.Quote{
		Time:   t,
		Open:   ohlcv.Open,
		High:   ohlcv.High,
		Low:    ohlcv.Low,
		Close:  ohlcv.Close,
		Volume: ohlcv.Volume}, nil
}

func getTableName(symbol string) string {
	tableName := strings.ToLower(symbol)
	tableName += "_quotes"
	return tableName
}

func nextCircularIndex(i, length int) int {
	return (i + 1) % length
}
