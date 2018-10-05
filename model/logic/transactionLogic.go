package logic

import (
	"git.kopilka.kz/BACKEND/golang_commons/errors"
	"git.kopilka.kz/BACKEND/golang_commons/model/entities"
	"github.com/jmoiron/sqlx"
	"math"
	"time"
)

const (
	AMOUNT_PERCENT_WITHOUT_WITHDRAW_PART = 1
	BONUS_PERCENT                        = 2
	FIXED                                = 3
	AMOUNT_PERCENT_WITH_WITHDRAW_PART    = 4
)

func getChargeFeeAmount(m entities.Merchant, bonusAmount, amount, bonusAmountToPay float64, typeId int) int64 {
	chargeFeeAmount := int64(0)

	if typeId == entities.DEFAULT_ACC_TYPE {
		switch m.ChargeFeeType {
		case AMOUNT_PERCENT_WITHOUT_WITHDRAW_PART:
			chargeFeeAmount = int64(amount - bonusAmountToPay*float64(m.ChargeFeeValue)/100)
			break
		case AMOUNT_PERCENT_WITH_WITHDRAW_PART:
			chargeFeeAmount = int64(amount * float64(m.ChargeFeeValue) / 100)
			break
		case BONUS_PERCENT:
			chargeFeeAmount = int64(bonusAmount * float64(m.ChargeFeeValue) / 100)
			break
		case FIXED:
			chargeFeeAmount = int64(m.ChargeFeeValue * 100)
			break
		}

	}

	return chargeFeeAmount
}

func GainBonusesBySrcItemList(tx *sqlx.Tx, crd entities.Card, itemList []entities.TrnItem, bonusAmountToPay float64, m entities.Merchant,
	mt entities.MerchantTerminal, tr entities.TransactionRequest, t entities.Transaction, commitState int, includeBlocked bool, processOnline bool) error {
	amount := entities.GetItemsOverallAmount(itemList)

	scheduledTime := time.Now().AddDate(0, 0, m.BlockDays)

	typesList := entities.GetAccountTypesAndSumsFromItems(itemList)

	for accType, bonusAmount := range typesList {

		clientAccount, err := entities.GetByCardAndType(tx, crd.Id, accType, true)
		if err != nil {
			return err
		}

		if clientAccount.Id == 0 {
			clientAccount, err = entities.CreateAndLinkNew(tx, false, crd.Id, accType, nil)
			if err != nil {
				return err
			}
		}

		chargeFeeAmount := getChargeFeeAmount(m, bonusAmount, amount, bonusAmountToPay, clientAccount.TypeId)

		if m.IsPrepaid {
			merchantAccount, err := entities.GetMerchantAccount(tx, m.Id, true)
			if err != nil {
				return err
			}

			err = entities.RegMerchantOperation(tx, t, entities.TRNOPERTYPE_TERMINAL_TO_TERMINAL_BLOCK,
				bonusAmountToPay, processOnline, merchantAccount, m, mt, nil, tr.Id,
				0, 0, entities.TRNOPERTYPEEX_PAY_WITH_BONUSES,
				entities.COMMITSTATE_INSTANT, 0)
			if err != nil {
				return err
			}

			err = entities.RegMerchantOperation(tx, t, entities.TRNOPERTYPE_TERMINAL_BLOCK_TO_CLIENT,
				bonusAmountToPay, processOnline, merchantAccount, m, mt, &scheduledTime,
				tr.Id, 0, 0,
				entities.TRNOPERTYPEEX_COMMIT, commitState, 0)
			if err != nil {
				return err
			}

		}

		err = entities.RegClientOperation(tx, t, entities.TRNOPERTYPE_TERMINAL_TO_CLIENT_BLOCK,
			bonusAmountToPay, processOnline, crd, clientAccount, m, mt, nil, tr.Id,
			0, chargeFeeAmount, entities.TRNOPERTYPEEX_GAIN_BONUSES,
			entities.COMMITSTATE_INSTANT, 0, "", includeBlocked, false)
		if err != nil {
			return err
		}
		err = entities.RegClientOperation(tx, t, entities.TRNOPERTYPE_CLIENT_BLOCK_TO_ACTIVE,
			bonusAmountToPay, processOnline, crd, clientAccount, m, mt, &scheduledTime, tr.Id,
			0, 0, entities.TRNOPERTYPEEX_COMMIT, commitState,
			0, "", includeBlocked, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func findLastUsedCardId(tx *sqlx.Tx, clientId int) (entities.Card, error) {
	// TODO
	return entities.Card{}, nil
}

func GetLastCardByClientId(tx *sqlx.Tx, clientId int, blockForUpdate bool) (entities.Card, error) {
	var card entities.Card
	cardList, err := entities.GetCardlistByClient(tx, clientId)
	if err != nil {
		return entities.Card{}, err
	}

	if len(cardList) == 0 {
		return entities.Card{}, errors.NotFoundError{}
	}

	if len(cardList) == 1 { /* если у нас только один кард-аккаунт, то без разницы. с какой карты пройдет оплата */
		return card, nil
	}

	card, err = findLastUsedCardId(tx, clientId)
	if err != nil {
		return entities.Card{}, err
	}

	if card.Id == 0 {
		return entities.Card{}, errors.NotFoundError{}
	}

	return card, err

}

func WithdrawBonusesByPriority(tx *sqlx.Tx, crd entities.Card, mt entities.MerchantTerminal, m entities.Merchant, ctr entities.TransactionRequest,
	ct entities.Transaction, amountToPay float64, needCommit int, includeBlocked bool, withdrawFeeAmount int64) error {
	withdrawFeeCalcCompleted := withdrawFeeAmount == -1

	for amountToPay > 0 {
		accountToPay, err := entities.GetAccountForWithdrawByPriority(tx, crd.Id, m.Id)
		if err != nil {
			return err
		}

		if accountToPay.Id == 0 {
			return errors.BalanceError{}
		}

		payOnThisIteration := math.Min(amountToPay, math.Max(accountToPay.Bonuses, 0))

		var withdrawFeeAmountOnThisIteration int64
		if withdrawFeeAmount == -1 && accountToPay.TypeId == entities.DEFAULT_ACC_TYPE {
			withdrawFeeAmountOnThisIteration = int64(float64(amountToPay) * m.WithdrawFeePercent / 100)
		} else if !withdrawFeeCalcCompleted {
			withdrawFeeAmountOnThisIteration = withdrawFeeAmount
			withdrawFeeCalcCompleted = true
		} else {
			withdrawFeeAmountOnThisIteration = 0
		}

		if m.IsPrepaid {
			merchantAccount, err := entities.GetMerchantAccount(tx, mt.MerchantId, true)
			if err != nil {
				return err
			}

			if needCommit == 1 {
				err = entities.RegMerchantOperation(tx, ct, entities.TRNOPERTYPE_CLIENT_TO_TERMINAL_BLOCK,
					payOnThisIteration, true, merchantAccount, m, mt, nil, ctr.Id, 0, 0,
					entities.TRNOPERTYPEEX_PAY_WITH_BONUSES, entities.COMMITSTATE_INSTANT, 0)
				if err != nil {
					return err
				}

				err = entities.RegMerchantOperation(tx, ct, entities.TRNOPERTYPE_TERMINAL_BLOCK_TO_TERMINAL,
					payOnThisIteration, false, merchantAccount, m, mt, nil, ctr.Id,
					withdrawFeeAmountOnThisIteration, 0, entities.TRNOPERTYPEEX_UNDEFINED,
					entities.COMMITSTATE_WAITFORCOMMIT, 0)
			} else {
				err = entities.RegMerchantOperation(tx, ct, entities.TRNOPERTYPE_CLIENT_TO_TERMINAL,
					payOnThisIteration, true, merchantAccount, m, mt, nil, ctr.Id,
					withdrawFeeAmountOnThisIteration, 0, entities.TRNOPERTYPEEX_UNDEFINED,
					entities.COMMITSTATE_WAITFORCOMMIT, 0)

			}
			if err != nil {
				return err
			}
		}

		if needCommit == 1 {
			err = entities.RegClientOperation(tx, ct, entities.TRNOPERTYPE_CLIENT_ACTIVE_TO_CLIENT_BLOCK,
				payOnThisIteration, true, crd, accountToPay, m, mt, nil, ctr.Id, 0, 0,
				entities.TRNOPERTYPEEX_PAY_WITH_BONUSES, entities.COMMITSTATE_INSTANT, 0, "", includeBlocked, false)
			if err != nil {
				return err
			}
			err = entities.RegClientOperation(tx, ct, entities.TRNOPERTYPE_CLIENT_BLOCK_TO_TERMINAL,
				payOnThisIteration, false, crd, accountToPay, m, mt, nil, ctr.Id,
				withdrawFeeAmountOnThisIteration, 0, entities.TRNOPERTYPEEX_UNDEFINED,
				entities.COMMITSTATE_WAITFORCOMMIT, 0, "", includeBlocked, false)
		} else {
			err = entities.RegClientOperation(tx, ct, entities.TRNOPERTYPE_CLIENT_ACTIVE_TO_TERMINAL,
				payOnThisIteration, true, crd, accountToPay, m, mt, nil, ctr.Id,
				withdrawFeeAmountOnThisIteration, 0, entities.TRNOPERTYPEEX_PAY_WITH_BONUSES,
				entities.COMMITSTATE_INSTANT, 0, "", includeBlocked, false)

		}
		if err != nil {
			return err
		}
		amountToPay = amountToPay - payOnThisIteration
	}

	return nil
}
