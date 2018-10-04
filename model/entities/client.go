package entities

import (
	"git.kopilka.kz/BACKEND/golang_commons/db"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type Client struct {
	Id           int       `db:"iid"`
	IsRegistered bool      `db:"bisregistered"`
	FirstName    string    `db:"sfirstname"`
	CellPhone    string    `db:"scellphone"`
	DtRegistered time.Time `db:"dtregistered"`
}

func GetClientById(tx *sqlx.Tx, id int) (Client, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		cli := Client{}
		err := tx.Get(&cli, `select iid, bisregistered, sfirstname, scellphone, dtregistered from ls.tclients where iid = $1`, id)
		if err != nil {
			log.Println(err)
			return cli, err
		}

		return cli, err
	}, tx)
	return res.(Client), err
}

func GetClientByCellPhone(tx *sqlx.Tx, phone string) (Client, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		cli := Client{}
		err := tx.Get(&cli, `select * from ls.tclients where scellphone = $1`, phone)
		if err != nil {
			log.Println(err)
			return cli, err
		}

		return cli, err
	}, tx)
	return res.(Client), err
}

func CreateEmptyClient(tx *sqlx.Tx, phone string) (Client, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		cli := Client{}
		pswSecurity, err := CreateNewSecurityEntry(tx, nil, 0)
		if err != nil {
			log.Println(err)
			return cli, err
		}
		pinSecurity, err := CreateNewSecurityEntry(tx, nil, 0)
		if err != nil {
			log.Println(err)
			return cli, err
		}
		rows, err := tx.Query(`INSERT INTO ls.tclients (scellphone, semail, sfirstName, slastName, ipinsecurityid, ipswsecurityid)
			VALUES ($1, NULL, NULL, NULL, $2, $3) RETURNING iid, dtregistered;`, phone, pinSecurity, pswSecurity)
		if err != nil {
			log.Println(err)
			return cli, err
		}

		defer rows.Close()

		if rows.Next() {
			err := rows.Scan(&cli.Id, cli.DtRegistered)
			if err != nil {
				log.Println(err)
				return cli, err
			}
		}

		return cli, err
	}, tx)
	return res.(Client), err
}

func GetPinById(tx *sqlx.Tx, id int) (int, error) {
	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		var pin int
		err := tx.Get(&pin, `select ipinsecurityid from ls.tclients where iid = $1`, id)
		if err != nil {
			log.Println(err)
			return pin, err
		}

		return pin, err
	}, tx)
	return res.(int), err
}
