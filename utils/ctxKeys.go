package utils

type CtxKey string

const (
	UserIdKey CtxKey = "userID"
	CsrfKey   CtxKey = "csrf"
	JwtKey    CtxKey = "jwtToken"
)
