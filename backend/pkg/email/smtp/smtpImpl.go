package smtp

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"net/smtp"
	"strings"

	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/email"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
)

type SmtpEmailPayload struct {
	FromEmail string
	ToEmail   []string
	Subject   string
	EmailBody string
}

type Smtp struct {
	Config   *email.EmailConfig
	SmtpAuth smtp.Auth
}

func NewSmtpClient(config *email.EmailConfig) *Smtp {
	auth := smtp.PlainAuth("", config.AccessKey, config.SecretKey, config.ServerUrl)
	return &Smtp{
		Config:   config,
		SmtpAuth: auth,
	}
}

func (s *Smtp) ValidateSMTPConfig(ctx context.Context) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>pkg>email>smtp: ValidateSMTPConfig started", reqID)

	// smtpInfo := config.SmtpValueJson

	// Setup authentication
	auth := smtp.PlainAuth("", s.Config.AccessKey, s.Config.SecretKey, s.Config.ServerUrl)

	// Determine the connection method based on the encryption type
	serverAddress := fmt.Sprintf("%s:%s", s.Config.ServerUrl, s.Config.Port)

	switch strings.ToLower(s.Config.EncryptionType) {
	case "tls/ssl":
		// Establish an SSL/TLS connection
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true, // Disable certificate verification (for testing purposes)
			ServerName:         s.Config.ServerUrl,
		}

		conn, err := tls.Dial("tcp", serverAddress, tlsConfig)
		if err != nil {
			log.Error(err.Error(), reqID)
			return errors.ResponseBadRequestError("wrong serverUrl or port")
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, s.Config.ServerUrl)
		if err != nil {
			log.Error(err.Error(), reqID)
			return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
		}
		defer client.Quit()

		// Authenticate with the SMTP server
		err = client.Auth(auth)
		if err != nil {
			log.Error(err.Error(), reqID)
			return errors.ResponseBadRequestError("invalid accessKey or secretKey")
		}

	case "starttls":
		// Connect without SSL/TLS
		client, err := smtp.Dial(serverAddress)
		if err != nil {
			log.Error(err.Error(), reqID)
			return errors.ResponseBadRequestError("wrong serverUrl or port")
		}
		defer client.Quit()

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         s.Config.ServerUrl,
		}

		err = client.StartTLS(tlsConfig)
		if err != nil {
			log.Error(err.Error(), reqID)
			return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
		}

		// Authenticate with the SMTP server
		err = client.Auth(auth)
		if err != nil {
			log.Error(err.Error(), reqID)
			return errors.ResponseBadRequestError("invalid accessKey or secretKey")
		}
	}

	log.Info("core>pkg>email>smtp: ValidateSMTPConfig completed", reqID)
	return nil
}

func (s *Smtp) Send(ctx context.Context, payload email.EmailPayload, icsFileContent string) (int, string, error) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>pkg>email>smtp: Send using SMTP started", reqID)

	go func() {
		// Create a boundary for separating parts of the email
		boundary := "my-boundary-12345"

		//add all recicpients into single array
		recipients := []string{payload.ToEmail}
		recipients = append(recipients, payload.CC...)
		recipients = append(recipients, payload.BCC...)

		// Create the email headers
		headers := make(map[string]string)
		headers["From"] = fmt.Sprintf("%s <%s>", payload.SenderName, payload.SenderEmail)
		headers["To"] = payload.ToEmail

		if len(payload.CC) > 0 {
			headers["Cc"] = strings.Join(payload.CC, ", ")
		}

		headers["Subject"] = mime.BEncoding.Encode("UTF-8", payload.EmailSubject)
		headers["MIME-Version"] = "1.0"
		headers["Content-Type"] = fmt.Sprintf("multipart/mixed; boundary=%s", boundary)

		// Construct the email body
		var messageBody strings.Builder
		for key, value := range headers {
			messageBody.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
		}
		messageBody.WriteString("\r\n--" + boundary + "\r\n")
		messageBody.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
		messageBody.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
		messageBody.WriteString(payload.HtmlContent + "\r\n")

		// Attach the .ics file if content is provided
		if icsFileContent != "" {
			encodedIcs := base64.StdEncoding.EncodeToString([]byte(icsFileContent))
			messageBody.WriteString("--" + boundary + "\r\n")
			messageBody.WriteString("Content-Type: text/calendar; name=\"invite.ics\"\r\n")
			messageBody.WriteString("Content-Disposition: attachment; filename=\"invite.ics\"\r\n")
			messageBody.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
			messageBody.WriteString(encodedIcs + "\r\n")
		}

		// Attach the PDF invoice if content is provided
		if payload.InvoicePDF != "" {
			encodedInvoice := base64.StdEncoding.EncodeToString([]byte(payload.InvoicePDF))
			messageBody.WriteString("--" + boundary + "\r\n")
			messageBody.WriteString("Content-Type: application/pdf; name=\"invoice.pdf\"\r\n")
			messageBody.WriteString("Content-Disposition: attachment; filename=\"invoice.pdf\"\r\n")
			messageBody.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
			messageBody.WriteString(encodedInvoice + "\r\n")
		}

		// Close the MIME message with the boundary
		messageBody.WriteString("--" + boundary + "--")

		// Connect to the SMTP server.
		err := smtp.SendMail(s.Config.ServerUrl+":"+s.Config.Port, s.SmtpAuth, payload.SenderEmail, recipients, []byte(messageBody.String()))
		if err != nil {
			log.Error(err.Error(), reqID)
			if strings.Contains(err.Error(), "451") {
				log.Error(err.Error(), reqID)
				return
			}
			return
		}
	}()

	log.Info("core>pkg>email>smtp: Send using SMTP completed", reqID)
	return http.StatusAccepted, "", nil
}
