package entities

import (
	"git.kopilka.kz/BACKEND/golang_commons"
	"github.com/jmoiron/sqlx"
	"log"
)

type MerchantTerminal struct {
	Id                  int  `db:"iid"`
	MerchantId          int  `db:"imerchantid"`
	NeedPostponedCommit bool `db:"bneedpostponecommit"`
}

func GetMerchantTerminal(tx *sqlx.Tx, actorId, terminalNum int, lock bool) (MerchantTerminal, error) {
	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		mt := MerchantTerminal{}
		forUpdStr := ""
		if lock {
			forUpdStr = " FOR UPDATE"
		}
		err := tx.Get(&mt, `select iid, imerchantid, bneedpostponecommit from ls.tmerchantterminals where isalespointid = $1 and iterminalnum = $2 `+forUpdStr,
			actorId, terminalNum)
		if err != nil {
			log.Println(err)
			return mt, err
		}

		return mt, err
	}, tx)
	return res.(MerchantTerminal), err
}
