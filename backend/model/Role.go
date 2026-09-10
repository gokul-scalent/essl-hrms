package model

import (
	"database/sql"
	"time"
)

type Role struct {
	ID        int          `db:"id"`
	Name      string       `db:"name"`
	Code      string       `db:"code"`
	Status    string       `db:"status"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

var RoleModelMap = map[string]FieldStruct{
	"ID":        {MySQLDatatype: "int", FieldName: "id"},
	"Name":      {MySQLDatatype: "varchar", FieldName: "name"},
	"Code":      {MySQLDatatype: "varchar", FieldName: "code"},
	"Status":    {MySQLDatatype: "enum", FieldName: "status"},
	"CreatedAt": {MySQLDatatype: "datetime", FieldName: "created_at"},
	"UpdatedAt": {MySQLDatatype: "datetime", FieldName: "updated_at"},
	"DeletedAt": {MySQLDatatype: "datetime", FieldName: "deleted_at"},
}
