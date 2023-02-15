package quote

import "fmt"

func (r *repository) Upload(symbol, filePath, delimiter string, timeframeSec int) error {
	tableName := getTableName(symbol)
	return r.dbGORM.Exec(fmt.Sprintf(`
	CREATE UNLOGGED TABLE IF NOT EXISTS temp_quotes (
		time double precision,
		oplcv bytea
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
