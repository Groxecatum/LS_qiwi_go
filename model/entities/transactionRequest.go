package entities

import (
	"git.kopilka.kz/BACKEND/golang_commons/db"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

const (
	COMMITSTATE_INSTANT       = 2
	COMMITSTATE_WAITFORCOMMIT = 0
	COMMITSTATE_COMMITTED     = 1
	COMMITSTATE_CANCELLED     = -1

	TRNREQTYPE_BONUS_TRN               = 1
	TRNREQTYPE_DECREASE_TRN            = 2
	TRNREQTYPE_BONUS_TRN_WITH_WITHDRAW = 3
	TRNREQTYPE_MERGECARDSOFONECLIENT   = 4
	TRNREQTYPE_MERCHANT_TO_MERCHANT    = 5
)

type TransactionRequest struct {
	Id            int64
	TransactionId int64
	DtCreated     time.Time
}

func GetTransactionRequestById(tx *sqlx.Tx, id int64) (TransactionRequest, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		trnReq := TransactionRequest{}
		err := tx.Get(&trnReq, `select * from ls.ttrnrequests where biid = $1`, id)
		if err != nil {
			log.Println(err)
			return trnReq, err
		}

		return trnReq, err
	}, tx)
	return res.(TransactionRequest), err
}

func findByRefInTimePeriodEx(tx *sqlx.Tx, ref string, dateFrom, dateTo time.Time, mt MerchantTerminal) (TransactionRequest, error) {
	return TransactionRequest{}, nil
}

func GetByRefUnified(tx *sqlx.Tx, ref string, mt MerchantTerminal, refDate *time.Time, blockDays *int) (TransactionRequest, error) {

	dateTo := refDate
	if blockDays == nil {
		*blockDays = 0
	}

	if dateTo == nil {
		*dateTo = time.Now().AddDate(0, 0, 1)
	}

	var dateFrom time.Time
	if refDate == nil {
		dateFrom = time.Now().AddDate(0, 0, -1)
	} else {
		if *blockDays == 0 {
			dateFrom = time.Now().AddDate(0, 0, (-1 * *blockDays))
		} else {
			dateFrom = time.Now().AddDate(0, 0, -14)
		}
	}

	return findByRefInTimePeriodEx(tx, ref, dateFrom, *dateTo, mt)
}

func CreateNewTransactionRequest(tx *sqlx.Tx, trnType int, trnId int64, merchantId int, merchantTerminalId int,
	actorId int, descr string, ref string, fullRef string, date time.Time, zRepId string, commitState int,
	checkId string, batchPeriodId string, cardId *int, acceptorMerchantId *int, acceptorActorId *int) (TransactionRequest, error) {

	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		trnReq := TransactionRequest{TransactionId: trnId}
		rows, err := tx.Query(`INSERT INTO ls.ttrnrequests (sitypeid, bitrnid,
				imerchantid, iterminalid, isalespointid, sdescr, sreference,
				sfullreference, dtcreatedext, szrepid, sicommitstate,
				scheckid, sbatchperiodid, ioriginalcardid, iacceptormerchantid, iacceptoractorid)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) returning biid, dtcreated;`,
			trnType, trnId, merchantId, merchantTerminalId, actorId, descr, ref, fullRef, date,
			zRepId, commitState, checkId, batchPeriodId, cardId, acceptorMerchantId, acceptorActorId)
		if err != nil {
			log.Println(err)
			return trnReq, err
		}

		defer rows.Close()
		if rows.Next() {
			err := rows.Scan(&trnReq.Id, &trnReq.DtCreated)
			if err != nil {
				log.Println(err)
				return trnReq, err
			}
		}

		return trnReq, err
	}, tx)
	return res.(TransactionRequest), err
}
