package entities

import (
	"git.kopilka.kz/BACKEND/golang_commons"
	"github.com/jmoiron/sqlx"
	"log"
)

type Actor struct {
	Id         int    `db:"iid"  xml:"id"`
	MerchantId int    `db:"imerchantid"  xml:"merchantId"`
	Title      string `db:"stitle"  xml:"title"`
	Email      string `db:"semail"  xml:"email"`
}

func GetActorById(tx *sqlx.Tx, id int) (Actor, error) {
	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
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
