package model

import (
	"database/sql"
	"time"
)

type Lead struct {
	ID                 int            `db:"id"`
	EmailListID        int            `db:"email_list_id"`
	Priority           string         `db:"priority"`
	Email              string         `db:"email"`
	EmailProvIDer      sql.NullString `db:"email_provider"`
	FirstName          sql.NullString `db:"first_name"`
	LastName           sql.NullString `db:"last_name"`
	JobTitle           sql.NullString `db:"job_title"`
	Company            sql.NullString `db:"company"`
	City               sql.NullString `db:"city"`
	Country            sql.NullString `db:"country"`
	Industry           sql.NullString `db:"industry"`
	IsSafe             string         `db:"is_safe"`
	FinalStatus        string         `db:"final_status"`
	IsReachable        sql.NullString `db:"is_reachable"`
	IsDisposable       sql.NullBool   `db:"is_disposable"`
	IsRoleAccount      sql.NullBool   `db:"is_role_account"`
	VerificationStatus string         `db:"verification_status"`
	VerifiedOn         sql.NullTime   `db:"verified_on"`
	CreatedAt          time.Time      `db:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at"`
	DeletedAt          sql.NullTime   `db:"deleted_at"`
	RetryCount         int            `db:"retry_count"`
	NextRetryAt        sql.NullTime   `db:"next_retry_at"`
}

var LeadModelMap = map[string]FieldStruct{
	"ID":                 {MySQLDatatype: "int", FieldName: "id"},
	"EmailListID":        {MySQLDatatype: "int", FieldName: "email_list_id"},
	"Priority":           {MySQLDatatype: "enum", FieldName: "priority"},
	"Email":              {MySQLDatatype: "varchar", FieldName: "email"},
	"EmailProvIDer":      {MySQLDatatype: "varchar", FieldName: "email_provider"},
	"FirstName":          {MySQLDatatype: "varchar", FieldName: "first_name"},
	"LastName":           {MySQLDatatype: "varchar", FieldName: "last_name"},
	"JobTitle":           {MySQLDatatype: "varchar", FieldName: "job_title"},
	"Company":            {MySQLDatatype: "varchar", FieldName: "company"},
	"City":               {MySQLDatatype: "varchar", FieldName: "city"},
	"Country":            {MySQLDatatype: "varchar", FieldName: "country"},
	"Industry":           {MySQLDatatype: "varchar", FieldName: "industry"},
	"IsSafe":             {MySQLDatatype: "enum", FieldName: "is_safe"},
	"FinalStatus":        {MySQLDatatype: "enum", FieldName: "final_status"},
	"IsReachable":        {MySQLDatatype: "enum", FieldName: "is_reachable"},
	"IsDisposable":       {MySQLDatatype: "tinyint", FieldName: "is_disposable"},
	"IsRoleAccount":      {MySQLDatatype: "tinyint", FieldName: "is_role_account"},
	"VerificationStatus": {MySQLDatatype: "enum", FieldName: "verification_status"},
	"VerifiedOn":         {MySQLDatatype: "datetime", FieldName: "verified_on"},
	"CreatedAt":          {MySQLDatatype: "datetime", FieldName: "created_at"},
	"UpdatedAt":          {MySQLDatatype: "datetime", FieldName: "updated_at"},
	"DeletedAt":          {MySQLDatatype: "datetime", FieldName: "deleted_at"},
	"RetryCount":         {MySQLDatatype: "int", FieldName: "retry_count"},
	"NextRetryAt":        {MySQLDatatype: "datetime", FieldName: "next_retry_at"},
}
