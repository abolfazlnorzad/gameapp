package mysql

import (
	"database/sql"
	"gameapp/entity"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
	"time"
)

func (d *MySQLDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	const op = "mysql.IsPhoneNumberUnique"
	row := d.db.QueryRow("select * from users where phone_number = ?", phoneNumber)

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

func (d *MySQLDB) Create(u entity.User) (entity.User, error) {
	const op = "mysql.Create"
	result, err := d.db.Exec("insert into users(name,phone_number,password) values (?,?,?)", u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return entity.User{}, richerror.New(op).WithErr(err).WithMessage(errmsg.CannotInsertCommand).WithKind(richerror.KindUnexpected)
	}

	id, _ := result.LastInsertId()
	u.ID = uint(id)
	return u, nil
}

func (d *MySQLDB) GetUserByPhoneNumber(phoneNumber string) (entity.User, error) {
	const op = "mysql.GetUserByPhoneNumber"
	row := d.db.QueryRow("select * from users where phone_number = ?", phoneNumber)
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

func (d *MySQLDB) GetUserProfile(userID uint) (entity.User, error) {
	const op = "mysql.GetUserProfile"
	row := d.db.QueryRow("select * from users where id = ?", userID)
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

func ScanUser(row *sql.Row) (entity.User, error) {
	var user entity.User
	var createdAt time.Time

	sErr := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password, &createdAt)
	return user, sErr
}
