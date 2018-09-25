package entities

import "database/sql"

type MerchantTerminal struct {
	Id         int
	MerchantId int
}

func GetMerchantTerminalById(tx *sql.Tx, id int) (*MerchantTerminal, error) {
	return nil, nil
}
