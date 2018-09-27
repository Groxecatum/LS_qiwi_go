package entities

import (
	"database/sql"
	"encoding/xml"
	"git.kopilka.kz/BACKEND/golang_commons"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	TOTAL_ITEM_SID = "#total#"
)

type TrnItem struct {
	XMLName          xml.Name `xml:"item"`
	ItemId           string   `json:"id"        xml:"id"`
	ItemName         string   `json:"itemName"  xml:",chardata"`
	IdInCheck        int      `json:"idInCheck" xml:"idInCheck"`
	Quantity         int      `json:"quantity"  xml:"quantity"`
	Price            int64    `json:"price"     xml:"price"`
	Amount           int64    `json:"amount"    xml:"amount"`
	TotalBonus       int64    `json:"totalBonusAmount"    xml:"totalBonusAmount"`
	Bonus            int64    `json:"bonusAmount"     xml:"bonusAmount"`
	Id               int
	Created          time.Time
	IsActual         bool
	TransactionId    int64
	TrnRequestId     int64
	AccountTypeId    int
	SourceTerminalId int
}

func GetItemsOverallAmount(list []TrnItem) int64 {
	var amount int64
	for _, item := range list {
		amount += item.Amount
	}
	return amount
}

func GetOverallBonusAmount(list []TrnItem, cardAccountType int, includeCampaigns bool) int64 {
	var amount int64
	for _, item := range list {
		amount += item.TotalBonus
	}
	return amount
}

func (item *TrnItem) save() error {
	_, err := golang_commons.DoX(func(tx *sqlx.Tx) (interface{}, error) {
		rows, err := tx.Query("INSERT INTO ls.tsrctransactionitems "+
			" (bitransactionid, sitemsid, sitemname, nitemquantity, npriceperitem, "+
			"  isactual, bitrnrequestid, nitemquantitychange, nbonusamount, "+
			"  nbonusamountchange, namount, namountchange, iidincheck, dtcreated, iAccountTypeId, isourceterminal) \n"+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) returning biid",
			item.TransactionId, item.ItemId, item.ItemName, item.Quantity, item.Price, item.IsActual, item.TrnRequestId,
			item.Quantity, item.Bonus, item.TotalBonus, item.Amount, item.Amount, item.IdInCheck, item.Created, item.AccountTypeId,
			item.SourceTerminalId)

		if rows.Next() {
			err = rows.Scan(&item.Id)
		}

		return nil, err
	}, nil)
	return err
}

func SetTrnRequestId(list []TrnItem, id int64) {
	for _, item := range list {
		item.TrnRequestId = id
	}
}

func SetTrnId(list []TrnItem, id int64) {
	for _, item := range list {
		item.TransactionId = id
	}
}

func SaveTrnItems(tx *sql.Tx, list []TrnItem) error {
	var err error
	for _, item := range list {
		err := item.save()
		if err != nil {
			break
		}
	}
	return err
}
