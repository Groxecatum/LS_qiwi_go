package requests

import (
	"encoding/json"
	"git.kopilka.kz/BACKEND/golang_commons"
	"net/http"
)

type CheckPinRequest struct {
	CellPhoneNum string `json:"cellPhoneNum" xml:"cellPhoneNum"`
	Pin          string `json:"pin" xml:"pin"`
	DeviceId     string `json:"deviceId" xml:"deviceId"`
	Manufacturer string `json:"manufacturer" xml:"manufacturer"`
	Model        string `json:"model" xml:"model"`
	PushToken    string `json:"pushToken" xml:"pushToken"`
}

func ParseCheckPinRequest(r *http.Request) (CheckPinRequest, error) {
	b, err := golang_commons.ParseReqByte(r)
	var req CheckPinRequest
	if err != nil {
		return req, err
	}

	err = json.Unmarshal(b, &req)
	return req, err
}
