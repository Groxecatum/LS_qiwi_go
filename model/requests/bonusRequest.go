package requests

import (
	"encoding/json"
	"encoding/xml"
	"git.kopilka.kz/BACKEND/golang_commons"
	"git.kopilka.kz/BACKEND/golang_commons/model/entities"
	"net/http"
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

type BonusRequest struct {
	CustomRequest
	SessionId       string              `json:"sessionId"        xml:"sessionId"`
	Login           string              `json:"login"            xml:"login"`
	Password        string              `json:"psw"              xml:"psw"`
	Amount          float64             `json:"amount"           xml:"amount"`
	Date            entities.CustomTime `json:"date"             xml:"date"`
	Ref             string              `json:"ref"              xml:"ref"`
	CheckId         string              `json:"checkId"          xml:"checkId"`
	Descr           string              `json:"description"      xml:"description"`
	Terminal        int                 `json:"terminal"         xml:"terminal"`
	Card            string              `json:"card"             xml:"card"`
	CellPhone       string              `json:"cellPhone"        xml:"cellPhone"`
	SecCode         string              `json:"secCode"          xml:"secCode"`
	Pin             string              `json:"pin"              xml:"pin"`
	BonusesPay      float64             `json:"bonusAmountToPay" xml:"bonusAmountToPay"`
	NeedCommit      int                 `json:"needCommit"       xml:"needCommit"`
	BonusesAcc      float64             `json:"bonusAmount"      xml:"bonusAmount"`
	SecureHashCode  string              `json:"secureHashCode"   xml:"secureHashCode"`
	SecureShortCode string              `json:"secureShortCode"  xml:"secureShortCode"`
	Items           []entities.TrnItem  `json:"items"            xml:"items>item"`
	ZRepId          string              `json:"zRepId"           xml:"zRepId"`
	BatchPeriodId   string              `json:"batchPeriodId"    xml:"batchPeriodId"`
}

type BonusResponse struct {
	CustomResponse
	TransactionId     int64   `json:"transactionId"     xml:"transactionId"`
	ResultBonusAmount float64 `json:"resultBonusAmount" xml:"resultBonusAmount"`
	ClientName        string  `json:"clientName"        xml:"clientName"`
	TokenCount        int     `json:"tokenCount"        xml:"tokenCount"`
	Discount          int     `json:"discount"          xml:"discount"`
	ReqChallenge      string  `json:"reqChallenge"      xml:"reqChallenge"`
}

func ParseBonusRequest(r *http.Request) (BonusRequest, error) {
	b, err := golang_commons.ParseReqByte(r)
	if err != nil {
		return NewBonusRequestStruct(), err
	}

	return BonusRequestFromBytes(b, golang_commons.GetFormatByRequest(r))
}

func BonusRequestFromBytes(b []byte, format string) (BonusRequest, error) {
	req := NewBonusRequestStruct()
	switch format {
	case "json":
		err := json.Unmarshal(b, &req)
		return req, err
	default:
		err := xml.Unmarshal(b, &req)
		return req, err
	}
}

func (req *BonusRequest) IsPayment() bool {
	return req.Pin != "" && req.Type == "mrct_PayWithBonuses"
}

func NewBonusRequestStruct() BonusRequest {
	return BonusRequest{Terminal: 1}
}
