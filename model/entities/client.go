package entities

import "database/sql"

type Client struct {
	Id           int
	IsRegistered bool
	FirstName    string
}

func GetClientById(tx *sql.Tx, id int) (*Client, error) {
	return nil, nil
}

func GetClientByCellPhone(tx *sql.Tx, phone string) (*Client, error) {
	return nil, nil
}

func CreateEmptyClient(tx *sql.Tx, phone string) (*Client, error) {
	return nil, nil
}
