package entities

import (
	"git.kopilka.kz/BACKEND/golang_commons"
	"git.kopilka.kz/BACKEND/golang_commons/errors"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type Account struct {
	Id                 int   `db:"iid" xml:"id"`
	TypeId             int   `db:"icardaccounttypeid"`
	Bonuses            int64 `db:"nbonuses"`
	BlockedBonuses     int64 `db:"nblockedbonuses"`
	IsPaymentAllowed   bool  `db:"bispaymentallowed"`
	IsTemporaryBlocked bool  `db:"bistemporaryblocked"`
	IsTest             bool  `db:"bistest"`
	IsBlocked          bool  `db:"bblocked"`
}

type AccountChange struct {
	ResultAmount        int64
	ResultBlockedAmount int64
}

const (
	DEFAULT_ACC_TYPE = 1
)

func RegAccountChange(tx *sqlx.Tx, accId int, amountChange, blockedAmountChange int64) (AccountChange, error) {

	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		accChng := AccountChange{}

		rows, err := tx.Query(`update ls.tcardaccounts set nbonuses = nbonuses + $1,
			nblockedbonuses = nblockedbonuses + $2, dtburn = CASE WHEN cat.bburndaterefreshable THEN
			        CASE WHEN cat.iburnafterdays = 0 OR cat.iburnafterdays IS NULL THEN CURRENT_DATE + INTERVAL '1 year ' ELSE CURRENT_DATE + cat.iburnafterdays END
			        ELSE dtburn END
			 FROM ls.tcardaccounttypes cat 
			 where tcardaccounts.iid = $3 AND cat.iid = icardaccounttypeid returning nbonuses, nblockedbonuses`,
			amountChange, blockedAmountChange, accId)
		if err != nil {
			log.Println(err)
			return accChng, err
		}

		defer rows.Close()

		if rows.Next() {
			err = rows.Scan(&accChng.ResultAmount, &accChng.ResultBlockedAmount)
			if err != nil {
				return accChng, err
			}
		}

		return accChng, err
	}, tx)
	return res.(AccountChange), err
}

func GetAccountById(tx *sqlx.Tx, id int, lock bool) (Account, error) {
	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		acc := Account{}

		forUpdStr := ""
		if lock {
			forUpdStr = " FOR UPDATE"
		}

		err := tx.Get(&acc, `select ca.iid, ca.icardaccounttypeid, ca.nbonuses, ca.nblockedbonuses, ca.bispaymentallowed,
				ca.bistemporaryblocked, ca.btest, ca.bblocked from ls.tcardaccounts ca  where iid = $1`+forUpdStr, id)
		if err != nil {
			log.Println(err)
			return acc, err
		}

		return acc, err
	}, tx)
	return res.(Account), err
}

func GetAccountForWithdrawByPriority(tx *sqlx.Tx, cardId, merchantId int) (Account, error) {
	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		acc := Account{}
		err := tx.Get(&acc, `select ca.iid, ca.icardaccounttypeid, ca.nbonuses, ca.nblockedbonuses, ca.bispaymentallowed,
				ca.bistemporaryblocked, ca.btest, ca.bblocked
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

func GetAccListByClientId(tx *sqlx.Tx, clientId int) ([]Account, error) {
	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		accs := []Account{}
		err := tx.Select(&accs, `select ca.iid, ca.icardaccounttypeid, ca.nbonuses, ca.nblockedbonuses, ca.bispaymentallowed,
					ca.bistemporaryblocked, ca.btest, ca.bblocked
				from ls.tcardaccounts ca
                    INNER JOIN ls.tcardaccounts_cards cac ON ca.iid = cac.icardaccid AND ca.btest IS FALSE AND ca.nbonuses > 0
                    INNER JOIN ls.tcards crd ON cac.icardid = crd.iid AND crd.iclientid = $1;`, clientId)

		return accs, err
	}, tx)
	return res.([]Account), err
}

func GetMerchantAccount(tx *sqlx.Tx, merchantId int, forUpdate bool) (Account, error) {
	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		acc := Account{}
		forUpdStr := ""
		if forUpdate {
			forUpdStr = " FOR UPDATE"
		}
		err := tx.Get(&acc, `select a.iid, a.icardaccounttypeid, a.nbonuses, a.nblockedbonuses, a.bispaymentallowed,
					a.bistemporaryblocked, a.btest, a.bblocked from ls.tmerchantaccounts m_a
				INNER JOIN ls.tcardaccounts a on a.iid = m_a.iaccountid
			WHERE imerchantid = $1 and bblocked is not true and bispaymentallowed is true`+forUpdStr, merchantId)

		return acc, err
	}, tx)
	return res.(Account), err
}

func GetByCardAndType(tx *sqlx.Tx, cardId, typeId int, forUpdate bool) (Account, error) {
	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		acc := Account{}
		forUpdStr := ""
		if forUpdate {
			forUpdStr = " FOR UPDATE"
		}
		err := tx.Get(&acc, `SELECT  a.iid, a.icardaccounttypeid, a.nbonuses, a.nblockedbonuses, a.bispaymentallowed,
					a.bistemporaryblocked, a.btest, a.bblocked 
			FROM ls.tcardaccounts a
				INNER JOIN ls.tcardaccounts_cards cac on a.iid = cac.icardaccid and cac.icardid = $1
			WHERE a.icardAccountTypeId = $2 AND a.bblocked IS NOT TRUE`+forUpdStr, cardId, typeId)

		return acc, err
	}, tx)
	return res.(Account), err
}

func CreateNewAccount(tx *sqlx.Tx, cardAccountTypeId int, test bool, externalBurnDate *time.Time) (Account, error) {
	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		var dtBurn time.Time
		if externalBurnDate == nil {
			typ, err := GetAccountTypeById(tx, cardAccountTypeId)
			if err != nil {
				return typ, err
			}
			if typ.BurnDays > 0 {
				dtBurn = time.Now().AddDate(0, 0, typ.BurnDays)
			} else {
				dtBurn = time.Now().AddDate(1, 0, 0)
			}
		} else {
			dtBurn = *externalBurnDate
		}
		acc := Account{TypeId: cardAccountTypeId, IsTest: test, IsPaymentAllowed: true, IsTemporaryBlocked: false, BlockedBonuses: 0,
			Bonuses: 0}
		rows, err := tx.Query(`insert into ls.tCardAccounts (nBonuses, nBlockedBonuses, bblocked, btest, icardaccounttypeid, dtburn) 
									values (0.0, 0.0, false, $1, $2, $3) returning iid`, test, cardAccountTypeId, dtBurn)
		if err != nil {
			log.Println(err)
			return acc, err
		}

		defer rows.Close()

		if rows.Next() {
			err := rows.Scan(&acc.Id)
			if err != nil {
				log.Println(err)
				return acc, err
			}
		}

		return acc, err
	}, tx)

	return res.(Account), err
}

func CreateAndLinkNew(tx *sqlx.Tx, test bool, cardId, cardAccountTypeId int, externalBurnDate *time.Time) (Account, error) {
	res, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		acc, err := CreateNewAccount(tx, cardAccountTypeId, test, externalBurnDate)
		if err != nil {
			log.Println(err)
			return acc, err
		}

		if acc.Id == 0 {
			return acc, errors.DBError{}
		}

		_, err = tx.Exec("insert into ls.tcardaccounts_cards (icardid, icardaccid) values (?, ?);",
			cardId, cardAccountTypeId)

		return acc, err
	}, tx)

	return res.(Account), err
}
