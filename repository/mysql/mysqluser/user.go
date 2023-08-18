package mysqluser

import (
	"database/sql"
	"gameapp/entity"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
	"gameapp/repository/mysql"
	"time"
)

func (d *DB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	const op = "mysql.IsPhoneNumberUnique"
	row := d.conn.Conn().QueryRow("select * from users where phone_number = ?", phoneNumber)

	_, sErr := ScanUser(row)
	if sErr != nil {
		if sErr == sql.ErrNoRows {
			return true, nil
		}
		return false, richerror.New(op).WithErr(sErr).WithMessage(errmsg.ErrorMsgCantScanQueryResult).
			WithKind(richerror.KindUnexpected).WithMeta(map[string]any{"phone_number": phoneNumber})
	}

	return false, nil
}

func (d *DB) Create(u entity.User) (entity.User, error) {
	const op = "mysql.Create"
	result, err := d.conn.Conn().Exec("insert into users(name,phone_number,password) values (?,?,?)", u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return entity.User{}, richerror.New(op).WithErr(err).WithMessage(errmsg.CannotInsertCommand).WithKind(richerror.KindUnexpected)
	}

	id, _ := result.LastInsertId()
	u.ID = uint(id)
	return u, nil
}

func (d *DB) GetUserByPhoneNumber(phoneNumber string) (entity.User, error) {
	const op = "mysql.GetUserByPhoneNumber"
	row := d.conn.Conn().QueryRow("select * from users where phone_number = ?", phoneNumber)
	user, sErr := ScanUser(row)
	if sErr != nil {
		if sErr == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(sErr).WithMessage(errmsg.ErrorMsgNotFound).
				WithKind(richerror.KindNotFound)
		}
		return entity.User{}, richerror.New(op).WithErr(sErr).WithMessage(errmsg.ErrorMsgCantScanQueryResult).
			WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func (d *DB) GetUserProfile(userID uint) (entity.User, error) {
	const op = "mysql.GetUserProfile"
	row := d.conn.Conn().QueryRow("select * from users where id = ?", userID)
	user, sErr := ScanUser(row)
	if sErr != nil {
		if sErr == sql.ErrNoRows {
			return entity.User{}, richerror.New(op).WithErr(sErr).WithMessage(errmsg.ErrorMsgNotFound).
				WithKind(richerror.KindNotFound).WithMeta(map[string]any{"userID": user})
		}
		return entity.User{}, richerror.New(op).WithErr(sErr).WithMessage(errmsg.ErrorMsgCantScanQueryResult).
			WithKind(richerror.KindUnexpected)
	}

	return user, nil
}

func ScanUser(scanner mysql.Scanner) (entity.User, error) {
	var user entity.User
	var createdAt time.Time
	var roleStr string
	sErr := scanner.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password, &createdAt, &roleStr)
	user.Role = entity.MapToRoleEntity(roleStr)
	return user, sErr
}
