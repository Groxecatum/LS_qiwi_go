package errors

type AuthError struct {
	Text string
}

const authErrorText = "Ошибка аутентификации"

func (e AuthError) Error() string {
	return authErrorText + "[" + e.Text + "]"
}
