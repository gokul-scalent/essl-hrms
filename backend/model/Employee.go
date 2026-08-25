package model

import (
	"database/sql"
)

type Employee struct {
	ID        int            `db:"id"`
	UID       int            `db:"uid"`
	EmpID     string         `db:"emp_id"`
	EmpName   string         `db:"emp_name"`
	Privilege sql.NullInt64  `db:"privilege"`
	Password  sql.NullString `db:"password"`
	GroupID   sql.NullString `db:"group_id"`
	Card      sql.NullString `db:"card"`
	CreatedAt sql.NullTime   `db:"created_at"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
}

var EmployeeModelMap = map[string]FieldStruct{
	"ID":        {MySQLDatatype: "int", FieldName: "id"},
	"UID":       {MySQLDatatype: "int", FieldName: "uid"},
	"EmpID":     {MySQLDatatype: "varchar", FieldName: "emp_id"},
	"EmpName":   {MySQLDatatype: "varchar", FieldName: "emp_name"},
	"Privilege": {MySQLDatatype: "int", FieldName: "privilege"},
	"Password":  {MySQLDatatype: "varchar", FieldName: "password"},
	"GroupID":   {MySQLDatatype: "varchar", FieldName: "group_id"},
	"Card":      {MySQLDatatype: "bigint", FieldName: "card"},
	"CreatedAt": {MySQLDatatype: "timestamp", FieldName: "created_at"},
	"DeletedAt": {MySQLDatatype: "datetime", FieldName: "deleted_at"},
}
