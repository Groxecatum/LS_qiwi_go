package model

import "database/sql"

type Card struct {
}

func GetCardById(tx *sql.Tx, id int) (*Card, error) {
	return nil, nil
}
