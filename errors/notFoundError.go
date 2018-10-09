package errors

type NotFoundError struct {
	CustomError
}

const notFoundText = "Сущность не найдена"

func (e NotFoundError) Error() string {
	return notFoundText + "[" + e.Text + "]"
}
