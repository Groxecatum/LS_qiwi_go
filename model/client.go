package model

import "database/sql"

type Client struct {
}

func GetClientById(tx *sql.Tx, id int) (*Client, error) {

}

func GetClientByCellPhone(tx *sql.Tx, phone string) (*Client, error) {

}
