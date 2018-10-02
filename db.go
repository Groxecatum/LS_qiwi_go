package golang_commons

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

var db *sqlx.DB

type DBExecutorFunc func(tx *sql.Tx) (interface{}, error)
type DBXExecutorFunc func(tx *sqlx.Tx) (interface{}, error)

func EstablishDB(dbhost, dbport, dbuser, dbpass, dbname string) {
	var err error

	db, err = sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpass, dbname))
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(40)
}

// Логика такая - кто создал транзакцию - тот ее и коммитит\роллбечит
func initTrn(extTx *sql.Tx) (tx *sql.Tx, err error, externalTrn bool) {
	externalTrn = extTx != nil
	if externalTrn {
		return extTx, nil, externalTrn
	} else {
		tx, err = db.Begin()
		return
	}
}

func Do(f DBExecutorFunc, extTx *sql.Tx) (intf interface{}, err error) {
	// Инициализируем новую транзакцию, если мы не в существующей
	tx, err, extTrn := initTrn(extTx)
	if err != nil {
		log.Fatal(err)
	}

	//defer tx.Rollback()

	intf, err = f(tx)

	// Откат при ошибке
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return intf, err
	} else {
		// Коммит, если все норм и мы начинали транзакцию
		if !extTrn {
			err = tx.Commit()
		}
		return
	}
}

func initTrnX(extTx *sqlx.Tx) (tx *sqlx.Tx, err error, externalTrn bool) {
	externalTrn = extTx != nil
	if externalTrn {
		return extTx, nil, externalTrn
	} else {
		tx, err = db.Beginx()
		return
	}
}

func DoX(f DBXExecutorFunc, extTx *sqlx.Tx) (intf interface{}, err error) {
	// Инициализируем новую транзакцию, если мы не в существующей
	tx, err, extTrn := initTrnX(extTx)
	if err != nil {
		log.Fatal(err)
	}

	//defer tx.Rollback()

	intf, err = f(tx)

	// Откат при ошибке
	if err != nil {
		tx.Rollback()
		log.Println(err)
		return intf, err
	} else {
		// Коммит, если все норм и мы начинали транзакцию
		if !extTrn {
			err = tx.Commit()
		}
		return
	}
}
