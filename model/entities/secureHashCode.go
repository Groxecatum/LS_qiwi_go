package entities

import (
	"git.kopilka.kz/BACKEND/golang_commons/db"
	"github.com/jmoiron/sqlx"
	"log"
)

type SecureHashCode struct {
	Id     int64
	CardId int
}

func GetSecureHashCodeById(tx *sqlx.Tx, id int64) (SecureHashCode, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		shc := SecureHashCode{}
		err := tx.Get(&shc, `select * from ls.tsecurehashcode where biid = $1`, id)
		if err != nil {
			log.Println(err)
			return shc, err
		}

		return shc, err
	}, tx)
	return res.(SecureHashCode), err
}

func GetSecureHashCodeByShortCode(tx *sqlx.Tx, shortCode string) (SecureHashCode, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		shc := SecureHashCode{}
		err := tx.Get(&shc, `select * from ls.tsecurehashcode where sshortcode = $1`, shortCode)
		if err != nil {
			log.Println(err)
			return shc, err
		}

		return shc, err
	}, tx)
	return res.(SecureHashCode), err
}

func GetSecureHashCodeByHash(tx *sqlx.Tx, secCode string) (SecureHashCode, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		shc := SecureHashCode{}
		err := tx.Get(&shc, `select * from ls.tsecurehashcode where hashcode = $1`, secCode)
		if err != nil {
			log.Println(err)
			return shc, err
		}

		return shc, err
	}, tx)
	return res.(SecureHashCode), err
}
