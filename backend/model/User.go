package model

import (
	"database/sql"
)

type User struct {
	ID            int            `db:"id"`
	Email         sql.NullString `db:"email"`
	Password      sql.NullString `db:"password"`
	IsPasswordSet string         `db:"is_password_set"`
	Status        sql.NullString `db:"status"`
	RoleID        sql.NullInt64  `db:"role_id"`
	RoleName      sql.NullString `db:"role_name"`
	RoleCode      sql.NullString `db:"role_code"`
	RoleStatus    sql.NullString `db:"role_status"`
	LastLoginAt   sql.NullTime   `db:"last_login_at"`
	SessionToken  sql.NullString `db:"session_token"`
	CreatedAt     sql.NullTime   `db:"created_at"`
	UpdatedAt     sql.NullTime   `db:"updated_at"`
	DeletedAt     sql.NullTime   `db:"deleted_at"`
}

var UserModelMap = map[string]FieldStruct{
	"ID":            {MySQLDatatype: "int", FieldName: "id"},
	"Email":         {MySQLDatatype: "varchar", FieldName: "email"},
	"Password":      {MySQLDatatype: "varchar", FieldName: "password"},
	"Status":        {MySQLDatatype: "enum", FieldName: "status"},
	"IsPasswordSet": {MySQLDatatype: "enum", FieldName: "is_password_set"},
	"LastLoginAt":   {MySQLDatatype: "datetime", FieldName: "last_login_at"},
	"SessionToken":  {MySQLDatatype: "varchar", FieldName: "session_token"},
	"CreatedAt":     {MySQLDatatype: "datetime", FieldName: "created_at"},
	"UpdatedAt":     {MySQLDatatype: "datetime", FieldName: "updated_at"},
	"DeletedAt":     {MySQLDatatype: "datetime", FieldName: "deleted_at"},
	"RoleID":        {MySQLDatatype: "int", FieldName: "ur.role_id"},
	"RoleName":      {MySQLDatatype: "varchar", FieldName: "r.name"},
	"RoleCode":      {MySQLDatatype: "varchar", FieldName: "r.code"},
	"RoleStatus":    {MySQLDatatype: "enum", FieldName: "r.status"},
}
