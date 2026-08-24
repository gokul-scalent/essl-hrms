package email

import (
	"context"
)

type EmailConfig struct {
	ApiKey         string
	ServerUrl      string
	Port           string
	AccessKey      string
	SecretKey      string
	EncryptionType string
}

type EmailPayload struct {
	SenderName          string
	SenderEmail         string
	ToName              string
	ToEmail             string
	EmailSubject        string
	HtmlContent         string
	EmailBody           string
	Location            string
	StartAt             string
	EndAt               string
	TimeZone            string
	CampaignCode        string
	AttendeeID          int
	CalendarTitle       string
	CalendarDescription string
	InvoicePDF          string // invoice pdf
	CC                  []string
	BCC                 []string
	BookCalendar        bool
	MicrosoftID         string
}

type EmailSender interface {
	Send(ctx context.Context, payload EmailPayload, icsContent string) (int, string, error)
}
