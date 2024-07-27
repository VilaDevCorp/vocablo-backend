package mail

import (
	"appname/conf"
	"net/smtp"
)

type MailSvc interface {
	SendMail(to string, subject string, message string) error
}

type MailSvcImpl struct{}

func (*MailSvcImpl) SendMail(to string, subject string, message string) error {
	conf := conf.Get()

	auth := smtp.PlainAuth("", conf.Mail.User, conf.Mail.Pass, conf.Mail.SmtpHost)
	toArray := []string{to}
	msg := []byte("From: " + conf.Mail.User + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		message + "\r\n")

	err := smtp.SendMail(conf.Mail.SmtpHost+":"+conf.Mail.SmtpPort, auth, conf.Mail.User, toArray, msg)
	if err != nil {
		return err
	}
	return nil
}
