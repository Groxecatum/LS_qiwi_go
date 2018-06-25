package service

import (
	"net/http"
	"encoding/xml"
)

func isTrnAllowed(platform string, account string, txn_id string, srcAddr string) {

}

func pay() {
	if (!LSInteraction.isTrnAllowed("web", data.get("account"), data.get("txn_id"), srcAddr)) throw new NotAllowedException();
	return LSInteraction.newQiwiTrn("web", data.get("account"), data.get("txn_id"), Double.parseDouble(data.get("sum")), srcAddr);
}

func check() {
	if (!LSInteraction.isTrnAllowed("web", data.get("account"), data.get("txn_id"), srcAddr)) throw new NotAllowedException();
	return new QiwiResponse(data.get("txn_id"), 0, "OK", "0", data.get("sum"));
}

func NewTransaction(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["command"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch (keys[0]) {
	case "pay":
		xml.NewEncoder(w).Encode(pay);
	default:
		xml.NewEncoder(w).Encode(check);
	}
}
