package entities

import "database/sql"

type Merchant struct {
	Id                int
	CreateVirtualUser bool
}

func GetMerchantById(tx *sql.Tx, id int) (*Merchant, error) {
	return nil, nil
}

func GetMerchantDataForTransaction(tx *sql.Tx, merchantId int) (Merchant, error) {
	return Merchant{}, nil
}
