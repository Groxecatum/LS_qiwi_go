package service

import (
	"net/http"
	"encoding/xml"
	"../model"
	"strconv"
	"fmt"
	"log"
	"io/ioutil"
	"../config"
	"bytes"
)

const clientVerifyServlet = "msrv_AllowBonusTransaction"
const bonusTransactionServlet = "mrct_BonusTransaction"

func isTrnAllowed(platform, account string) (int, string, error) {
	hc := http.Client{}
	vd := model.VerifyData{FrontEnd: platform, Type: bonusTransactionServlet, Account: account};
	buf, err := xml.Marshal(vd)

	req, err := http.NewRequest("POST", config.Config.LSURL + config.Config.LSPATH + clientVerifyServlet,
		bytes.NewBuffer(buf))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/xml")

	resp, err := hc.Do(req)

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	LSResp := model.VerifyDataResponse{}
	xml.Unmarshal(respBody, &LSResp)

	return LSResp.Result, LSResp.Description, err

}

func newQiwiTrn(platform, account, txn_id string, sum float64) (int, string, string, error) {
	hc := http.Client{}
	bt := model.BonusTransaction{Psw: config.Config.LS_PASSWORD, FrontEnd: platform, Type: bonusTransactionServlet,
		Account: account, Description: "Qiwi trn", Ref: txn_id, CheckId: txn_id, Amount: strconv.Itoa(int(sum * 100))};

	buf, err := xml.Marshal(bt)

	req, err := http.NewRequest("POST", config.Config.LSURL + config.Config.LSPATH + clientVerifyServlet,
		bytes.NewBuffer(buf))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/xml")

	resp, err := hc.Do(req)

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	LSResp := model.BonusTransactionResponse{}
	err = xml.Unmarshal(respBody, &LSResp)

	return LSResp.Result, LSResp.Description, LSResp.TransactionId, err
}

func mapLSCode(code int) int {
    return 1
}

func pay(account, txn_id string, sum float64) model.Response {
	code, descr, err := isTrnAllowed("web", account);
	if err != nil || code != 0 {
		return model.NewErrorResponse(mapLSCode(code), descr)
	}
	code, descr, trnId, err := newQiwiTrn("web", account, txn_id, sum * 100);
	if err != nil || code != 0 {
		return model.NewErrorResponse(mapLSCode(code), descr)
	}
	return model.Response{Result: 0, Sum: fmt.Sprintf("%.2f", sum), Prv_txn: trnId, Osmp_txn_id: txn_id, Comment: "OK"}
}

func check(account, txn_id string) model.Response {
	code, descr, err := isTrnAllowed("web", account);
	if err != nil || code != 0 {
		return model.NewErrorResponse(mapLSCode(code), descr)
	}
	return model.Response{Result: 0, Sum: "0", Prv_txn: "0", Osmp_txn_id: txn_id, Comment: "OK"}
}

func NewTransaction(w http.ResponseWriter, r *http.Request) {
	keys, _ := r.URL.Query()["command"]
	txn_id, _ := r.URL.Query()["txn_id"]
	account, _ := r.URL.Query()["account"]

	switch (keys[0]) {
		case "pay":
			sum, _ := r.URL.Query()["sum"]
			sumf, err := strconv.ParseFloat(sum[0], 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			xml.NewEncoder(w).Encode(pay(account[0], txn_id[0], sumf));
		default:
			xml.NewEncoder(w).Encode(check(account[0], txn_id[0]));
	}
}
