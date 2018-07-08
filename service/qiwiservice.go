package service

import (
	"net/http"
	"encoding/xml"
	"../model"
	"strconv"
	"fmt"
	"strings"
	"log"
	"io/ioutil"
	"../config"
)

const clientVerifyServlet = "msrv_AllowBonusTransaction"
const bonusTransactionServlet = "mrct_BonusTransaction"

func isTrnAllowed(platform, account, txn_id, srcAddr string) (int, string, error) {
	hc := http.Client{}
	

	req, err := http.NewRequest("POST", config.Config.LSURL + config.Config.LSPATH + clientVerifyServlet,
		strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/xml")

	resp, err := hc.Do(req)

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	resp := model.LSResponse{}
	xml.Unmarshal(respBody, &resp)

	return resp.Code, resp.Description,

}

func newQiwiTrn(platform, account, txn_id, srcAddr string, sum float64) (int, string, int, error) {

	//body := "<cellPhone>" + account + "</cellPhone>" +
	//	"<login>" + ApplicationProperties.QIWI_LOGIN + "</login>" +
	//	"<psw>" + ApplicationProperties.QIWI_PASSWORD + "</psw>" +
	//	"<description>Пополнение qiwi</description>" +
	//	"<date>" + new Timestamp(new Date().getTime()) + "</date>" +
	//	"<ref>" + txn_id + "</ref>" +
	//	"<amount>" + sum +  "</amount>" +
	//	"<bonusAmount>" + sum + "</bonusAmount>" +
	//	"<checkId>" + txn_id + "</checkId>";
	//ResponseEntity<String> response = sendPostXMLRequest(BONUS_TRANSACTION_SERVLET,
	//	formLSRequest(BONUS_TRANSACTION_SERVLET, originatorPlatform, body, srcAddr));
	//String trnId = getFieldFromXMLResponse(response, "transactionId", BONUS_TRANSACTION_SERVLET, txn_id);
	//return (!trnId.isEmpty() ?
	//new QiwiResponse(txn_id, SUCCESS, "OK", trnId, sum.toString()) :
	//new QiwiResponse(txn_id, FAIL, "Trn failed", "0", "0"));

	//code, descr, err :=
	//if err != nil || code != 0 {
	//	return model.NewErrorResponse(mapLSCode(code), descr)
	//}
}

func mapLSCode(code int) int {
    return 1
}

func pay(account, txn_id, srcAddr string, sum float64) model.Response {
	code, descr, err := isTrnAllowed("web", account, txn_id, srcAddr);
	if err != nil || code != 0 {
		return model.NewErrorResponse(mapLSCode(code), descr)
	}
	code, descr, trnId, err := newQiwiTrn("web", account, txn_id, srcAddr, sum * 100);
	if err != nil || code != 0 {
		return model.NewErrorResponse(mapLSCode(code), descr)
	}
	return model.Response{Result: 0, Sum: fmt.Sprintf("%.2f", sum), Prv_txn: strconv.Itoa(trnId), Osmp_txn_id: txn_id, Comment: "OK"}
}

func check(account, txn_id, srcAddr string) model.Response {
	code, descr, err := isTrnAllowed("web", account, txn_id, srcAddr);
	if err != nil || code != 0 {
		return model.NewErrorResponse(mapLSCode(code), descr)
	}
	return model.Response{Result: 0, Sum: "0", Prv_txn: "0", Osmp_txn_id: txn_id, Comment: "OK"}
}

func NewTransaction(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["command"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch (keys[0]) {
		case "pay":
			xml.NewEncoder(w).Encode(pay());
		default:
			xml.NewEncoder(w).Encode(check());
	}
}
