package errors

type MultipleChoicesError struct {
	CustomError
}

const MultipleChoicesText = "Запрос вернул более одной сущности"

func (e MultipleChoicesError) Error() string {
	return MultipleChoicesText + "[" + e.Text + "]"
}
