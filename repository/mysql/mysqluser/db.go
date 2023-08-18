package mysqluser

import "gameapp/repository/mysql"

type DB struct {
	conn *mysql.MySQLDB
}

func New(c *mysql.MySQLDB) *DB {
	return &DB{
		conn: c,
	}
}
