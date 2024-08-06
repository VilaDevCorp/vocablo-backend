package customerrors

const (
	NOT_CSRF_TOKEN               = "NOT_CSRF_TOKEN"
	NOT_JWT_TOKEN                = "NOT_JWT_TOKEN"
	INVALID_TOKEN                = "INVALID_TOKEN"
	USERNAME_ALREADY_IN_USE      = "USERNAME_ALREADY_IN_USE"
	EMAIL_ALREADY_IN_USE         = "EMAIL_ALREADY_IN_USE"
	ALREADY_VALIDATED_ACCOUNT    = "ALREADY_VALIDATED_ACCOUNT"
	INVALID_CREDENTIALS          = "INVALID_CREDENTIALS"
	NOT_VALIDATED_ACCOUNT        = "NOT_VALIDATED_ACCOUNT"
	EXPIRED_VALIDATION_CODE      = "EXPIRED_VALIDATION_CODE"
	INCORRECT_VALIDATION_CODE    = "INCORRECT_VALIDATION_CODE"
	ALREADY_USED_VALIDATION_CODE = "ALREADY_USED_VALIDATION_CODE"
	NOT_ALLOWED_RESOURCE         = "NOT_ALLOWED_RESOURCE"
	NOT_ENOUGH_WORDS_FOR_QUIZ    = "NOT_ENOUGH_WORDS_FOR_QUIZ"
)

type AlreadyUsedValidationCodeError struct{}

func (e AlreadyUsedValidationCodeError) Error() string {
	return "Already used validation code"
}

type EmailAlreadyInUseError struct{}

func (e EmailAlreadyInUseError) Error() string {
	return "Email already in use"
}

type EmptyFormFieldsError struct{}

func (e EmptyFormFieldsError) Error() string {
	return "Empty form fields"
}

type ExpiredValidationCodeError struct{}

func (e ExpiredValidationCodeError) Error() string {
	return "Expired validation code"
}

type IncorrectValidationCodeError struct{}

func (e IncorrectValidationCodeError) Error() string {
	return "Incorrect validation code"
}

type InvalidCredentialsError struct{}

func (e InvalidCredentialsError) Error() string {
	return "Invalid credentials"
}

type NotAllowedResourceError struct{}

func (e NotAllowedResourceError) Error() string {
	return "Not allowed resource"
}

type NotValidatedAccountError struct{}

func (e NotValidatedAccountError) Error() string {
	return "Not validated account"
}

type SendEmailError struct{}

func (e SendEmailError) Error() string {
	return "Send email error"
}

type UserAlreadyValidatedError struct{}

func (e UserAlreadyValidatedError) Error() string {
	return "User already validated"
}

type UsernameAlreadyInUseError struct{}

func (e UsernameAlreadyInUseError) Error() string {
	return "Username already in use"
}

type NotFoundError struct {
	Resource string
}

func (e NotFoundError) Error() string {
	return e.Resource + " not found: "
}

type NotEnoughWordsForQuizError struct{}

func (e NotEnoughWordsForQuizError) Error() string {
	return "Not enough words to create a quiz"
}
