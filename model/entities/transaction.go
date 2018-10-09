package entities

import (
	"git.kopilka.kz/BACKEND/golang_commons/db"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type Transaction struct {
	Id        int64     `db:"biid"`
	DtCreated time.Time `db:"dtcreated"`
}

func GetTransactionById(tx *sqlx.Tx, id int64) (Transaction, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		trn := Transaction{}
		err := tx.Get(&trn, `select biid, dtcreated from ls.ttransactions where biid = $1`, id)
		if err != nil {
			log.Println(err)
			return trn, err
		}

		return trn, err
	}, tx)
	return res.(Transaction), err
}

func CreateNewTransaction(tx *sqlx.Tx) (Transaction, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		trn := Transaction{}
		rows, err := tx.Query(`INSERT INTO ls.ttransactions DEFAULT VALUES returning biid, dtcreated;`)
		if err != nil {
			log.Println(err)
			return trn, err
		}

		defer rows.Close()
		if rows.Next() {
			err := rows.Scan(&trn.Id, &trn.DtCreated)
			if err != nil {
				log.Println(err)
				return trn, err
			}
		}

		return trn, err
	}, tx)

	return res.(Transaction), err
}
