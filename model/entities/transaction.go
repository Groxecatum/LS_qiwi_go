package entities

import "database/sql"

type Transaction struct {
	Id int64
}

func GetTransactionById(tx *sql.Tx, id int64) (*Transaction, error) {
	return nil, nil
}

func CreateNewTransaction(tx *sql.Tx) (Transaction, error) {
	return Transaction{}, nil
}
