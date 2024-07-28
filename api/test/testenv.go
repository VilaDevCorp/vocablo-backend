package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"vocablo/api"
	"vocablo/api/test/mocks"
	"vocablo/ent"
	"vocablo/ent/enttest"
	"vocablo/svc"
	"vocablo/svc/auth"
	"vocablo/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
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

// General function to setup a test, that creates a user validated or not and accepts a custom setup function
// for more specific setup
func SetupTest(t *testing.T, isValidated bool, customSetup func(client *ent.Client, t *testing.T, ctx context.Context)) (*ent.Client,
	func(t *testing.T), context.Context) {
	ctx := context.Background()
	client, teardown := StartTest(t)
	bytesPass, err := bcrypt.GenerateFromPassword([]byte(testUserForm.Password), 14)
	if err != nil {
		t.Errorf("Error hashing password: %s", err)
	}
	_, err = client.User.Create().SetUsername(testUserForm.Username).SetEmail(testUserForm.Email).SetPassword(string(bytesPass[:])).SetValidated(isValidated).Save(context.Background())
	if err != nil {
		t.Errorf("Error creating user: %s", err)
	}
	if isValidated {
		loginResult, err := svc.Get().Auth.Login(ctx, auth.LoginForm{Username: testUserForm.Username, Password: testUserForm.Password})
		if err != nil {
			t.Errorf("Error logging in user: %s", err)
		}
		ctx = context.WithValue(ctx, utils.CsrfKey, loginResult.HashedCsrfToken)
		ctx = context.WithValue(ctx, utils.JwtKey, loginResult.JWTToken)
	}
	if customSetup != nil {
		customSetup(client, t, ctx)
	}
	return client, teardown, ctx
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

func (testEnv TestEnvironment) MakeAuthRequest(method string, path string, body *string,
	ctx context.Context) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	var bodyReader *strings.Reader
	if body == nil {
		bodyReader = strings.NewReader("")
	} else {
		bodyReader = strings.NewReader(*body)
	}
	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-CSRF", ctx.Value(utils.CsrfKey).(string))
	req.AddCookie(&http.Cookie{Name: "JWT_TOKEN", Value: ctx.Value(utils.JwtKey).(string), Path: "/", HttpOnly: true, Secure: false, SameSite: http.SameSiteLaxMode})
	testEnv.Router.ServeHTTP(recorder, req)
	return recorder
}
