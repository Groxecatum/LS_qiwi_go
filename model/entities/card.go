package entities

import "database/sql"

type Card struct {
	Id       int
	ClientId int
}

func GetCardById(tx *sql.Tx, id int) (*Card, error) {
	return nil, nil
}

func (crd *Card) GetClient(tx *sql.Tx) (*Client, error) {
	return GetClientById(tx, crd.Id)
}

func ExtractCardNum(fullNum string) (string, error) {
	return "", nil
}

func GenerateCardOnline(tx *sql.Tx, virtual bool, clientId int) error {
	return nil
}

func GetCardByNum(tx *sql.Tx, num string, blockForUpdate bool) (*Card, error) {
	return nil, nil
}
