package model

import (
	"database/sql"
	"time"
)

type EmailList struct {
	ID               int          `db:"id"`
	UserID           int          `db:"user_id"`
	Name             string       `db:"name"`
	TotalRecords     int          `db:"total_records"`
	ProcessedRecords int          `db:"processed_records"`
	SafeCount        int          `db:"safe_count"`
	RiskyCount       int          `db:"risky_count"`
	InvalIDCount     int          `db:"invalid_count"`
	UnknownCount     int          `db:"unknown_count"`
	PendingCount     int          `db:"pending_count"`
	Priority         string       `db:"priority"`
	CreatedAt        time.Time    `db:"created_at"`
	UpdatedAt        sql.NullTime `db:"updated_at"`
	DeletedAt        sql.NullTime `db:"deleted_at"`
}

var EmailListModelMap = map[string]FieldStruct{
	"ID":               {MySQLDatatype: "int", FieldName: "id"},
	"UserID":           {MySQLDatatype: "int", FieldName: "user_id"},
	"Name":             {MySQLDatatype: "varchar", FieldName: "name"},
	"TotalRecords":     {MySQLDatatype: "int", FieldName: "total_records"},
	"ProcessedRecords": {MySQLDatatype: "int", FieldName: "processed_records"},
	"SafeCount":        {MySQLDatatype: "int", FieldName: "safe_count"},
	"RiskyCount":       {MySQLDatatype: "int", FieldName: "risky_count"},
	"InvalIDCount":     {MySQLDatatype: "int", FieldName: "invalid_count"},
	"UnknownCount":     {MySQLDatatype: "int", FieldName: "unknown_count"},
	"PendingCount":     {MySQLDatatype: "int", FieldName: "pending_count"},
	"Priority":         {MySQLDatatype: "enum", FieldName: "priority"},
	"CreatedAt":        {MySQLDatatype: "datetime", FieldName: "created_at"},
	"UpdatedAt":        {MySQLDatatype: "datetime", FieldName: "updated_at"},
	"DeletedAt":        {MySQLDatatype: "datetime", FieldName: "deleted_at"},
}
