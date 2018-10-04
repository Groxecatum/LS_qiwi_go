package entities

import (
	"git.kopilka.kz/BACKEND/golang_commons/db"
	"github.com/jmoiron/sqlx"
	"log"
)

type Merchant struct {
	Id                    int     `db:"iid"`
	CreateVirtualUser     bool    `db:"bcreatevirtualuser"`
	AllowPayWithoutPin    bool    `db:"ballowpaywithoutpin"`
	IsPrepaid             bool    `db:"bisprepaid"`
	WithdrawFeePercent    float64 `db:"nwithdrawfeepercent"`
	BlockDays             int     `db:"siblockdays"`
	ChargeFeeType         int     `db:"sichargefeetype"`
	ChargeFeeValue        float32 `db:"nchargefeevalue"`
	AllowNegativeDecrease bool    `db:"ballownegativedecrease"`
	AllowNegativeBalance  bool    `db:"ballownegativebalance"`
}

func GetMerchantById(tx *sqlx.Tx, id int) (Merchant, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		mrct := Merchant{}
		err := tx.Get(&mrct, `select * from ls.tcards where iid = $1`, id)
		if err != nil {
			log.Println(err)
			return mrct, err
		}

		return mrct, err
	}, tx)
	return res.(Merchant), err
}

// без всяких жирных текстовых полей
func GetMerchantDataForTransaction(tx *sqlx.Tx, merchantId int) (Merchant, error) {
	res, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		mrct := Merchant{}
		err := tx.Get(&mrct, `select iid, bcreatevirtualuser, ballowpaywithoutpin, bisprepaid, nwithdrawfeepercent,
											siblockdays, sichargefeetype, nchargefeevalue, ballownegativedecrease, ballownegativebalance
										from ls.tmerchants where iid = $1`, merchantId)
		if err != nil {
			log.Println(err)
			return mrct, err
		}

		return mrct, err
	}, tx)
	return res.(Merchant), err
}
