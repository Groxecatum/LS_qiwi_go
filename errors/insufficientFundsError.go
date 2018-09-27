package errors

type InsufficientFundsError struct {
}

const insufficientFundsText = "Недостаточно средств"

var InsufficientFundsErr = InsufficientFundsError{}

func (e InsufficientFundsError) Error() string {
	return insufficientFundsText
}
