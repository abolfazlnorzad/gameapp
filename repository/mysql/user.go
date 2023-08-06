package mysql

import (
	"database/sql"
	"fmt"
	"gameapp/entity"
)

func (d *MySQLDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	var user entity.User
	var createdAt []uint8
	row := d.db.QueryRow("select * from users where phone_number = ?", phoneNumber)

	sErr := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password, &createdAt)
	if sErr != nil {
		if sErr == sql.ErrNoRows {
			return true, nil
		}
		return false, fmt.Errorf("can't scan the QueryRow : %w", sErr)
	}

	return false, nil
}

func (d *MySQLDB) Create(u entity.User) (entity.User, error) {
	result, err := d.db.Exec("insert into users(name,phone_number,password) values (?,?,?)", u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("can't exec insert command : %w", err)
	}

	id, _ := result.LastInsertId()
	u.ID = uint(id)
	return u, nil
}
