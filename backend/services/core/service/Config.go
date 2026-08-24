package service

type Config struct {
	IsEmailSendingEnabled       string `yaml:"isEmailSendingEnabled"`
	SmtpHost                    string `yaml:"smtpHost"`
	SmtpPort                    string `yaml:"smtpPort"`
	SmtpEmail                   string `yaml:"smtpEmail"`
	SmtpEncryptionType          string `yaml:"smtpEncryptionType"`
	SmtpPassword                string `yaml:"smtpPassword"`
	SmtpUsername                string `yaml:"smtpUsername"`
	SenderEmail                 string `yaml:"senderEmail"`
	SenderName                  string `yaml:"senderName"`
	TemplatePath                string `yaml:"templatePath"`
	WelcomeUserTemplateConstant string `yaml:"welcomeUserTemplateConstant"`
	WelcomeUserEmailSubject     string `yaml:"welcomeUserEmailSubject"`
}
