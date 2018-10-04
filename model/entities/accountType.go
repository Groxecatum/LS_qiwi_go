package entities

import (
	"git.kopilka.kz/BACKEND/golang_commons/db"
	"github.com/jmoiron/sqlx"
	"log"
)

type AccountType struct {
	Id       int `db:"iid" xml:"id"`
	BurnDays int `db:"iburnafterdays" xml:"burndays"`
}

func GetAccountTypeById(tx *sqlx.Tx, id int) (AccountType, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		crd := AccountType{}

		err := tx.Get(&crd, `select iid, iburnafterdays from ls.tcardaccounttypes where iid = $1`, id)
		if err != nil {
			log.Println(err)
			return crd, err
		}

		return crd, err
	}, tx)
	return res.(AccountType), err
}
