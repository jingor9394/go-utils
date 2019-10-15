package utils

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

type DB struct {
	dsn string
	db  *sql.DB
}

func NewDB(dsn string) *DB {
	d := &DB{
		dsn,
		nil,
	}
	return d
}

func (d *DB) Open() (err error) {
	if d.db == nil || d.db.Ping() != nil {
		d.db, err = sql.Open("mysql", d.dsn)
	}
	return
}

func (d *DB) Query(query string, values ...interface{}) (*sql.Rows, error) {
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

func (d *DB) Fetch(rows *sql.Rows) ([]map[string]string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var list []map[string]string
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		row := make(map[string]string)
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		for i, col := range values {
			var value string
			switch val := col.(type) {
			case []uint8:
				value = string(val)
			case int64:
				value = strconv.Itoa(int(val))
			case nil:
				value = ""
			default:
				errMsg := fmt.Sprintf("unexpected type %T", value)
				return nil, errors.New(errMsg)
			}
			row[columns[i]] = value
		}
		list = append(list, row)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (d *DB) Exec(querySql string, values ...interface{}) (int64, error) {
	err := d.Open()
	if err != nil {
		return 0, err
	}
	stmt, err := d.db.Prepare(querySql)
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
