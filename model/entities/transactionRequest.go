package entities

import (
	"database/sql"
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
}

func GetTransactionRequestById(tx *sql.Tx, id int64) (*TransactionRequest, error) {
	return nil, nil
}

func GetByRefUnified() (*TransactionRequest, error) {
	return nil, nil
}

func CreateNewTransactionRequest(tx *sql.Tx, trnType int, trnId int64, merchantId int, merchantTerminalId int,
	actorId int, descr string, ref string, fullRef string, date time.Time, requestId *int64, zRepId string, commitState int,
	checkId string, batchPeriodId string, cardId *int, acceptorMerchantId *int, acceptorActorId *int) (TransactionRequest, error) {
	return TransactionRequest{}, nil
}
