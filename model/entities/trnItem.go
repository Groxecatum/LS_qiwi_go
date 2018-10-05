package entities

import (
	"encoding/xml"
	"git.kopilka.kz/BACKEND/golang_commons/db"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	TOTAL_ITEM_SID = "#total#"
)

type TrnItem struct {
	XMLName   xml.Name `xml:"item"`
	ItemId    string   `json:"id"        xml:"id,attr"`
	IdInCheck int      `json:"idInCheck" xml:"idInCheck,attr"`
	Quantity  float64  `json:"quantity"  xml:"quantity,attr"`
	Price     float64  `json:"price"     xml:"price,attr"`
	Amount    float64  `json:"amount"    xml:"amount,attr"`
	//TotalBonus       int64  `json:"totalBonusAmount"    xml:"totalBonusAmount,attr"`
	Bonus            float64 `json:"bonusAmount"     xml:"bonusAmount,attr"`
	ItemName         string  `json:"itemName"  xml:",chardata"`
	Id               int
	Created          time.Time
	IsActual         bool
	TransactionId    int64
	TrnRequestId     int64
	AccountTypeId    int
	SourceTerminalId int
}

func GetItemsOverallAmount(list []TrnItem) float64 {
	var amount float64
	for _, item := range list {
		amount += item.Amount
	}
	return amount
}

func GetOverallBonusAmount(list []TrnItem) float64 {
	var amount float64
	for _, item := range list {
		amount += item.Bonus //TotalBonus
	}
	return amount
}

func (item *TrnItem) save(tx *sqlx.Tx) error {
	_, err := db.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		rows, err := tx.Query("INSERT INTO ls.tsrctransactionitems "+
			" (bitransactionid, sitemsid, sitemname, nitemquantity, npriceperitem, "+
			"  isactual, bitrnrequestid, nitemquantitychange, nbonusamount, "+
			"  nbonusamountchange, namount, namountchange, iidincheck, dtcreated, iAccountTypeId, isourceterminal) \n"+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) returning biid",
			item.TransactionId, item.ItemId, item.ItemName, item.Quantity, item.Price*100, item.IsActual, item.TrnRequestId,
			item.Quantity, item.Bonus*100, item.Bonus*100, item.Amount*100, item.Amount*100, item.IdInCheck, item.Created, item.AccountTypeId,
			item.SourceTerminalId)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		if rows.Next() {
			err = rows.Scan(&item.Id)
		}

		return nil, err
	}, tx)
	return err
}

func SetTrnRequestId(list []TrnItem, id int64) {
	for i, _ := range list {
		list[i].TrnRequestId = id
	}
}

func SetTrnId(list []TrnItem, id int64) {
	for i, _ := range list {
		list[i].TransactionId = id
	}
}

func SetActual(list []TrnItem, val bool) {
	for i, _ := range list {
		list[i].IsActual = val
	}
}

func SetDtCreated(list []TrnItem, val time.Time) {
	for i, _ := range list {
		list[i].Created = val
	}
}

func SetAccountTypeId(list []TrnItem, id int) {
	for i, _ := range list {
		list[i].AccountTypeId = id
	}
}

func SetTerminalId(list []TrnItem, id int) {
	for i, _ := range list {
		list[i].SourceTerminalId = id
	}
}

func SaveTrnItems(tx *sqlx.Tx, list []TrnItem) error {
	var err error
	for _, item := range list {
		err := item.save(tx)
		if err != nil {
			break
		}
	}
	return err
}

func GetAccountTypesAndSumsFromItems(list []TrnItem) map[int]float64 {
	res := make(map[int]float64)
	for _, item := range list {
		res[item.AccountTypeId] = res[item.AccountTypeId] + item.Bonus
	}
	return res
}
