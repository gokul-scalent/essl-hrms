package model

type UserRoleModel struct {
	UserID     int    `db:"user_id"`
	RoleID     int    `db:"role_id"`
	RoleCode   string `db:"role_code"`
	RoleName   string `db:"role_name"`
	RoleStatus string `db:"role_status"`
}
