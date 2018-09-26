package entities

import (
	"database/sql"
	"git.kopilka.kz/BACKEND/golang_commons"
	"git.kopilka.kz/BACKEND/golang_commons/errors"
	"log"
)

type Actor struct {
	Id         int `xml:"id"`
	MerchantId int
	Title      string `xml:"title"`
	Email      string `xml:"email"`
}

func getAllActorFields() string {
	return " iid, imerchantid, stitle, email "
}

func actorFromRows(rows *sql.Rows) (Actor, error) {
	actor := Actor{}
	err := rows.Scan(&actor.Id, &actor.MerchantId, &actor.Title, &actor.Email)
	return actor, err
}

func GetActorById(tx *sql.Tx, id int) (Actor, error) {
	res, err := golang_commons.Do(func(tx *sql.Tx) (interface{}, error) {
		actor := Actor{}
		rows, err := tx.Query("select "+getAllActorFields()+" from ls.tactors where iid = $1", id)

		if err != nil {
			log.Println(err)
			return actor, err
		}
		defer rows.Close()

		if rows.Next() {
			actor, err = actorFromRows(rows)
			if err != nil {
				log.Println(err)
				return actor, err
			}
		} else {
			return actor, errors.AuthErr
		}

		return actor, err
	}, tx)
	return res.(Actor), err
}
