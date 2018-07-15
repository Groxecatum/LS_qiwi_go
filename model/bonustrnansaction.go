package model

import "encoding/xml"

//body := "<cellPhone>" + account + "</cellPhone>" +
//"<login>" + ApplicationProperties.QIWI_LOGIN + "</login>" +
//"<psw>" + ApplicationProperties.QIWI_PASSWORD + "</psw>" +
//"<description>Пополнение qiwi</description>" +
//"<date>" + new Timestamp(new Date().getTime()) + "</date>" +
//"<ref>" + txn_id + "</ref>" +
//"<amount>" + sum +  "</amount>" +
//"<bonusAmount>" + sum + "</bonusAmount>" +
//"<checkId>" + txn_id + "</checkId>";

type BonusTransaction struct {
	XMLName     xml.Name `xml:"request"`
	FrontEnd    string   `xml:"frontEnd,attr"`
	Type        string   `xml:"type,attr"`
	Account     string   `xml:"cellPhone"`
	Psw         string   `xml:"psw"`
	Description string   `xml:"description"`
	Ref         string   `xml:"ref"`
	Amount      string   `xml:"amount"`
	CheckId     string   `xml:"checkId"`
}

type BonusTransactionResponse struct {
	XMLName       xml.Name `xml:"response"`
	Result        int      `xml:"result,attr"`
	Description   string   `xml:"resultDescription,attr"`
	TransactionId string   `xml:"transactionId"`
}