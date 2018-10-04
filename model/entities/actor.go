package entities

import (
	"git.kopilka.kz/BACKEND/golang_commons/db"
	"github.com/jmoiron/sqlx"
	"log"
)

type Actor struct {
	Id         int    `db:"iid"  xml:"id"`
	MerchantId int    `db:"imerchantid"  xml:"merchantId"`
	Title      string `db:"stitle"  xml:"title"`
}

func GetActorById(tx *sqlx.Tx, id int) (Actor, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		actor := Actor{}
		err := tx.Get(&actor, `select * from ls.tactors where iid = $1`, id)
		if err != nil {
			log.Println(err)
			return actor, err
		}

		return actor, err
	}, tx)
	return res.(Actor), err
}

func GetActorByLogin(tx *sqlx.Tx, login string) (Actor, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		actor := Actor{}
		//err := tx.Get(&actor, `select * from ls.tactors where slogin = $1`, login)
		err := tx.Get(&actor, `select iid, imerchantid, stitle from ls.tactors where slogin = $1`, login)
		if err != nil {
			log.Println(err)
			return actor, err
		}

		return actor, err
	}, tx)
	return res.(Actor), err
}
