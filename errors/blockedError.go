package errors

type BlockedError struct {
}

const blockedText = "Оплата заблокирована"

var BlockedErr = BlockedError{}

func (e BlockedError) Error() string {
	return blockedText
}
