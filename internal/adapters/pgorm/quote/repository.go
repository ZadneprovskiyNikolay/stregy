package quote

import (
	"database/sql"
	"fmt"
	"stregy/internal/domain/quote"
	"stregy/pkg/utils"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"gorm.io/gorm"
)

type repository struct {
	db     *sql.DB
	dbGorm *gorm.DB
}

func NewRepository(client *gorm.DB) quote.Repository {
	connStr := "user=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return &repository{db: db, dbGorm: client}
}

func (r repository) Get(
	dest []quote.Quote,
	symbol string,
	startTime, endTime time.Time,
	limit, timeframeSec int,
) ([]quote.Quote, error) {

	tableName := getTableName(symbol, timeframeSec)
	startTimeStr := utils.FormatTime(startTime)

	rows, _ := r.db.Query(fmt.Sprintf("SELECT * FROM \"%s\" WHERE time >= '%s' ORDER BY time LIMIT %d", tableName, startTimeStr, limit))
	defer rows.Close()

	for rows.Next() {
		var t time.Time
		var o, h, l, c float64
		var v int32
		rows.Scan(&t, &o, &h, &l, &c, &v)
		if t.After(endTime) {
			break
		}
		dest = append(dest, quote.Quote{Time: t, Open: o, High: h, Low: l, Close: c, Volume: v})
	}

	return dest, nil
}

// GetAndPushToChan implements quote.Repository
func (r *repository) GetAndPushToChan(
	dest chan<- quote.Quote,
	symbol string,
	startTime time.Time,
	endTime time.Time,
	limit int,
	timeframeSec int,
) (err error, lastQuoteTime time.Time) {
	tableName := getTableName(symbol, timeframeSec)
	startTimeStr := utils.FormatTime(startTime)

	rows, _ := r.db.Query(fmt.Sprintf("SELECT * FROM \"%s\" WHERE time >= '%s' ORDER BY time LIMIT %d", tableName, startTimeStr, limit))
	defer rows.Close()

	for rows.Next() {
		var t time.Time
		var o, h, l, c float64
		var v int32
		rows.Scan(&t, &o, &h, &l, &c, &v)
		if t.After(endTime) {
			break
		}
		lastQuoteTime = t
		dest <- quote.Quote{Time: t, Open: o, High: h, Low: l, Close: c, Volume: v}
	}
	return nil, lastQuoteTime
}

func (r repository) Load(symbol, filePath, delimiter string, timeframeSec int) error {
	tableName := getTableName(symbol, timeframeSec)
	return r.dbGorm.Exec(fmt.Sprintf(`
	CREATE UNLOGGED TABLE IF NOT EXISTS temp_quotes (
		time double precision,
		open double precision,
		high double precision,
		low double precision,
		close double precision,
		volume int
	 );

	COPY temp_quotes FROM '%v' DELIMITERS '%v' CSV;

	ALTER TABLE temp_quotes
	ALTER time TYPE timestamp without time zone
		USING (to_timestamp(time) AT TIME ZONE 'UTC');

	CREATE TABLE IF NOT EXISTS %v (LIKE quotes INCLUDING ALL);

	INSERT INTO %v SELECT * FROM temp_quotes ON CONFLICT DO NOTHING;

	DROP TABLE temp_quotes;`,
		filePath, delimiter, tableName, tableName)).Error
}

func getTableName(symbol string, timeframeSec int) string {
	tableName := strings.ToLower(symbol)
	if timeframeSec < 60 {
		tableName += "_s1_quotes"
	} else {
		tableName += "_m1_quotes"
	}

	return tableName
}
