package model

import "database/sql"

type SecureHashCode struct {
	Id     int64
	CardId int
}

func GetSecureHashCodeById(tx sql.Tx, id int) {

}

func GetSecureHashCodeByShortCode(tx *sql.Tx, shortCode string) (SecureHashCode, error) {

}

func GetSecureHashCodeByHash(tx *sql.Tx, secCode string) (SecureHashCode, error) {

}
