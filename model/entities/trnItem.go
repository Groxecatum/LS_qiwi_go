package entities

import (
	"encoding/xml"
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
	TotalBonus       int      `json:"totalBonusAmount"    xml:"totalBonusAmount"`
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
	return 0
}

func GetOverallBonusAmount(list []TrnItem, cardAccountType int, includeCampaigns bool) int64 {
	return 0
}
