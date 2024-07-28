package test

import (
	"vocablo/ent"
	"vocablo/schema"
	"vocablo/svc/user"
	"vocablo/svc/userword"
	"vocablo/utils"
)

var testUserForm1 user.CreateForm = user.CreateForm{Username: "test", Password: "test", Email: "test@gmail.com"}
var testUserForm2 user.CreateForm = user.CreateForm{Username: "test2", Password: "test2", Email: "test2@gmail.com"}
var verificationCodeForm *ent.VerificationCode = &ent.VerificationCode{Code: "123456", Type: utils.VALIDATION_TYPE}

const INCORRECT_VERIFICATION_CODE string = "654321"
const NEW_PASSWORD string = "newpassword"

var testWordForm1 = userword.CreateForm{Term: "bad", Lang: "en",
	Definitions: []schema.Definition{{Definition: "not good", Example: "drug is bad"}}}
var testWordForm2 = userword.CreateForm{Term: "good", Lang: "en",
	Definitions: []schema.Definition{{Definition: "not bad", Example: "vegetables are good"}}}
var otherUserWordForm = userword.CreateForm{Term: "regular", Lang: "en",
	Definitions: []schema.Definition{{Definition: "normal", Example: "regular exercise"}}}

const UPDATED_TERM = "adverse"
