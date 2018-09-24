package errors

type BalanceError struct {
}

const balanceErrorText = "Ошибка получения баланса"

var BalanceErr = BalanceError{}

func (e BalanceError) Error() string {
	return balanceErrorText
}
