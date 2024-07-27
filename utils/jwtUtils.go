package utils

import (
	"appname/conf"
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JWTClaim struct {
	Id       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Username string    `json:"username"`
	Mail     string    `json:"mail"`
	Csrf     string    `json:"csrf"`
	jwt.RegisteredClaims
}

func GenerateJWT(id string, mail string, username string, csrf string) (tokenString string, err error) {
	conf := conf.Get()
	expirationTime := time.Now().Add(24 * time.Hour)
	fmt.Println(expirationTime)
	fmt.Println(jwt.NewNumericDate(expirationTime))
	uuidId, _ := uuid.Parse(id)
	issuer := conf.Dev.CookieHost
	if conf.Env == "prod" {
		issuer = conf.Prod.CookieHost
	}
	claims := &JWTClaim{
		Id:       uuidId,
		Mail:     mail,
		Username: username,
		Csrf:     csrf,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(conf.JwtKey))
	return
}
func ValidateToken(signedToken string) (tokenClaims JWTClaim, err error) {
	conf := conf.Get()

	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.JwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if int64(claims.ExpiresAt.Compare(time.Now().Local())) <= 0 {
		err = errors.New("token expired")
		return
	}
	return *claims, nil
}
