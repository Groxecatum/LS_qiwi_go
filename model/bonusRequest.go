package model

import (
	"encoding/json"
	"io/ioutil"
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
	Id        string `json:"id"`
	IdInCheck int    `json:"idInCheck"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
	Amount    int    `json:"amount"`
	Bonus     int    `json:"bonus"`
}

type BonusRequest struct {
	Type       string    `json:"type"`
	SessionId  string    `json:"sessionId"`
	Amount     int       `json:"amount"`
	CommitType int       `json:"commitType"`
	Date       time.Time `json:"date"`
	Ref        string    `json:"ref"`
	CheckId    string    `json:"checkId"`
	Terminal   int       `json:"terminal"`
	Card       string    `json:"card"`
	SecCode    string    `json:"secCode"`
	Pin        string    `json:"pin"`
	BonusesPay int       `json:"bonusesPay"`
	BonusesAcc int       `json:"bonusesAccumulate"`
	Items      []Item    `json:"items"`
}

func ParseReqByte(r *http.Request) ([]byte, error) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	return b, err
}

func ParseBonusRequest(r *http.Request) (BonusRequest, error) {
	b, err := ParseReqByte(r)
	var req BonusRequest
	if err != nil {
		return req, err
	}

	err = json.Unmarshal(b, &req)
	return req, err
}

func RequestFromBytes(b []byte) (BonusRequest, error) {
	var req BonusRequest

	err := json.Unmarshal(b, &req)
	return req, err
}
