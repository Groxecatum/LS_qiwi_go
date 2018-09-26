package errors

type NotFoundError struct {
}

const notFoundText = "Сущность не найдена"

var NotFoundErr = NotFoundError{}

func (e NotFoundError) Error() string {
	return notFoundText
}
