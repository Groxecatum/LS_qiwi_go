package entities

import (
	"git.kopilka.kz/BACKEND/golang_commons/db"
	"github.com/jmoiron/sqlx"
	"log"
)

type MerchantTerminal struct {
	Id                  int     `db:"iid"`
	MerchantId          int     `db:"imerchantid"`
	NeedPostponedCommit bool    `db:"bneedpostponecommit"`
	IsTest              bool    `db:"btest"`
	AllowedMinimum      float64 `db:"nallowedminimum"`
}

func GetMerchantTerminal(tx *sqlx.Tx, actorId, terminalNum int, lock bool) (MerchantTerminal, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		mt := MerchantTerminal{}
		forUpdStr := ""
		if lock {
			forUpdStr = " FOR UPDATE"
		}
		err := tx.Get(&mt, `select iid, imerchantid, bneedpostponecommit, btest, nallowedminimum from ls.tmerchantterminals 
								where isalespointid = $1 and iterminalnum = $2 `+forUpdStr, actorId, terminalNum)
		if err != nil {
			log.Println(err)
			return mt, err
		}

		return mt, err
	}, tx)
	return res.(MerchantTerminal), err
}

func GetMerchantTerminalById(tx *sqlx.Tx, terminalId int, lock bool) (MerchantTerminal, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		mt := MerchantTerminal{}
		forUpdStr := ""
		if lock {
			forUpdStr = " FOR UPDATE"
		}
		err := tx.Get(&mt, `select iid, imerchantid, bneedpostponecommit, btest from ls.tmerchantterminals where iid = $1 `+forUpdStr,
			terminalId)
		if err != nil {
			log.Println(err)
			return mt, err
		}

		return mt, err
	}, tx)
	return res.(MerchantTerminal), err
}
