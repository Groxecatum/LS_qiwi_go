package entities

import (
	"github.com/jmoiron/sqlx"
)

type SecureHashCode struct {
	Id     int64
	CardId int
}

func GetSecureHashCodeById(tx sqlx.Tx, id int) {

}

func GetSecureHashCodeByShortCode(tx *sqlx.Tx, shortCode string) (SecureHashCode, error) {
	return SecureHashCode{}, nil
}

func GetSecureHashCodeByHash(tx *sqlx.Tx, secCode string) (SecureHashCode, error) {
	return SecureHashCode{}, nil
}
