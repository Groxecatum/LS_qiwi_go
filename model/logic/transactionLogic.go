package logic

import (
	"git.kopilka.kz/BACKEND/golang_commons/model/entities"
	"github.com/jmoiron/sqlx"
)

func GainBonusesBySrcItemList(tx *sqlx.Tx, crd entities.Card, itemList []entities.TrnItem, bonuses int64, terminal entities.MerchantTerminal,
	request entities.TransactionRequest, transaction entities.Transaction, state int, includeBlocked bool) error {
	return nil
}

func GetLastCardByClientId(tx *sqlx.Tx, clientId int, blockForUpdate bool) (*entities.Card, error) {
	return nil, nil
}

func WithdrawBonusesByPriority(tx *sqlx.Tx, crd entities.Card, mt entities.MerchantTerminal, ctr entities.TransactionRequest,
	ct entities.Transaction, bonusAmountToPay int64, needCommit int, includeBlocked bool) error {
	return nil
}
