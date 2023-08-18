package mysqlacl

import (
	"gameapp/entity"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
	"gameapp/pkg/slice"
	"gameapp/repository/mysql"
	"strings"
	"time"
)

func (d *DB) GetUserPermissionTitle(userID uint, role entity.Role) ([]entity.PermissionTitle, error) {
	const op = "mysqlacl.GetUserPermissionTitle"
	// get rule permissions

	roleAcl := make([]entity.Acl, 0)
	rows, err := d.conn.Conn().Query("select * from acls where actor_type=? and actor_id=?",
		entity.RoleActorType, role)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).WithMessage(errmsg.SomethingWentWrong).
			WithKind(richerror.KindUnexpected)
	}
	defer rows.Close()
	for rows.Next() {
		acl, err := ScanAcl(rows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).
				WithMessage(errmsg.SomethingWentWrong).WithKind(richerror.KindUnexpected)
		}
		roleAcl = append(roleAcl, acl)
	}
	err = rows.Err()
	if err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errmsg.SomethingWentWrong).WithKind(richerror.KindUnexpected)
	}

	//userAcl := make([]entity.Acl, 0)
	//uRows, err := d.conn.Conn().Query("select * from acls where actor_type=? and actor_id=?", userAcl, userID)
	//if err != nil {
	//	return nil, richerror.New(op).WithErr(err).WithMessage(errmsg.SomethingWentWrong).
	//		WithKind(richerror.KindUnexpected)
	//}
	//defer uRows.Close()
	//for uRows.Next() {
	//	acl, err := ScanAcl(uRows)
	//	if err != nil {
	//		return nil, richerror.New(op).WithErr(err).
	//			WithMessage(errmsg.SomethingWentWrong).WithKind(richerror.KindUnexpected)
	//	}
	//
	//	userAcl = append(userAcl, acl)
	//}
	//err = uRows.Err()
	//if err != nil {
	//	return nil, richerror.New(op).WithErr(err).
	//		WithMessage(errmsg.SomethingWentWrong).WithKind(richerror.KindUnexpected)
	//}
	permissionIDs := make([]uint, 0)

	for _, r := range roleAcl {
		if !slice.DoesExist(permissionIDs, r.PermissionID) {
			permissionIDs = append(permissionIDs, r.PermissionID)
		}
	}

	//for _, r := range userAcl {
	//	if !slice.DoesExist(permissionIDs, r.PermissionID) {
	//		permissionIDs = append(permissionIDs, r.PermissionID)
	//	}
	//}

	if len(permissionIDs) == 0 {
		return nil, nil
	}

	var args = make([]any, len(permissionIDs))

	for i, pID := range permissionIDs {
		args[i] = pID
	}

	query := "select * from permissions where id in (?" + strings.Repeat(",?", len(permissionIDs)-1) + ")"
	pRows, err := d.conn.Conn().Query(query, args...)
	if err != nil {
		return nil, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected).WithMessage(errmsg.SomethingWentWrong)
	}
	defer pRows.Close()
	permissionTitles := make([]entity.PermissionTitle, 0)
	for pRows.Next() {
		permission, err := scanPermission(pRows)
		if err != nil {
			return nil, richerror.New(op).WithErr(err).WithKind(richerror.KindUnexpected).WithMessage(errmsg.SomethingWentWrong)
		}
		permissionTitles = append(permissionTitles, permission.Title)
	}
	if err = pRows.Err(); err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithMessage(errmsg.SomethingWentWrong).WithKind(richerror.KindUnexpected)
	}
	return permissionTitles, nil

}

func ScanAcl(s mysql.Scanner) (entity.Acl, error) {
	var CreatedAt time.Time
	var acl entity.Acl
	err := s.Scan(&acl.ID, &acl.ActorID, &acl.ActorType, &acl.PermissionID, &CreatedAt)
	return acl, err
}
