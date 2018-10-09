package entities

import (
	"database/sql"
	"git.kopilka.kz/BACKEND/golang_commons/db"
	"git.kopilka.kz/BACKEND/golang_commons/errors"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

const (
	TRNOPERTYPE_TERMINAL_TO_CLIENT_BLOCK           = 1 // перечисление с терминала в блок. Используется когда коммерс начисляет бонусы, но они держатся в блоке 2 недели)
	TRNOPERTYPE_CLIENT_BLOCK_TO_ACTIVE             = 2 // перечисление из блока в актив. Используется, когда срок удержания бонусов в блоке заканчивается и надо их отдать клиенту на использование. Обычно выполняется по расписанию (спустя 2 недели).
	TRNOPERTYPE_CLIENT_BLOCK_TO_TERMINAL           = 3 // перечисление из блока в терминал. Используется, когда клиент возвращает товар, а бонусы, начисленные за этот товар, лежат у клиента в блоке т.е. бонусы забираются из блока.
	TRNOPERTYPE_CLIENT_ACTIVE_TO_TERMINAL          = 4 // перечисление с кардсчета на терминал. Используется при оплате бонусами. Либо когда надо вернуть начисленные бонусы, но в блоке они не удерживаются (0 дней блокировки).
	TRNOPERTYPE_TERMINAL_TO_CLIENT_ACTIVE          = 5 // перечисление с терминала на кардсчет. Используется, когда клиент платил бонусами и вернул товар. коммерсант возаращет оплаченные клиентом бонусы на карту.
	TRNOPERTYPE_CLIENT_ACTIVE_TO_CLIENT_BLOCK      = 6 // перечисление с терминала на кардсчет. Используется, когда клиент платил бонусами и вернул товар. коммерсант возаращет оплаченные клиентом бонусы на карту.
	TRNOPERTYPE_CANCEL_OPERATION                   = 7 // отмена операции
	TRNOPERTYPE_UNCANCEL_OPERATION                 = 8
	TRNOPERTYPE_FOREIGN_TERMINAL_TO_TERMINAL_BLOCK = 9  // перечисление с терминала внешнего партнера в блок терминала
	TRNOPERTYPE_TERMINAL_BLOCK_TO_TERMINAL         = 10 // перечисление с блока терминала в актив счета терминала
	TRNOPERTYPE_FOREIGN_TERMINAL_TO_TERMINAL       = 11 // перечисление с терминала в терминал напрямую, минуя блок
	TRNOPERTYPE_TERMINAL_BLOCK_TO_FOREIGN_TERMINAL = 12 // перечисление с блока терминала в терминал внешнего партнера
	TRNOPERTYPE_TERMINAL_TO_TERMINAL_BLOCK         = 13 // перечисление с актива терминала в блок терминала
	TRNOPERTYPE_TERMINAL_TO_FOREIGN_TERMINAL       = 14 // перечисление с актива терминала во внешний терминал минуя блок
	TRNOPERTYPE_TERMINAL_BLOCK_TO_CLIENT           = 15 // перечисление с блока терминала клиенту
	TRNOPERTYPE_CLIENT_TO_TERMINAL                 = 16 // перечисление с клиента на терминал
	TRNOPERTYPE_CLIENT_TO_TERMINAL_BLOCK           = 17 // перечисление с клиента на блок терминала

	TRNOPERTYPEEX_UNDEFINED            = 0
	TRNOPERTYPEEX_GAIN_BONUSES         = 1 // накопление бонусов
	TRNOPERTYPEEX_PAY_WITH_BONUSES     = 2 // трата бонусов
	TRNOPERTYPEEX_GAIN_BONUSES_REV     = 3 // возврат накопленных бонусов
	TRNOPERTYPEEX_PAY_WITH_BONUSES_REV = 4 // возврат оплчаенных бонусов
	TRNOPERTYPEEX_COMMIT               = 5
)

type TransactionOperation struct {
	Id                  int64     `db:"biid"`
	TransactionId       int64     `db:"bitrnid"`
	DtCreated           time.Time `db:"dtcreated"`
	TypeId              int       `db:"sitype"`
	ReferredOperationId int64     `db:"bireferredoperationid"`
	AccountId           int       `db:"icardaccountid"`
	AmountChange        float64   `db:"namountchange"`
	BlockedAmountChange float64   `db:"nblockedamountchange"`
}

func GetOperationById(tx *sqlx.Tx, id int64) (TransactionOperation, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		op := TransactionOperation{}
		err := tx.Get(&op, `select biid, bitrnid, dtcreated, sitype, bireferredoperationid, icardaccountid, namountchange * 0.01 as namountchange, 
										nblockedamountchange * 0.01 as nblockedamountchange
									  from ls.ttrnoperations where biid = $1`, id)
		if err != nil {
			log.Println(err)
			return op, err
		}

		return op, err
	}, tx)
	return res.(TransactionOperation), err
}

func setOperationCancelled(tx *sqlx.Tx, opId int64, newState bool) error {
	_, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		_, err := tx.Exec(`UPDATE ls.ttrnoperations SET bcancelled = $1 WHERE biid = $2 and siprocessed = 0`,
			newState, opId)
		return nil, err
	}, tx)

	return err
}

func ProcessOperation(tx *sqlx.Tx, operation TransactionOperation) error {
	var crcr AccountChange
	var err error
	if operation.TypeId == TRNOPERTYPE_CANCEL_OPERATION || operation.TypeId == TRNOPERTYPE_UNCANCEL_OPERATION {
		if operation.ReferredOperationId > 0 {
			newCancelState := operation.TypeId == TRNOPERTYPE_CANCEL_OPERATION
			err := setOperationCancelled(tx, operation.ReferredOperationId, newCancelState)
			if err != nil {
				return err
			}
		}
	} else {
		crcr, err = RegAccountChange(tx, operation.AccountId, operation.AmountChange, operation.BlockedAmountChange)
		if err != nil {
			return err
		}
	}

	_, err = db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		_, err := tx.Exec(`update ls.ttrnoperations a
			 set siprocessed = 1, nresultamount = $1, nresultblockedamount = $2,
			 dtprocessed = CURRENT_TIMESTAMP where biid = $3`,
			crcr.ResultAmount, crcr.ResultBlockedAmount, operation.Id)
		return nil, err
	}, tx)

	return err
}

func RegMerchantOperation(tx *sqlx.Tx, t Transaction, operationTypeId int,
	amount float64, processOnline bool, account Account,
	m Merchant, terminal MerchantTerminal, _scheduledTime *time.Time,
	trnRequestId int64, withdrawFeeAmount int64, chargeFeeAmount int64, typeEx int,
	commitState int, referredOperationId int64) error {

	if operationTypeId < TRNOPERTYPE_UNCANCEL_OPERATION {
		return errors.WrongOpError{"Попытка создать операцию партнера с типом операции клиента"}
	}

	var amountChange float64
	var blockedAmountChange float64

	if operationTypeId == TRNOPERTYPE_FOREIGN_TERMINAL_TO_TERMINAL_BLOCK {
		amountChange = 0
		blockedAmountChange = amount
	} else if operationTypeId == TRNOPERTYPE_TERMINAL_BLOCK_TO_TERMINAL {
		amountChange = amount
		blockedAmountChange = -amount
	} else if operationTypeId == TRNOPERTYPE_FOREIGN_TERMINAL_TO_TERMINAL {
		amountChange = amount
		blockedAmountChange = 0
	} else if operationTypeId == TRNOPERTYPE_TERMINAL_BLOCK_TO_FOREIGN_TERMINAL {
		amountChange = 0
		blockedAmountChange = -amount
		amount = -amount
	} else if operationTypeId == TRNOPERTYPE_TERMINAL_TO_TERMINAL_BLOCK {

		if (account.Bonuses < amount) && (!m.AllowNegativeBalance) {
			return errors.InsufficientFundsError{}
		}

		amountChange = -amount
		blockedAmountChange = amount
		amount = -amount
	} else if operationTypeId == TRNOPERTYPE_TERMINAL_TO_FOREIGN_TERMINAL {
		amountChange = -amount
		blockedAmountChange = 0
		amount = -amount
	} else if operationTypeId == TRNOPERTYPE_TERMINAL_BLOCK_TO_CLIENT {
		amountChange = 0
		blockedAmountChange = -amount
	} else if operationTypeId == TRNOPERTYPE_CLIENT_TO_TERMINAL {
		amountChange = amount
		blockedAmountChange = 0
	} else if operationTypeId == TRNOPERTYPE_CLIENT_TO_TERMINAL_BLOCK {
		amountChange = 0
		blockedAmountChange = amount
	} else if operationTypeId == TRNOPERTYPE_CANCEL_OPERATION || operationTypeId == TRNOPERTYPE_UNCANCEL_OPERATION {
		amountChange = 0
		blockedAmountChange = 0
	} else {
		return errors.WrongOpError{"Неверный тип операции"}
	}

	if commitState == 0 {
		processOnline = false
	}

	scheduledTime := time.Now()
	if _scheduledTime != nil {
		scheduledTime = *_scheduledTime
	}

	processStatus := 0
	if processOnline {
		processStatus = 1
	}

	_, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		rows, err := tx.Query(`INSERT INTO ls.ttrnoperations(bitrnid,sitype, namountchange, nblockedamountchange,
			siprocessed, icardaccountid, icardid, iterminalid, dtScheduledTime, bitrnrequestid, nchargefeeamount,
			nwithdrawfeeamount, sitypeex, sicommitstate, namount, bireferredoperationid, sdescr)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) returning biid;`,
			t.Id, operationTypeId, amountChange*100, blockedAmountChange*100, processStatus, account.Id, sql.NullInt64{}, terminal.Id,
			scheduledTime, trnRequestId, chargeFeeAmount*100, withdrawFeeAmount*100, typeEx, commitState, amount*100, referredOperationId, "")
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		operation := TransactionOperation{TransactionId: t.Id, DtCreated: time.Now(), TypeId: operationTypeId,
			AccountId: account.Id, AmountChange: amountChange, BlockedAmountChange: blockedAmountChange,
			ReferredOperationId: referredOperationId}
		if rows.Next() {
			err = rows.Scan(&operation.Id)
			if err != nil {
				return nil, err
			}

			rows.Close()

			if processOnline {
				// Вот тут мы прикроем принудительно rows, иначе он пошлет нас при попытке дернуть данные обработки
				err = ProcessOperation(tx, operation)
				if err != nil {
					return nil, err
				}
			}
		}

		return nil, err
	}, tx)
	return err
}

func RegClientOperation(tx *sqlx.Tx, t Transaction, operationTypeId int,
	amount float64, processOnline bool, card Card, account Account,
	m Merchant, terminal MerchantTerminal, _scheduledTime *time.Time,
	trnRequestId int64, withdrawFeeAmount int64, chargeFeeAmount int64, typeEx int,
	commitState int, referredOperationId int64, descr string,
	includeBlockedCards bool, fromDecrease bool) error {
	if operationTypeId > TRNOPERTYPE_UNCANCEL_OPERATION {
		return errors.WrongOpError{"Попытка создать операцию клиента с типом операции партнера"}
	}

	if card.IsTest != account.IsTest || card.IsTest != terminal.IsTest || account.IsTest != terminal.IsTest {
		return errors.WrongOpError{"Попытка использовать тестовый терминал с боевой картой(или наоборот, боевой терминал с тестовой)"}
	}

	if !includeBlockedCards {
		if card.IsBlocked {
			return errors.WrongOpError{"Карта заблокирована"}
		} else if account.IsBlocked {
			return errors.WrongOpError{"Аккаунт заблокирован"}
		}
	}

	var amountChange float64
	var blockedAmountChange float64
	if operationTypeId == TRNOPERTYPE_TERMINAL_TO_CLIENT_BLOCK {
		amountChange = 0
		blockedAmountChange = amount
	} else if operationTypeId == TRNOPERTYPE_CLIENT_BLOCK_TO_ACTIVE {
		amountChange = amount
		blockedAmountChange = -amount
	} else if operationTypeId == TRNOPERTYPE_CLIENT_BLOCK_TO_TERMINAL {
		amountChange = 0
		blockedAmountChange = -amount
		amount = -amount
	} else if operationTypeId == TRNOPERTYPE_CLIENT_ACTIVE_TO_TERMINAL {
		if ((account.Bonuses-terminal.AllowedMinimum) < float64(amount) && typeEx < 10) &&
			(!m.AllowNegativeDecrease || !fromDecrease) {
			return errors.WrongOpError{"Ограничение терминала"}
		}
		amountChange = -amount
		blockedAmountChange = 0
		amount = -amount
	} else if operationTypeId == TRNOPERTYPE_CLIENT_ACTIVE_TO_CLIENT_BLOCK {
		if ((account.Bonuses-terminal.AllowedMinimum) < amount && typeEx < 10) &&
			(!m.AllowNegativeDecrease || !fromDecrease) {
			return errors.InsufficientFundsError{}
		}
		amountChange = -amount
		blockedAmountChange = amount
		amount = -amount
	} else if operationTypeId == TRNOPERTYPE_TERMINAL_TO_CLIENT_ACTIVE {
		amountChange = amount
		blockedAmountChange = 0
	} else if operationTypeId == TRNOPERTYPE_CANCEL_OPERATION || operationTypeId == TRNOPERTYPE_UNCANCEL_OPERATION {
		amountChange = 0
		blockedAmountChange = 0
	} else {
		return errors.WrongOpError{"Неверный тип операции"}
	}

	if commitState == 0 {
		processOnline = false
	}

	processStatus := 0
	if processOnline {
		processStatus = 1
	}

	scheduledTime := time.Now()
	if _scheduledTime != nil {
		scheduledTime = *_scheduledTime
	}

	_, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		rows, err := tx.Query(`INSERT INTO ls.ttrnoperations(bitrnid,sitype, namountchange, nblockedamountchange,
			siprocessed, icardaccountid, icardid, iterminalid, dtScheduledTime, bitrnrequestid, nchargefeeamount,
			nwithdrawfeeamount, sitypeex, sicommitstate, namount, bireferredoperationid, sdescr)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) returning biid;`,
			t.Id, operationTypeId, amountChange*100, blockedAmountChange*100, processStatus, account.Id, card.Id, terminal.Id,
			scheduledTime, trnRequestId, chargeFeeAmount*100, withdrawFeeAmount*100, typeEx, commitState, amount*100, referredOperationId,
			descr)
		if err != nil {
			return nil, err
		}

		defer rows.Close()
		operation := TransactionOperation{TransactionId: t.Id, DtCreated: time.Now(), TypeId: operationTypeId,
			AccountId: account.Id, AmountChange: amountChange, BlockedAmountChange: blockedAmountChange,
			ReferredOperationId: referredOperationId}
		if rows.Next() {
			err = rows.Scan(&operation.Id)
			if err != nil {
				return nil, err
			}

			rows.Close()

			if processOnline {
				err = ProcessOperation(tx, operation)
				if err != nil {
					return nil, err
				}
			}
		}

		return nil, err
	}, tx)
	return err
}

func GetNextOperation(tx *sqlx.Tx) (TransactionOperation, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		ops := []TransactionOperation{}
		err := tx.Select(&ops, `select biid, bitrnid, dtcreated, sitype, bireferredoperationid, icardaccountid, namountchange * 0.01 as namountchange, 
										nblockedamountchange * 0.01 as nblockedamountchange from ls.ttrnoperations 
					where siprocessed = 0 and dtscheduledtime <= CURRENT_TIMESTAMP and bcancelled = false and sicommitstate in (1,2)
					order by dtscheduledtime limit 1 for update;`)
		if err != nil {
			return TransactionOperation{}, err
		}

		if len(ops) == 0 {
			return TransactionOperation{}, nil
		}

		return ops[0], nil
	}, tx)

	return res.(TransactionOperation), err

}
