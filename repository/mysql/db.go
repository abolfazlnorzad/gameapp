package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type MySQLDB struct {
	db *sql.DB
}

func New() *MySQLDB {
	db, err := sql.Open("mysql", "root:password@/gameapp_db")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MySQLDB{
		db: db,
	}
}