package errors

type WrongOpError struct {
	Text string
}

const WronOpText = "Ошибка создания операции"

func (e WrongOpError) Error() string {
	return WronOpText + "[" + e.Text + "]"
}
