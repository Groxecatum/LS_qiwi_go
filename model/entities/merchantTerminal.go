package entities

import (
	"database/sql"
	"git.kopilka.kz/BACKEND/golang_commons"
	"git.kopilka.kz/BACKEND/golang_commons/errors"
	"log"
)

type MerchantTerminal struct {
	Id                  int
	MerchantId          int
	NeedPostponedCommit bool
}

func GetMerchantTerminalById(tx *sql.Tx, id int) (*MerchantTerminal, error) {
	// TODO
	return nil, nil
}

func getAllTerminalFields() string {
	return " iid, imerchantid, bpostponedcommit "
}

func terminalFromRows(rows *sql.Rows) (MerchantTerminal, error) {
	mt := MerchantTerminal{}
	err := rows.Scan(&mt.Id, &mt.MerchantId, &mt.NeedPostponedCommit)
	return mt, err
}

func GetMerchantTerminal(tx *sql.Tx, actorId, terminalNum int, lock bool) (MerchantTerminal, error) {
	res, err := golang_commons.Do(func(tx *sql.Tx) (interface{}, error) {
		mt := MerchantTerminal{}
		forUpdStr := ""
		if lock {
			forUpdStr = " FOR UPDATE"
		}
		rows, err := tx.Query("select "+getAllTerminalFields()+" from ls.tmerchantterminals \n"+
			" where isalespointid = ? and iterminalnum = ? "+forUpdStr, actorId, terminalNum)

		if err != nil {
			log.Println(err)
			return mt, err
		}
		defer rows.Close()

		if rows.Next() {
			mt, err = terminalFromRows(rows)
			if err != nil {
				log.Println(err)
				return mt, err
			}
		} else {
			return mt, errors.AuthErr
		}

		return mt, err
	}, tx)
	return res.(MerchantTerminal), err
}
