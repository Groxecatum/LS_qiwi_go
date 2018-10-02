package logic

import (
	"git.kopilka.kz/BACKEND/golang_commons/errors"
	"git.kopilka.kz/BACKEND/golang_commons/model/entities"
	"github.com/jmoiron/sqlx"
	"math"
)

func GainBonusesBySrcItemList(tx *sqlx.Tx, crd entities.Card, itemList []entities.TrnItem, bonuses int64, terminal entities.MerchantTerminal,
	request entities.TransactionRequest, transaction entities.Transaction, state int, includeBlocked bool) error {
	amount := entities.GetItemsOverallAmount(itemList)

	scheduledTime := CommonFunctions.getScheduledTime(initialTerminal.getMerchant(conn).getBlockDays());
	HashMap<Integer, HashMap<Integer,Long>> operationsDataWithTerminals = itemList.getOperationsDataWithTerminals();
	for (Map.Entry<Integer, HashMap<Integer,Long>> entry : operationsDataWithTerminals.entrySet()) {
		Account clientAccount;
		// Получение, или создание аккаунта для начисления бонусов
		int cardAccountTypeId = entry.getKey();
		for (Map.Entry<Integer, Long> data : entry.getValue().entrySet()) {

			MerchantTerminal mt = new MerchantTerminal(conn, data.getKey());
			long bonusAmount = data.getValue();

			clientAccount = Account.getByCardAndType(conn, crd.getId(), cardAccountTypeId, true);
			if (clientAccount == null) {
				clientAccount = Account.createAndLinkNew(conn, Boolean.FALSE,
					crd.getId(), cardAccountTypeId);
			}


			long chargeFeeAmount = getChargeeFeeAmount(mt.getMerchant(conn), bonusAmount,
				amount, bonusAmountToPay, clientAccount.getCardAccountTypeId());

			if (mt.getMerchant().getPrepaid()) {
				Account merchantAccount = Account.getMerchantAccount(conn, mt.getMerchantId(), true);

				TrnOperation.regMerchantOperation(conn, ct, TrnOperation.TRNOPERTYPE_TERMINAL_TO_TERMINAL_BLOCK,
					bonusAmount, false, merchantAccount, mt, null, ctr.getId(),
					0L, 0L, TrnOperation.TRNOPERTYPEEX_PAY_WITH_BONUSES,
					TrnOperation.COMMITSTATE_INSTANT, 0);

				TrnOperation.regMerchantOperation(conn, ct, TrnOperation.TRNOPERTYPE_TERMINAL_BLOCK_TO_CLIENT,
					bonusAmount, false, merchantAccount, mt, scheduledTime,
					ctr.getId(), 0L, 0L,
					TrnOperation.TRNOPERTYPEEX_COMMIT, trnReqCommitState, 0);

			}

			TrnOperation.regClientOperation(conn, ct, TrnOperation.TRNOPERTYPE_TERMINAL_TO_CLIENT_BLOCK,
				bonusAmount, false, crd, clientAccount, mt, null, ctr.getId(),
				0L, chargeFeeAmount, TrnOperation.TRNOPERTYPEEX_GAIN_BONUSES,
				TrnOperation.COMMITSTATE_INSTANT, 0, "", includeBlockedCards);
			TrnOperation.regClientOperation(conn, ct, TrnOperation.TRNOPERTYPE_CLIENT_BLOCK_TO_ACTIVE,
				bonusAmount, false, crd, clientAccount, mt,
				scheduledTime, ctr.getId(), 0L, 0L,
				TrnOperation.TRNOPERTYPEEX_COMMIT,
				trnReqCommitState, 0, "", includeBlockedCards);
		}
	}
}

func GetLastCardByClientId(tx *sqlx.Tx, clientId int, blockForUpdate bool) (entities.Card, error) {
	var card entities.Card
	cardList := GetCardlistByClient(tx, clientId);
	if len(cardList) == 0 {
		return entities.Card{}, errors.NotFoundError{}
	}

	if len(cardList) == 1 { /* если у нас только один кард-аккаунт, то без разницы. с какой карты пройдет оплата */
		return card, nil
	}

	// Учитываем, что у клиента могут быть 2+ карты на 1 кард-аккаунте
	caList := GetAccListByClientId(tx, clientId);
	card, err := entities.GetCardById(tx, cardList[0].Id, blockForUpdate)

	if len(caList) == 1 { /* если у нас только один кард-аккаунт, то без разницы. с какой карты пройдет оплата */
		return card, nil
	}

	cardId, err := findLastUsedCardId(tx, clientId)
	if err != nil {
		return  entities.Card{}, err
	}

	if cardId == 0 {
		return entities.Card{}, errors.NotFoundError{}
	}

	card, err = entities.GetCardById(tx, cardId, blockForUpdate)
	if err != nil {
		return  entities.Card{}, err
	}

}

func WithdrawBonusesByPriority(tx *sqlx.Tx, crd entities.Card, mt entities.MerchantTerminal, m entities.Merchant, ctr entities.TransactionRequest,
	ct entities.Transaction, bonusAmountToPay int64, needCommit int, includeBlocked bool) error {
	withdrawFeeCalcCompleted := withdrawFeeAmount == -1;

	for (amount > 0) {
		accountToPay, err := entities.GetAccountForWithdrawByPriority(tx, crd.Id, m.Id);
		if err != nil {
			return err
		}

		if (accountToPay.Id == 0) {
			return errors.BalanceError{}
		}

		payOnThisIteration := math.Min(amount, math.Max(accountToPay.getBonuses(), 0));

		var withdrawFeeAmountOnThisIteration int64
		if (withdrawFeeAmount == -1 && accountToPay.getCardAccountTypeId() == entities.DEFAULT_ACC_TYPE) {
			withdrawFeeAmountOnThisIteration = (long) (amount * mt.getMerchant(conn).getWithdrawFeePercent() / 100)
		} else if (!withdrawFeeCalcCompleted) {
			withdrawFeeAmountOnThisIteration = withdrawFeeAmount
			withdrawFeeCalcCompleted = true
		} else {
			withdrawFeeAmountOnThisIteration = 0
		}

		if (mt.getMerchant(conn).getPrepaid()) {
			merchantAccount := Account.getMerchantAccount(tx, mt.getMerchantId(), true);
			if (needCommit == 1) {
				TrnOperation.regMerchantOperation(tx, ct, TrnOperation.TRNOPERTYPE_CLIENT_TO_TERMINAL_BLOCK,
					payOnThisIteration, true, merchantAccount, mt, null, ctr.getId(), 0, 0,
					TrnOperation.TRNOPERTYPEEX_PAY_WITH_BONUSES, TrnOperation.COMMITSTATE_INSTANT, 0);
				TrnOperation.regMerchantOperation(tx, ct, TrnOperation.TRNOPERTYPE_TERMINAL_BLOCK_TO_TERMINAL,
					payOnThisIteration, false, merchantAccount, mt, null, ctr.getId(),
					withdrawFeeAmountOnThisIteration, 0, TrnOperation.TRNOPERTYPEEX_UNDEFINED,
					TrnOperation.COMMITSTATE_WAITFORCOMMIT, 0);
			} else {
				TrnOperation.regMerchantOperation(tx, ct, TrnOperation.TRNOPERTYPE_CLIENT_TO_TERMINAL,
					payOnThisIteration, false, merchantAccount, mt, null, ctr.getId(),
					withdrawFeeAmountOnThisIteration, 0, TrnOperation.TRNOPERTYPEEX_UNDEFINED,
					TrnOperation.COMMITSTATE_WAITFORCOMMIT, 0);
			}
		}

		if (needCommit == 1) {
			TrnOperation.regClientOperation(tx, ct, TrnOperation.TRNOPERTYPE_CLIENT_ACTIVE_TO_CLIENT_BLOCK,
				payOnThisIteration, true, crd, accountToPay, mt, null, ctr.getId(), 0, 0,
				TrnOperation.TRNOPERTYPEEX_PAY_WITH_BONUSES, TrnOperation.COMMITSTATE_INSTANT, 0, "", includeBlocked);
			TrnOperation.regClientOperation(tx, ct, TrnOperation.TRNOPERTYPE_CLIENT_BLOCK_TO_TERMINAL,
				payOnThisIteration, false, crd, accountToPay, mt, null, ctr.getId(),
				withdrawFeeAmountOnThisIteration, 0, TrnOperation.TRNOPERTYPEEX_UNDEFINED,
				TrnOperation.COMMITSTATE_WAITFORCOMMIT, 0, "", includeBlocked);
		} else {
			TrnOperation.regClientOperation(tx, ct, TrnOperation.TRNOPERTYPE_CLIENT_ACTIVE_TO_TERMINAL,
				payOnThisIteration, true, crd, accountToPay, mt, null, ctr.getId(),
				withdrawFeeAmountOnThisIteration, 0, TrnOperation.TRNOPERTYPEEX_PAY_WITH_BONUSES,
				TrnOperation.COMMITSTATE_INSTANT, 0, "", includeBlocked);

		}
		amount = amount - payOnThisIteration
	}
}
