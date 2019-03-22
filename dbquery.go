package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Host     string `json:"Host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"dbName"`
	Timeout  string `json:"timeout"`
}

type DB struct {
	dsn string
	db  *sql.DB
}

func NewDB(config *DBConfig) *DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%s", config.User, config.Password, config.Host, config.DbName, config.Timeout)
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
