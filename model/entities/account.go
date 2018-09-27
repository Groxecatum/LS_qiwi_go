package entities

import (
	"git.kopilka.kz/BACKEND/golang_commons"
	"git.kopilka.kz/BACKEND/golang_commons/errors"
	"github.com/jmoiron/sqlx"
	"log"
)

type Account struct {
	Id                 int   `db:"iid" xml:"id"`
	TypeId             int   `db:"sitype"`
	Bonuses            int64 `db:"nbonuses"`
	BlockedBonuses     int64 `db:"nblockedbonuses"`
	IsPaymentAllowed   bool  `db:"bispaymentallowed"`
	IsTemporaryBlocked bool  `db:"nbonuses"`
}

const (
	DEFAULT_ACC_TYPE = 1
)

func GetAccountForWithdrawByPriority(tx *sqlx.Tx, cardId, merchantId int) (Account, error) {
	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		acc := Account{}
		err := tx.Get(&acc, `select ca.*
			from ls.tcardaccounts ca 
				INNER JOIN ls.tcardaccounts_cards cac on ca.iid = cac.icardaccid and cac.icardid = $1 
				INNER JOIN ls.tcardaccounttypes cat ON cat.iid = ca.icardaccounttypeid 
				left join ls.tCardAccountTypeMerchants catm ON ca.iCardAccountTypeId = catm.iCardAccountTypeId and catm.imerchantid = $2 
			where (catm.iid is not null or ca.iCardAccountTypeId=1 or cat.bforallmerchants is true) and ca.nbonuses > 0  AND ca.bblocked IS NOT TRUE 
			order by case when ca.iCardAccountTypeId = 1 then 99999 
				else (select count(*) from ls.tCardAccountTypeMerchants where icardaccounttypeid = ca.iCardAccountTypeId ) end  
			limit 1 for update of ca;`, cardId, merchantId) // с флагом bforallmerchants будет высший приоритет, если не вносить записей в ls.tCardAccountTypeMerchants

		if err != nil {
			log.Println(err)
			return acc, err
		}

		if acc.Id == 0 {
			return acc, errors.InsufficientFundsErr
		}

		if (!acc.IsPaymentAllowed) || acc.IsTemporaryBlocked {
			return acc, errors.BlockedErr
		}

		return acc, err
	}, tx)
	return res.(Account), err
}
