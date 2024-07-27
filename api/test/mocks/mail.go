package mocks

import (
	"fmt"
)

type MailSvcMock struct{}

func (s *MailSvcMock) SendMail(to string, subject string, message string) error {
	fmt.Println("Sending mail to: ", to, " with subject: ", subject, " and message: ", message)
	return nil
}
