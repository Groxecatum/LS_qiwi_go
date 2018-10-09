package errors

type WrongFormatError struct {
	CustomError
}

const wrongFormatText = "Неверный формат запроса"

func (e WrongFormatError) Error() string {
	return wrongFormatText + "[" + e.Text + "]"
}
