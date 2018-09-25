package entities

import "database/sql"

type Account struct {
	Id     int `xml:"id"`
	TypeId int
}

const (
	DEFAULT_ACC_TYPE = 1
)

func GetAccountForWithdrawByPriority(tx *sql.Tx, cardId, merchantId int) (Account, error) {
	return Account{}, nil
}
