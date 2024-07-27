package auth

import (
	"appname/customerrors"
	"appname/ent"
	"appname/ent/user"
	"appname/svc/verificationcode"
	"appname/utils"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type AuthSvc interface {
	Login(ctx context.Context, form LoginForm) (*LoginResult, error)
	SignUp(ctx context.Context, form SignUpForm) (*ent.User, error)
}

type LoginResult struct {
	HashedCsrfToken string
	JWTToken        string
}

type AuthSvcImpl struct {
	DB                  *ent.Client
	VerificationCodeSvc verificationcode.VerificationCodeSvc
}

func checkPassword(hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

func (s *AuthSvcImpl) Login(ctx context.Context, form LoginForm) (*LoginResult, error) {
	if form.Username == "" || form.Password == "" {
		return nil, customerrors.EmptyFormFieldsError{}
	}
	loginUser, err := s.DB.User.Query().Where(user.UsernameEQ(form.Username)).Only(ctx)
	if err != nil {
		return nil, customerrors.InvalidCredentialsError{}
	}

	if !checkPassword(loginUser.Password, form.Password) {
		return nil, customerrors.InvalidCredentialsError{}
	}
	if !loginUser.Validated {
		return nil, customerrors.NotValidatedAccountError{}
	}
	csrfToken, err := utils.GenerateRandomToken(64)
	if err != nil {
		return nil, err
	}
	// hash csrf
	hashedCsrfToken, err := utils.HashAndSalt(csrfToken)
	if err != nil {
		return nil, err
	}

	tokenString, err := utils.GenerateJWT(loginUser.ID.String(), loginUser.Email, loginUser.Username, csrfToken)
	if err != nil {
		return nil, err
	}
	return &LoginResult{HashedCsrfToken: hashedCsrfToken, JWTToken: tokenString}, nil
}

func (s *AuthSvcImpl) SignUp(ctx context.Context, form SignUpForm) (createdUser *ent.User, err error) {
	if form.Username == "" || form.Email == "" || form.Password == "" {
		return nil, customerrors.EmptyFormFieldsError{}
	}
	bytesPass, err := bcrypt.GenerateFromPassword([]byte(form.Password), 14)

	if err != nil {
		return nil, err
	}

	alreadyExistUser, _ := s.DB.User.Query().Where(user.UsernameEQ(form.Username)).First(ctx)
	if alreadyExistUser != nil {
		return nil, customerrors.UsernameAlreadyInUseError{}
	}
	alreadyExistMail, _ := s.DB.User.Query().Where(user.EmailEQ(form.Email)).First(ctx)
	if alreadyExistMail != nil {
		return nil, customerrors.EmailAlreadyInUseError{}
	}
	clientTx, err := s.DB.Tx(ctx)
	if err != nil {
		return nil, err
	}
	user, err := clientTx.User.Create().SetUsername(form.Username).SetPassword(string(bytesPass[:])).SetEmail(form.Email).Save(ctx)
	if err != nil {
		return nil, err
	}

	err = s.VerificationCodeSvc.Create(ctx, verificationcode.CreateForm{Username: user.Username, Type: utils.VALIDATION_TYPE}, clientTx)
	if err != nil {
		clientTx.Rollback()
		return nil, err
	}
	err = clientTx.Commit()
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}
