package errors

type DBError struct {
}

const DBText = "Ошибка БД"

var DBErr = DBError{}

func (e DBError) Error() string {
	return DBText
}
