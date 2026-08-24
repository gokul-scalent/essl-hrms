package validation

import "strings"

func GetEmailProviderFromMX(mxRecord string) string {
	mxRecord = strings.ToLower(mxRecord)

	switch {
	// Google
	case strings.Contains(mxRecord, "google"),
		strings.Contains(mxRecord, "googlemail"),
		strings.Contains(mxRecord, "aspmx"):
		return "Google"

	// Microsoft
	case strings.Contains(mxRecord, "outlook"),
		strings.Contains(mxRecord, "protection.outlook"),
		strings.Contains(mxRecord, "office365"),
		strings.Contains(mxRecord, "hotmail"),
		strings.Contains(mxRecord, "live.com"):
		return "Microsoft"

	// Zoho (customizable provider)
	case strings.Contains(mxRecord, "zoho"):
		return "Zoho"

	// Yahoo
	case strings.Contains(mxRecord, "yahoo"),
		strings.Contains(mxRecord, "ymail"),
		strings.Contains(mxRecord, "yahoodns"):
		return "Yahoo"

	// Apple
	case strings.Contains(mxRecord, "icloud"),
		strings.Contains(mxRecord, "me.com"),
		strings.Contains(mxRecord, "mac.com"):
		return "Apple iCloud"

	// Rackspace
	case strings.Contains(mxRecord, "emailsrvr"):
		return "Rackspace"

	// Proton
	case strings.Contains(mxRecord, "protonmail"),
		strings.Contains(mxRecord, "proton"):
		return "Proton Mail"

	// Fastmail
	case strings.Contains(mxRecord, "messagingengine"),
		strings.Contains(mxRecord, "fastmail"):
		return "Fastmail"

	// Amazon SES / WorkMail
	case strings.Contains(mxRecord, "amazonses"),
		strings.Contains(mxRecord, "awsapps"):
		return "Amazon"

	// Mailgun
	case strings.Contains(mxRecord, "mailgun"):
		return "Mailgun"

	// SendGrid
	case strings.Contains(mxRecord, "sendgrid"):
		return "SendGrid"

	// SparkPost
	case strings.Contains(mxRecord, "sparkpost"):
		return "SparkPost"

	// Mimecast
	case strings.Contains(mxRecord, "mimecast"):
		return "Mimecast"

	// Proofpoint
	case strings.Contains(mxRecord, "pphosted"):
		return "Proofpoint"

	// GoDaddy
	case strings.Contains(mxRecord, "secureserver"):
		return "GoDaddy"

	// OVH
	case strings.Contains(mxRecord, "ovh"):
		return "OVH"

	// Yandex (regional and international provider)
	case strings.Contains(mxRecord, "yandex"):
		return "Yandex"

	// Tencent (regional and international provider)
	case strings.Contains(mxRecord, "qq.com"),
		strings.Contains(mxRecord, "exmail"):
		return "Tencent"

	default:
		return "Unknown"
	}
}
