package model

import (
	"database/sql"
	"time"
)

type UserSetting struct {
	ID                   int          `db:"id"`
	UserID               int          `db:"user_id"`
	VerificationInterval int          `db:"verification_interval"`
	CreatedAt            time.Time    `db:"created_at"`
	UpdatedAt            sql.NullTime `db:"updated_at"`
	DeletedAt            sql.NullTime `db:"deleted_at"`
}

var UserSettingModelMap = map[string]FieldStruct{
	"ID":                   {MySQLDatatype: "int", FieldName: "id"},
	"UserID":               {MySQLDatatype: "int", FieldName: "user_id"},
	"VerificationInterval": {MySQLDatatype: "int", FieldName: "verification_interval"},
	"CreatedAt":            {MySQLDatatype: "datetime", FieldName: "created_at"},
	"UpdatedAt":            {MySQLDatatype: "datetime", FieldName: "updated_at"},
	"DeletedAt":            {MySQLDatatype: "datetime", FieldName: "deleted_at"},
}
