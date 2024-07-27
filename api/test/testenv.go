package test

import (
	"appname/api"
	"appname/api/test/mocks"
	"appname/ent"
	"appname/ent/enttest"
	"appname/svc"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

type TestEnvironment struct {
	Router *gin.Engine
}

var testEnv TestEnvironment

func Get() *TestEnvironment {
	return &testEnv
}

func StartTest(t *testing.T) (*ent.Client, func(t *testing.T)) {
	log.Info().Msg("Setup tests ")
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")

	svc.Setup(client, &(mocks.MailSvcMock{}))
	testEnv.Router = api.GetRouter()

	// Return a function to teardown the test
	return client, func(t *testing.T) {
		log.Info().Msg("Teardown tests ")
		defer client.Close()
	}
}

func (testEnv TestEnvironment) MakeRequest(method string, path string, body *string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	var bodyReader *strings.Reader
	if body == nil {
		bodyReader = strings.NewReader("")
	} else {
		bodyReader = strings.NewReader(*body)
	}
	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	testEnv.Router.ServeHTTP(recorder, req)
	return recorder
}

func (testEnv TestEnvironment) MakeAuthRequest(method string, path string, body *string, csrf *string, jwtToken *string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	var bodyReader *strings.Reader
	if body == nil {
		bodyReader = strings.NewReader("")
	} else {
		bodyReader = strings.NewReader(*body)
	}
	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-CSRF", *csrf)
	req.AddCookie(&http.Cookie{Name: "JWT_TOKEN", Value: *jwtToken, Path: "/", HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode})
	testEnv.Router.ServeHTTP(recorder, req)
	return recorder
}
