package errors

type WrongOpError struct {
}

const WronOpText = "Ошибка создания операции"

var WrongOpErr = WrongOpError{}

func (e WrongOpError) Error() string {
	return WronOpText
}
