package test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
	"vocablo/customerrors"
	"vocablo/ent"
	entuser "vocablo/ent/user"
	"vocablo/ent/verificationcode"
	"vocablo/svc/user"
	"vocablo/utils"

	"entgo.io/ent/dialect/sql"
	"github.com/stretchr/testify/assert"
)

func TestLoginOk(t *testing.T) {
	_, teardown, _ := SetupTest(t, true, nil)
	defer teardown(t)
	body, _ := json.Marshal(&testUserForm1)
	resp := testEnv.MakeRequest("POST", "/api/public/login", utils.GetStringPointer(string(body)))
	assert.True(t, strings.Contains(resp.Header().Get("Set-Cookie"), "JWT_TOKEN"))
	assert.Equal(t, 200, resp.Code, "Response status should be 200")
}

func TestLoginIncorrectCredentials(t *testing.T) {
	_, teardown, _ := SetupTest(t, true, nil)
	defer teardown(t)
	body, _ := json.Marshal(&testUserForm2)
	resp := testEnv.MakeRequest("POST", "/api/public/login", utils.GetStringPointer(string(body)))
	assert.Equal(t, 401, resp.Code, "Response status should be 401")
}

func TestLoginEmptyForm(t *testing.T) {
	_, teardown, _ := SetupTest(t, true, nil)
	defer teardown(t)
	body, _ := json.Marshal(&user.CreateForm{})
	resp := testEnv.MakeRequest("POST", "/api/public/login", utils.GetStringPointer(string(body)))
	assert.Equal(t, 400, resp.Code, "Response status should be 400")
}

func TestLoginNotValidatedAccount(t *testing.T) {
	_, teardown, _ := SetupTest(t, false, nil)
	defer teardown(t)
	body, _ := json.Marshal(&testUserForm1)
	resp := testEnv.MakeRequest("POST", "/api/public/login", utils.GetStringPointer(string(body)))
	assert.Equal(t, 403, resp.Code, "Response status should be 403")
}

func TestSignupOk(t *testing.T) {
	client, teardown := StartTest(t)
	defer teardown(t)
	body, _ := json.Marshal(&testUserForm2)
	resp := testEnv.MakeRequest("POST", "/api/public/register", utils.GetStringPointer(string(body)))
	assert.Equal(t, 200, resp.Code, "Response status should be 200")
	//We check that the verificationcode was created
	existValidationCode, err := client.VerificationCode.Query().Where(verificationcode.And(
		verificationcode.HasUserWith(entuser.UsernameEQ(testUserForm2.Username)),
		verificationcode.TypeEQ(utils.VALIDATION_TYPE))).Exist(context.Background())
	if err != nil {
		t.Errorf("Error checking validation code: %s", err)
	}
	assert.True(t, existValidationCode, "Validation code should have been created")
}

func TestSignupEmptyForm(t *testing.T) {
	_, teardown := StartTest(t)
	defer teardown(t)
	body, _ := json.Marshal(&user.CreateForm{})
	resp := testEnv.MakeRequest("POST", "/api/public/register", utils.GetStringPointer(string(body)))
	assert.Equal(t, 400, resp.Code, "Response status should be 400")
}

func TestSignupMailAlreadyInUse(t *testing.T) {
	_, teardown, _ := SetupTest(t, true, nil)
	defer teardown(t)
	form := testUserForm2
	form.Email = testUserForm1.Email
	body, _ := json.Marshal(&form)
	resp := testEnv.MakeRequest("POST", "/api/public/register", utils.GetStringPointer(string(body)))
	assert.Equal(t, 409, resp.Code, "Response status should be 409")
	bodyRes := resp.Body.Bytes()
	var bodyResObj utils.ResponseBody
	err := json.Unmarshal(bodyRes, &bodyResObj)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, customerrors.EMAIL_ALREADY_IN_USE, *bodyResObj.ErrorCode, "Error code should be EMAIL_ALREADY_IN_USE")
}

func TestSignupUsernameAlreadyInUse(t *testing.T) {
	_, teardown, _ := SetupTest(t, true, nil)
	defer teardown(t)
	form := testUserForm2
	form.Username = testUserForm1.Username
	body, _ := json.Marshal(&form)
	resp := testEnv.MakeRequest("POST", "/api/public/register", utils.GetStringPointer(string(body)))
	assert.Equal(t, 409, resp.Code, "Response status should be 409")
	bodyRes := resp.Body.Bytes()
	var bodyResObj utils.ResponseBody
	err := json.Unmarshal(bodyRes, &bodyResObj)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, customerrors.USERNAME_ALREADY_IN_USE, *bodyResObj.ErrorCode, "Error code should be USERNAME_ALREADY_IN_USE")
}

func TestValidateAccountOk(t *testing.T) {
	client, teardown, ctx := SetupTest(t, false, nil)
	defer teardown(t)
	expirationDate := time.Now().Add(time.Minute * 15)
	verificationCode, err := client.VerificationCode.Create().SetCode(verificationCodeForm.Code).SetUserID(client.User.Query().Where(
		entuser.UsernameEQ(testUserForm1.Username)).OnlyX(ctx).ID).
		SetType(utils.VALIDATION_TYPE).SetExpireDate(expirationDate).Save(ctx)
	if err != nil {
		t.Errorf("Error creating verification code: %s", err)
	}
	resp := testEnv.MakeRequest("POST", "/api/public/validate/"+testUserForm1.Username+"/"+verificationCode.Code, nil)

	assert.Equal(t, 200, resp.Code, "Response status should be 200")
	//We check that the user is validated
	user, err := client.User.Query().Where(entuser.Username(testUserForm1.Username)).Only(ctx)
	if err != nil {
		t.Errorf("Error getting user: %s", err)
	}
	assert.True(t, user.Validated, "User should be validated")
}

func TestValidateAccountWrongCode(t *testing.T) {
	client, teardown, ctx := SetupTest(t, false, nil)
	defer teardown(t)
	expirationDate := time.Now().Add(time.Minute * 15)
	_, err := client.VerificationCode.Create().SetCode(verificationCodeForm.Code).SetUserID(client.User.Query().Where(
		entuser.UsernameEQ(testUserForm1.Username)).OnlyX(ctx).ID).
		SetType(utils.VALIDATION_TYPE).SetExpireDate(expirationDate).Save(ctx)
	if err != nil {
		t.Errorf("Error creating verification code: %s", err)
	}
	resp := testEnv.MakeRequest("POST", "/api/public/validate/"+testUserForm1.Username+"/"+INCORRECT_VERIFICATION_CODE, nil)

	assert.Equal(t, 401, resp.Code, "Response status should be 401")
	bodyRes := resp.Body.Bytes()
	var bodyResObj utils.ResponseBody
	err = json.Unmarshal(bodyRes, &bodyResObj)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, customerrors.INCORRECT_VALIDATION_CODE, *bodyResObj.ErrorCode, "Error code should be INCORRECT_VALIDATION_CODE")
}

func TestValidateAccountAlreadyUsedCode(t *testing.T) {
	client, teardown, ctx := SetupTest(t, false, nil)
	defer teardown(t)
	expirationDate := time.Now().Add(time.Minute * 15)
	verificationCode, err := client.VerificationCode.Create().SetCode(verificationCodeForm.Code).SetUserID(client.User.Query().Where(
		entuser.UsernameEQ(testUserForm1.Username)).OnlyX(ctx).ID).
		SetType(utils.VALIDATION_TYPE).SetExpireDate(expirationDate).Save(ctx)
	if err != nil {
		t.Errorf("Error creating verification code: %s", err)
	}
	testEnv.MakeRequest("POST", "/api/public/validate/"+testUserForm1.Username+"/"+verificationCode.Code, nil)
	resp := testEnv.MakeRequest("POST", "/api/public/validate/"+testUserForm1.Username+"/"+verificationCode.Code, nil)

	assert.Equal(t, 409, resp.Code, "Response status should be 409")
	bodyRes := resp.Body.Bytes()
	var bodyResObj utils.ResponseBody
	err = json.Unmarshal(bodyRes, &bodyResObj)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, customerrors.ALREADY_USED_VALIDATION_CODE, *bodyResObj.ErrorCode, "Error code should be ALREADY_USED_VALIDATION_CODE")
}

func TestValidateAccountExpiredCode(t *testing.T) {
	client, teardown, ctx := SetupTest(t, false, nil)
	defer teardown(t)
	expirationDate := time.Now().Add(time.Minute * -5)
	verificationCode, err := client.VerificationCode.Create().SetCode(verificationCodeForm.Code).SetUserID(client.User.Query().Where(
		entuser.UsernameEQ(testUserForm1.Username)).OnlyX(ctx).ID).
		SetType(utils.VALIDATION_TYPE).SetExpireDate(expirationDate).Save(ctx)
	if err != nil {
		t.Errorf("Error creating verification code: %s", err)
	}
	resp := testEnv.MakeRequest("POST", "/api/public/validate/"+testUserForm1.Username+"/"+verificationCode.Code, nil)

	assert.Equal(t, 410, resp.Code, "Response status should be 410")
	bodyRes := resp.Body.Bytes()
	var bodyResObj utils.ResponseBody
	err = json.Unmarshal(bodyRes, &bodyResObj)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, customerrors.EXPIRED_VALIDATION_CODE, *bodyResObj.ErrorCode, "Error code should be EXPIRED_VALIDATION_CODE")
}

func TestValidateAccountResendCode(t *testing.T) {
	client, teardown, ctx := SetupTest(t, false, nil)
	defer teardown(t)
	expirationDate := time.Now().Add(time.Minute * 15)
	//We create a verification code with a creation date in the past to ensure that the just created code is the one used
	verificationCode, err := client.VerificationCode.Create().SetCode(verificationCodeForm.Code).SetCreationDate(time.Now().Add(time.Minute * -5)).SetUserID(client.User.Query().Where(
		entuser.UsernameEQ(testUserForm1.Username)).OnlyX(ctx).ID).
		SetType(utils.VALIDATION_TYPE).SetExpireDate(expirationDate).Save(ctx)
	if err != nil {
		t.Errorf("Error creating verification code: %s", err)
	}

	resp := testEnv.MakeRequest("POST", "/api/public/validate/"+testUserForm1.Username+"/resend", nil)
	assert.Equal(t, 200, resp.Code, "Response status should be 200")

	//We check that the verification code was created
	validationCode, err := client.VerificationCode.Query().Where(verificationcode.And(verificationcode.HasUserWith(entuser.UsernameEQ(testUserForm1.Username)),
		verificationcode.TypeEQ(utils.VALIDATION_TYPE))).Order(verificationcode.ByCreationDate(sql.OrderDesc())).First(ctx)
	if err != nil {
		t.Errorf("Error getting verification code: %s", err)
	}
	//We check that the previous verification code is not working
	verif, _ := client.VerificationCode.Query().All(ctx)
	fmt.Println(verif)
	resp = testEnv.MakeRequest("POST", "/api/public/validate/"+testUserForm1.Username+"/"+verificationCode.Code, nil)
	assert.Equal(t, 401, resp.Code, "Response status should be 401")
	bodyRes := resp.Body.Bytes()
	var bodyResObj utils.ResponseBody
	err = json.Unmarshal(bodyRes, &bodyResObj)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, customerrors.INCORRECT_VALIDATION_CODE, *bodyResObj.ErrorCode, "Error code should be INCORRECT_VALIDATION_CODE")
	//We check that the new verification code is working
	resp = testEnv.MakeRequest("POST", "/api/public/validate/"+testUserForm1.Username+"/"+validationCode.Code, nil)
	assert.Equal(t, 200, resp.Code, "Response status should be 200")
	//We check that the user is validated
	user, err := client.User.Query().Where(entuser.Username(testUserForm1.Username)).Only(ctx)
	if err != nil {
		t.Errorf("Error getting user: %s", err)
	}
	assert.True(t, user.Validated, "User should be validated")
}

func TestResetPasswordOk(t *testing.T) {
	client, teardown, ctx := SetupTest(t, true, nil)
	defer teardown(t)
	resp := testEnv.MakeRequest("POST", "/api/public/forgotten-password/"+testUserForm1.Username, nil)
	assert.Equal(t, 200, resp.Code, "Response status should be 200")

	//We check that the verification code was created and we get the code number
	validationCode, err := client.VerificationCode.Query().Where(verificationcode.And(verificationcode.HasUserWith(entuser.UsernameEQ(testUserForm1.Username)),
		verificationcode.TypeEQ(utils.RESET_TYPE))).Order(verificationcode.ByCreationDate(sql.OrderDesc())).First(ctx)
	if err != nil {
		t.Errorf("Error getting verification code: %s", err)
	}
	//We reset the password using the code
	resp = testEnv.MakeRequest("POST", "/api/public/reset-password/"+testUserForm1.Username+"/"+validationCode.Code, utils.GetStringPointer(NEW_PASSWORD))
	assert.Equal(t, 200, resp.Code, "Response status should be 200")

	//We check that the password was updated
	newPassTestUserForm := testUserForm1
	newPassTestUserForm.Password = NEW_PASSWORD
	//First, we check that the previous password is not working
	body, err := json.Marshal(&testUserForm1)
	if err != nil {
		t.Errorf("Error marshalling body: %s", err)
	}
	resp = testEnv.MakeRequest("POST", "/api/public/login", utils.GetStringPointer(string(body)))
	assert.Equal(t, 401, resp.Code, "Response status should be 401")
	//Then, we check that the new password is working
	body, err = json.Marshal(&newPassTestUserForm)
	if err != nil {
		t.Errorf("Error marshalling body: %s", err)
	}
	resp = testEnv.MakeRequest("POST", "/api/public/login", utils.GetStringPointer(string(body)))
	assert.Equal(t, 200, resp.Code, "Response status should be 200")
}

func TestSelf(t *testing.T) {
	_, teardown, ctx := SetupTest(t, true, nil)
	defer teardown(t)
	body, _ := json.Marshal(&testUserForm1)
	resp := testEnv.MakeRequest("POST", "/api/public/login", utils.GetStringPointer(string(body)))

	bodyRes := resp.Body.Bytes()
	var bodyResObj utils.ResponseBody
	err := json.Unmarshal(bodyRes, &bodyResObj)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	resp = testEnv.MakeAuthRequest("GET", "/api/self", nil, ctx)
	assert.Equal(t, 200, resp.Code, "Response status should be 200")
	bodyRes = resp.Body.Bytes()
	err = json.Unmarshal(bodyRes, &bodyResObj)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	var resultUser ent.User
	userStr, err := json.Marshal(bodyResObj.Data)
	if err != nil {
		t.Errorf("Error marshalling response body: %s", err)
	}
	err = json.Unmarshal(userStr, &resultUser)
	if err != nil {
		t.Errorf("Error unmarshalling response body: %s", err)
	}
	assert.Equal(t, testUserForm1.Username, resultUser.Username, "Username should be the same")

}
