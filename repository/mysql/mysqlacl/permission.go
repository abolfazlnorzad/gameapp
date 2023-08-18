package mysqlacl

import (
	"gameapp/entity"
	"gameapp/repository/mysql"
	"time"
)

func scanPermission(s mysql.Scanner) (entity.Permission, error) {
	var createdAt time.Time
	var p entity.Permission
	err := s.Scan(&p.ID, &p.Title, &createdAt)

	return p, err
}
