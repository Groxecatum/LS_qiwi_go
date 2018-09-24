package model

import (
	"encoding/json"
	"encoding/xml"
	"git.kopilka.kz/BACKEND/golang_commons"
	"net/http"
	"time"
)

/*  <request type=”mrct_BonusTransaction”>
      <sessionId>[Идентификатор сессии точки продаж]</sessionId> ( ИЛИ <login>[логин]</login><psw>[пароль]</psw> )
      <date>[дата операции]</date>
      <ref>[номер запроса]</ref>
      <checkId>[номер чека. Default = ref]</checkId>
	  <terminal>[номер кассы ( =1, если отсутствует)]</terminal>
	  <card>[Номер карты]</card>
	  <pin>[пин-код клиента - если требуется оплата бонусами]</pin>
	  <bonusAmountToPay>[сумма бонусов к оплат. Если отсутствует и указан пин, то = bonusAmount) ]</ bonusAmountToPay>
	  <zRepId>[z-report id]</zRepId>
	  <batchPeriodId>[batch-период, может быть пустым]</batchPeriodId>
	  <needCommit>[=1 если для операции необходим последующий вызов Commit, по умолчанию = 0]</needCommit >
	  <items> [список товаров по чеку]
	     <item id=”[id]” idInCheck=”[порядковый номер в чеке]” quantity="[кол-во]" price="[стоимость единицы]" amount=”[общая сумма]”  bonusAmount=”[общая сумма бонусов к начислению]”>описание</item>
	     <item id=”[id]” idInCheck=”[порядковый номер в чеке]” quantity="[кол-во]" price="[стоимость единицы]" amount=”[общая сумма]” bonusAmount=”[общая сумма бонусов к начислению]>описание</item>
	     ...
	  </items>
    </request>
 Ответ:
    <response result=”[код результата. 0 - Ok]” resultDescr=”[описание результата]”>
      <transactionID>[Идентификатор транзакции в системе лояльности]</transactionID>
      <resultBonusAmount>[сумма бонусов, которая предполагается на счете после операции]<resultBonusAmount>
*   </response>*/

type Item struct {
	Id        string `json:"id"        xml:"id"`
	IdInCheck int    `json:"idInCheck" xml:"idInCheck"`
	Quantity  int    `json:"quantity"  xml:"quantity"`
	Price     int    `json:"price"     xml:"price"`
	Amount    int    `json:"amount"    xml:"amount"`
	Bonus     int    `json:"bonus"     xml:"bonus"`
}

type BonusRequest struct {
	CustomRequest
	SessionId  string    `json:"sessionId"  xml:"sessionId"`
	Amount     int       `json:"amount"     xml:"amount"`
	CommitType int       `json:"commitType" xml:"commitType"`
	Date       time.Time `json:"date"       xml:"date"`
	Ref        string    `json:"ref"        xml:"ref"`
	CheckId    string    `json:"checkId"    xml:"checkId"`
	Terminal   int       `json:"terminal"   xml:"terminal"`
	Card       string    `json:"card"       xml:"card"`
	SecCode    string    `json:"secCode"    xml:"secCode"`
	Pin        string    `json:"pin"        xml:"pin"`
	BonusesPay int       `json:"bonusesPay" xml:"bonusesPay"`
	BonusesAcc int       `json:"bonuses"    xml:"bonuses"`
	Items      []Item    `json:"items"      xml:"items"`
}

type BonusResponse struct {
	CustomResponse
	TransactionId     string `json:"transactionId"     xml:"transactionId"`
	ResultBonusAmount int    `json:"resultBonusAmount" xml:"resultBonusAmount"`
	ClientName        string `json:"clientName"        xml:"clientName"`
	TokenCount        int    `json:"tokenCount"        xml:"tokenCount"`
	Discount          int    `json:"discount"          xml:"discount"`
	ReqChallenge      string `json:"reqChallenge"      xml:"reqChallenge"`
}

func ParseBonusRequest(r *http.Request) (BonusRequest, error) {
	b, err := golang_commons.ParseReqByte(r)
	if err != nil {
		return BonusRequest{}, err
	}

	return RequestFromBytes(b, golang_commons.GetFormatByRequiest(r))
}

func RequestFromBytes(b []byte, format string) (BonusRequest, error) {
	var req BonusRequest
	switch format {
	case "json":
		return req, json.Unmarshal(b, &req)
	default:
		return req, xml.Unmarshal(b, &req)
	}

}

func (req *BonusRequest) IsPayment() bool {
	return false
}
