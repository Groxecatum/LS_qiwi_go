package entities

import (
	"git.kopilka.kz/BACKEND/golang_commons/db"
	"git.kopilka.kz/BACKEND/golang_commons/errors"
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
	Id               int64     `db:"biid"`
	TransactionId    int64     `db:"bitrnid"`
	DtCreated        time.Time `db:"dtcreated"`
	OriginalCardId   int       `db:"ioriginalcardid"`
	InitiatorActorId int       `db:"isalespointid"`
}

func GetTransactionRequestById(tx *sqlx.Tx, id int64) (TransactionRequest, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		trnReq := TransactionRequest{}
		err := tx.Get(&trnReq, `select biid, dtcreated, dtcreated, ioriginalcardid from ls.ttrnrequests where biid = $1`, id)
		if err != nil {
			log.Println(err)
			return trnReq, err
		}

		return trnReq, err
	}, tx)
	return res.(TransactionRequest), err
}

func findByRefInTimePeriodEx(tx *sqlx.Tx, reference string, dateFrom, dateTo *time.Time, mt MerchantTerminal) (TransactionRequest, error) {
	ref := reference

	// TODO: Выпилить нахой - мы формируем фуллреф со временем
	if dateFrom != nil && (dateTo == nil) {
		ref = dateFrom.Format("060102") + reference
	}

	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		trnReqs := []TransactionRequest{}
		var err error

		if dateFrom != nil && dateTo != nil {
			err = tx.Select(&trnReqs, `SELECT biid, dtcreated, bitrnid, ioriginalcardid, isalespointid FROM ls.ttrnrequests
				WHERE iterminalid = $1 AND sreference = $2 AND dtcreatedext >= $3 AND dtcreatedext <= $4
				ORDER BY biid DESC;`, mt.Id, ref, dateFrom, dateTo)
		} else {
			err = tx.Select(&trnReqs, `SELECT biid, dtcreated, bitrnid, ioriginalcardid FROM ls.ttrnrequests
				WHERE iterminalid = $1 AND sfullreference = $2
				ORDER BY biid DESC;`, mt.Id, ref)
		}

		if err != nil {
			log.Println(err)
			return TransactionRequest{}, err
		}

		if len(trnReqs) > 1 {
			return TransactionRequest{}, errors.MultipleChoicesError{CustomError: errors.CustomError{"Транзакция с референсом " + ref + ": найдено больше одной"}}
		}

		if len(trnReqs) == 0 {
			return TransactionRequest{}, errors.NotFoundError{CustomError: errors.CustomError{"Транзакции с референсом " + ref + ": не найдено"}}
		}

		return trnReqs[0], err
	}, tx)
	return res.(TransactionRequest), err
}

func GetByRefUnified(tx *sqlx.Tx, ref string, mt MerchantTerminal, refDate *time.Time, blockDays *int) (TransactionRequest, error) {

	dateTo := refDate
	if blockDays == nil {
		newBlockDays := 0
		blockDays = &newBlockDays
	}

	if dateTo == nil {
		newDateTo := time.Now().AddDate(0, 0, 1)
		dateTo = &newDateTo
	}

	var dateFrom time.Time
	if refDate == nil {
		dateFrom = time.Now().AddDate(0, 0, -1)
	} else {
		if *blockDays != 0 {
			dateFrom = time.Now().AddDate(0, 0, -1**blockDays)
		} else {
			dateFrom = time.Now().AddDate(0, 0, -14)
		}
	}

	return findByRefInTimePeriodEx(tx, ref, &dateFrom, dateTo, mt)
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

func (request *TransactionRequest) Commit(tx *sqlx.Tx, mt MerchantTerminal, m Merchant, transactionId int64, commit int) error {

	start := time.Now()
	newCommitState := COMMITSTATE_COMMITTED
	if commit != COMMITSTATE_COMMITTED {
		newCommitState = COMMITSTATE_CANCELLED
	}

	err := request.updateOperationsState(tx, newCommitState)
	if err != nil {
		return err
	}

	log.Println("c.ops", time.Since(start))

	err = request.updateCommitState(tx, newCommitState)
	if err != nil {
		return err
	}

	log.Println("c.state", time.Since(start))

	// До разбирательства
	err = UpdateTrnItemsState(tx, transactionId)
	if err != nil {
		return err
	}

	log.Println("c.items_state", time.Since(start))

	card, err := GetCardById(tx, request.OriginalCardId, false)
	if err != nil {
		return err
	}

	client, err := GetClientById(tx, card.ClientId)
	if err != nil {
		return err
	}

	log.Println("c.card_client", time.Since(start))

	if newCommitState == COMMITSTATE_COMMITTED && (client.IsRegistered || m.CreateVirtualUser) {
		// TODO: Отправка в микросервис кампаний
		// Начисление акционных бонусов
		//CampaignActionProcessor.processOnPaymentCampaignsForSingleClient(conn, initialRequest.getId(), client.getId());

		//для виртуальной карты без точки продаж,
		//устанавливаем точку продаж этой транзакции, если она первая,
		if (card.IsVirtual) && (card.DistrSalesPointId == 0) {
			card.SetDistrSalesPointId(tx, request.InitiatorActorId, time.Now())
		}
	}

	log.Println("c.end", time.Since(start))

	return nil

	// TODO: Cancel
	//if (newCommitState == TrnOperation.COMMITSTATE_CANCELLED) {
	//	CardTrnOperationList initialRequestOperations = CardTrnOperationList.getOperationsByTrnReqId(conn, initialRequest.getId());
	//	for (TrnOperation o : initialRequestOperations)
	//	if (o.getCommitState() == TrnOperation.COMMITSTATE_COMMITTED || o.getCommitState() == TrnOperation.COMMITSTATE_INSTANT) {
	//		short oppositeOperTypeId = TrnOperation.getOppositeOperType((short) o.getTypeId());
	//		TrnOperation.regClientOperation(conn, ct, oppositeOperTypeId, Math.abs(o.getAmount()), false, card, Account.getById(conn,o.getCardAccountId()), // изменил на получение cardAccount с опрерации, т.к. getAcc возвращает только дефолтный кард акк.
	//		o.getMerchantTerminal(conn), null, initialRequest.getId(), -o.getWithdrawFeeAmount(), -o.getChargeFeeAmount(), o.getTypeEx(),
	//		TrnOperation.COMMITSTATE_INSTANT, o.getReferredOperationId(), "", true);
	//}
	//}
}

func (request *TransactionRequest) updateOperationsState(tx *sqlx.Tx, state int) error {
	_, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		_, err := tx.Exec("UPDATE ls.ttrnoperations set sicommitstate = $1 where bitrnrequestid = $2 and sicommitstate = $3;",
			// TODO: CANCEL
			//(commitState == COMMITSTATE_CANCELLED ? ", bcancelled = true \n" : "") +
			state, request.Id, COMMITSTATE_WAITFORCOMMIT)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return nil, err
	}, tx)
	return err
}

func (request *TransactionRequest) updateCommitState(tx *sqlx.Tx, state int) error {
	_, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		_, err := tx.Exec("UPDATE ls.ttrnrequests SET sicommitstate = $1 WHERE biid = $2;",
			state, request.Id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return nil, err
	}, tx)
	return err
}
