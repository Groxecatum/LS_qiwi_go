package errors

type WrongFormatError struct {
}

const wrongFormatText = "Сущность не найдена"

var WrongFormatErr = WrongFormatError{}

func (e WrongFormatError) Error() string {
	return notFoundText
}
