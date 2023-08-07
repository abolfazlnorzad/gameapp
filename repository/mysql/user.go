package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"gameapp/entity"
)

func (d *MySQLDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	row := d.db.QueryRow("select * from users where phone_number = ?", phoneNumber)

	_, sErr := ScanUser(row)
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

func (d *MySQLDB) GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error) {
	row := d.db.QueryRow("select * from users where phone_number = ?", phoneNumber)
	user, sErr := ScanUser(row)
	if sErr != nil {
		if sErr == sql.ErrNoRows {
			return entity.User{}, false, nil
		}
		return entity.User{}, false, fmt.Errorf("can't scan the QueryRow : %w", sErr)
	}

	return user, true, nil
}

func (d *MySQLDB) GetUserProfile(userID uint) (entity.User, error) {
	row := d.db.QueryRow("select * from users where id = ?", userID)

	user, sErr := ScanUser(row)
	if sErr != nil {
		if sErr == sql.ErrNoRows {
			return entity.User{}, errors.New("record not found")
		}
		return entity.User{}, fmt.Errorf("can't scan the QueryRow : %w", sErr)
	}

	return user, nil
}

func ScanUser(row *sql.Row) (entity.User, error) {
	var user entity.User
	var createdAt []uint8

	sErr := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password, &createdAt)
	return user, sErr
}
