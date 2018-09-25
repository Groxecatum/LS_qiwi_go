package entities

import "database/sql"

type Actor struct {
	Id         int `xml:"id"`
	MerchantId int
	Title      string `xml:"title"`
	Email      string `xml:"email"`
}

func GetActorById(tx *sql.Tx, id int) (*Actor, error) {
	return nil, nil
}

func (actor *Actor) GetMerchantTerminal(tx *sql.Tx, id int, lock bool) (MerchantTerminal, error) {
	return MerchantTerminal{}, nil
}
