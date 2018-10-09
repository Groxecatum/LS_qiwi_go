package errors

type NotAllowedError struct {
	CustomError
}

const notAllowedText = "Клиент не может участвовать в данный момент"

func (e NotAllowedError) Error() string {
	return notAllowedText + "[" + e.Text + "]"
}
