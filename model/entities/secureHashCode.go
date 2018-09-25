package entities

import "database/sql"

type SecureHashCode struct {
	Id     int64
	CardId int
}

func GetSecureHashCodeById(tx sql.Tx, id int) {

}

func GetSecureHashCodeByShortCode(tx *sql.Tx, shortCode string) (SecureHashCode, error) {
	return SecureHashCode{}, nil
}

func GetSecureHashCodeByHash(tx *sql.Tx, secCode string) (SecureHashCode, error) {
	return SecureHashCode{}, nil
}
