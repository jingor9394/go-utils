package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DBQuery struct {
	dsn string
	db  *sql.DB
}

func NewDBQuery(config *DBConfig) *DBQuery {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%s", config.User, config.Password, config.Host, config.DbName, config.Timeout)
	d := &DBQuery{
		dsn,
		nil,
	}
	return d
}

func (d *DBQuery) Open() (err error) {
	if d.db == nil || d.db.Ping() != nil {
		d.db, err = sql.Open("mysql", d.dsn)
	}
	return
}

func (d *DBQuery) Query(query string, values ...interface{}) (*sql.Rows, error) {
	err := d.Open()
	if err != nil {
		return nil, err
	}
	rows, err := d.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *DBQuery) Fetch(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var list []map[string]interface{}
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values{
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		row := make(map[string]interface{})
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		for i, col := range values {
			row[columns[i]] = col
		}
		list = append(list, row)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (d *DBQuery) Exec(query string, values ...interface{}) (int64, error) {
	err := d.Open()
	if err != nil {
		return 0, err
	}
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	ret, err := stmt.Exec(values...)
	if err != nil {
		return 0, err
	}
	affectedRows, err := ret.RowsAffected()
	return affectedRows, nil
}
