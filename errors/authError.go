package errors

type AuthError struct {
}

const authErrorText = "Ошибка аутентификации"

var AuthErr = AuthError{}

func (e AuthError) Error() string {
	return authErrorText
}
