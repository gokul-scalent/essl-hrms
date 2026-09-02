package model

import (
	"database/sql"
	"time"
)

type AttendanceLog struct {
	ID              int          `db:"id"`
	UID             int          `db:"uid"`
	EmpID           string       `db:"emp_id"`
	EmpName         string       `db:"emp_name"`
	Timestamp       time.Time    `db:"timestamp"`
	Status          int          `db:"status"`
	Punch           int          `db:"punch"`
	AttendanceState string       `db:"attendance_state"`
	DeviceName      string       `db:"device_name"`
	CreatedAt       sql.NullTime `db:"created_at"`
}

var AttendanceLogModelMap = map[string]FieldStruct{
	"ID":              {MySQLDatatype: "int", FieldName: "id"},
	"UID":             {MySQLDatatype: "int", FieldName: "uid"},
	"EmpID":           {MySQLDatatype: "varchar", FieldName: "emp_id"},
	"EmpName":         {MySQLDatatype: "varchar", FieldName: "emp_name"},
	"Timestamp":       {MySQLDatatype: "datetime", FieldName: "timestamp"},
	"Status":          {MySQLDatatype: "int", FieldName: "status"},
	"Punch":           {MySQLDatatype: "int", FieldName: "punch"},
	"AttendanceState": {MySQLDatatype: "varchar", FieldName: "attendance_state"},
	"DeviceName":      {MySQLDatatype: "varchar", FieldName: "device_name"},
	"CreatedAt":       {MySQLDatatype: "timestamp", FieldName: "created_at"},
}

type DailyAttendanceLog struct {
	EmpID     string        `db:"emp_id"`
	EmpName   string        `db:"emp_name"`
	LogDate   time.Time     `db:"log_date"`
	Timestamp sql.NullTime  `db:"timestamp"`
	Punch     sql.NullInt64 `db:"punch"`
}
