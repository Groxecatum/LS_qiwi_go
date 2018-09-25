package errors

type NotAllowedError struct {
}

const notAllowedText = "Клиент не может участвовать в данный момент"

var NotAllowedErr = NotAllowedError{}

func (e NotAllowedError) Error() string {
	return notAllowedText
}
