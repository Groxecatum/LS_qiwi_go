package errors

type CustomError struct {
	Text string
}

const customText = "Неизвестная ошибка"

func (e CustomError) Error() string {
	return customText + "[" + e.Text + "]"
}
