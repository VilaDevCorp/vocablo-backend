package svc

import (
	"vocablo/ent"
	"vocablo/svc/auth"
	"vocablo/svc/mail"
	"vocablo/svc/user"
	"vocablo/svc/userword"
	"vocablo/svc/verificationcode"
	"vocablo/svc/word"
)

type Service struct {
	User             user.UserSvc
	Auth             auth.AuthSvc
	VerificationCode verificationcode.VerificationCodeSvc
	UserWord         userword.UserWordSvc
	Word             word.WordSvc
}

var svc Service

func Get() *Service {
	return &svc
}

func Setup(client *ent.Client, mailSvc mail.MailSvc) {
	svc = Service{
		User:             &user.UserSvcImpl{DB: client},
		Auth:             &auth.AuthSvcImpl{DB: client, VerificationCodeSvc: &verificationcode.VerificationCodeSvcImpl{DB: client, Mail: mailSvc}},
		VerificationCode: &verificationcode.VerificationCodeSvcImpl{DB: client, Mail: mailSvc},
		UserWord:         &userword.UserWordSvcImpl{DB: client},
		Word:             &word.WordSvcImpl{DB: client},
	}
}
