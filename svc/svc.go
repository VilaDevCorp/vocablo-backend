package svc

import (
	"appname/ent"
	"appname/svc/auth"
	"appname/svc/mail"
	"appname/svc/user"
	"appname/svc/verificationcode"
)

type Service struct {
	// Activity         activity.Svc
	User             user.UserSvc
	Auth             auth.AuthSvc
	VerificationCode verificationcode.VerificationCodeSvc
}

var svc Service

func Get() *Service {
	return &svc
}

func Setup(client *ent.Client, mailSvc mail.MailSvc) {
	svc = Service{
		// Activity:         &activity.Store{DB: client},
		User:             &user.UserSvcImpl{DB: client},
		Auth:             &auth.AuthSvcImpl{DB: client, VerificationCodeSvc: &verificationcode.VerificationCodeSvcImpl{DB: client, Mail: mailSvc}},
		VerificationCode: &verificationcode.VerificationCodeSvcImpl{DB: client, Mail: mailSvc},
	}
}
