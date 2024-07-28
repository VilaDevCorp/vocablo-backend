package auth

import (
	"net/http"
	"vocablo/conf"
	"vocablo/customerrors"
	"vocablo/svc"
	"vocablo/svc/auth"
	"vocablo/svc/verificationcode"
	"vocablo/utils"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	conf := conf.Get()
	var form auth.LoginForm
	err := c.ShouldBind(&form)
	if err != nil {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer(err.Error()), nil)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}
	svc := svc.Get()
	loginResponse, err := svc.Auth.Login(c.Request.Context(), form)
	var res utils.HttpResponse
	if err != nil {
		switch err.(type) {
		case customerrors.EmptyFormFieldsError:
			res = utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer("No username or password present"), nil)
		case customerrors.InvalidCredentialsError:
			res = utils.ErrorResponse(http.StatusUnauthorized, utils.GetStringPointer("Invalid credentials"), utils.GetStringPointer(customerrors.INVALID_CREDENTIALS))
		case customerrors.NotValidatedAccountError:
			res = utils.ErrorResponse(http.StatusForbidden, utils.GetStringPointer("Not validated account"), utils.GetStringPointer(customerrors.NOT_VALIDATED_ACCOUNT))
		default:
			res = utils.InternalError(err)
		}
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
		if conf.Env == "prod" {
			c.SetCookie("JWT_TOKEN", loginResponse.JWTToken, 30*24*60*60*1000, "/", conf.Prod.CookieHost, conf.Prod.CookieSecure, conf.Prod.CookieHttpOnly)
		} else {
			c.SetCookie("JWT_TOKEN", loginResponse.JWTToken, 30*24*60*60*1000, "/", conf.Dev.CookieHost, conf.Dev.CookieSecure, conf.Dev.CookieHttpOnly)
		}
		res = utils.SuccessResponse(loginResponse.HashedCsrfToken)
	}
	c.JSON(res.Status, res.Body)
}

func SignUp(c *gin.Context) {
	var form auth.SignUpForm
	err := c.ShouldBind(&form)
	if err != nil {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer(err.Error()), nil)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}
	svc := svc.Get()
	user, err := svc.Auth.SignUp(c.Request.Context(), form)
	var res utils.HttpResponse
	if err != nil {
		switch err.(type) {
		case customerrors.EmptyFormFieldsError:
			res = utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer("No username, email or password present"), nil)
		case customerrors.UsernameAlreadyInUseError:
			res = utils.ErrorResponse(http.StatusConflict, utils.GetStringPointer("Username already exist"), utils.GetStringPointer(customerrors.USERNAME_ALREADY_IN_USE))
		case customerrors.EmailAlreadyInUseError:
			res = utils.ErrorResponse(http.StatusConflict, utils.GetStringPointer("Email already exist"), utils.GetStringPointer(customerrors.EMAIL_ALREADY_IN_USE))
		default:
			res = utils.InternalError(err)
		}
	} else {
		res = utils.SuccessResponse(user)
	}
	c.JSON(res.Status, res.Body)
}

func ResendValidationCode(c *gin.Context) {
	username, _ := c.Params.Get("username")
	svc := svc.Get()
	err := svc.VerificationCode.Create(c.Request.Context(), verificationcode.CreateForm{Username: username, Type: utils.VALIDATION_TYPE}, nil)

	var res utils.HttpResponse
	if err != nil {
		res = utils.InternalError(err)
	} else {
		res = utils.SuccessResponse(nil)
	}
	c.JSON(res.Status, res.Body)
}

func SendForgottenPasswordCode(c *gin.Context) {
	username, _ := c.Params.Get("username")
	svc := svc.Get()
	err := svc.VerificationCode.Create(c.Request.Context(), verificationcode.CreateForm{Username: username, Type: utils.RESET_TYPE}, nil)
	var res utils.HttpResponse
	if err != nil {
		res = utils.InternalError(err)

	} else {
		res = utils.SuccessResponse(nil)
	}
	c.JSON(res.Status, res.Body)
}

func ValidateAccount(c *gin.Context) {
	username, _ := c.Params.Get("username")
	code, _ := c.Params.Get("code")
	form := verificationcode.UseForm{Username: username, Code: code, Type: utils.VALIDATION_TYPE}
	svc := svc.Get()
	err := svc.VerificationCode.UseCode(c.Request.Context(), form)
	var res utils.HttpResponse
	if err != nil {
		switch err.(type) {
		case customerrors.AlreadyUsedValidationCodeError:
			res = utils.ErrorResponse(http.StatusConflict, utils.GetStringPointer("Validation code already used"), utils.GetStringPointer(customerrors.ALREADY_USED_VALIDATION_CODE))
		case customerrors.ExpiredValidationCodeError:
			res = utils.ErrorResponse(http.StatusGone, utils.GetStringPointer("Validation code expired"), utils.GetStringPointer(customerrors.EXPIRED_VALIDATION_CODE))
		case customerrors.IncorrectValidationCodeError:
			res = utils.ErrorResponse(http.StatusUnauthorized, utils.GetStringPointer("Invalid validation code"), utils.GetStringPointer(customerrors.INCORRECT_VALIDATION_CODE))
		default:
			res = utils.InternalError(err)
		}
	} else {
		res = utils.SuccessResponse(nil)
	}
	c.JSON(res.Status, res.Body)
}

func ResetPassword(c *gin.Context) {
	username, _ := c.Params.Get("username")
	code, _ := c.Params.Get("code")
	body, err := c.GetRawData()
	if err != nil {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer(err.Error()), nil)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}
	newPassStr := string(body)

	if newPassStr == "" {
		res := utils.ErrorResponse(http.StatusBadRequest, utils.GetStringPointer("No new password present"), nil)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}

	form := verificationcode.UseForm{Username: username, Code: code, Type: utils.RESET_TYPE, NewPass: newPassStr}
	svc := svc.Get()
	err = svc.VerificationCode.UseCode(c.Request.Context(), form)
	var res utils.HttpResponse
	if err != nil {
		switch err.(type) {
		case customerrors.AlreadyUsedValidationCodeError:
			res = utils.ErrorResponse(http.StatusConflict, utils.GetStringPointer("Validation code already used"), utils.GetStringPointer(customerrors.ALREADY_USED_VALIDATION_CODE))
		case customerrors.ExpiredValidationCodeError:
			res = utils.ErrorResponse(http.StatusGone, utils.GetStringPointer("Validation code expired"), utils.GetStringPointer(customerrors.EXPIRED_VALIDATION_CODE))
		case customerrors.IncorrectValidationCodeError:
			res = utils.ErrorResponse(http.StatusUnauthorized, utils.GetStringPointer("Invalid validation code"), utils.GetStringPointer(customerrors.INCORRECT_VALIDATION_CODE))
		default:
			res = utils.InternalError(err)
		}
	} else {
		res = utils.SuccessResponse(nil)
	}
	c.JSON(res.Status, res.Body)
}

func Self(c *gin.Context) {
	jwtCookie, err := c.Cookie("JWT_TOKEN")
	if err != nil {
		res := utils.InternalError(err)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}
	claims, err := utils.ValidateToken(jwtCookie)
	if err != nil {
		res := utils.InternalError(err)
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}
	svc := svc.Get()
	user, err := svc.User.Get(c.Request.Context(), claims.Id)
	if err != nil {
		res := utils.NotFound("user", claims.Id.String())
		c.AbortWithStatusJSON(res.Status, res.Body)
		return
	}
	res := utils.SuccessResponse(user)
	c.JSON(res.Status, res.Body)
}
