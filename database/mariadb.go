package database

import (
	"database/sql"
	"fmt"
)

func GetMariaDBConnectionString(username, password, host, dbname string, port int) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=1&charset=utf8mb4&collation=utf8mb4_unicode_ci", username, password, host, port, dbname)
}

func ConnectMariaDB(uri string) (*sql.DB, error) {
	db, err := sql.Open("mysql", uri)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)

	return db, nil
}
