package requests

import (
	"encoding/json"
	"encoding/xml"
	"git.kopilka.kz/BACKEND/golang_commons"
	"net/http"
	"time"
)

type CommitRequest struct {
	CustomRequest
	SessionId     string    `json:"sessionId"     xml:"sessionId"`
	Login         string    `json:"login"         xml:"login"`
	Password      string    `json:"psw"           xml:"psw"`
	Date          time.Time `json:"date"          xml:"date"`
	Terminal      int       `json:"terminal"      xml:"terminal"`
	InitialRef    string    `json:"initialRef"    xml:"initialRef"`
	Commit        int       `json:"commit"        xml:"commit"`
	ReqChallenge  string    `json:"reqChallenge"  xml:"reqChallenge"`
	AuthorizeCode string    `json:"authorizeCode" xml:"authorizeCode"`
}

type CommitResponse struct {
	CustomResponse
	TransactionId int64 `json:"transactionId"     xml:"transactionId"`
}

func NewCommitRequestStruct() CommitRequest {
	return CommitRequest{Terminal: 1, Commit: 1}
}

func CommitRequestFromBytes(b []byte, format string) (CommitRequest, error) {
	req := NewCommitRequestStruct()
	switch format {
	case "json":
		err := json.Unmarshal(b, &req)
		return req, err
	default:
		err := xml.Unmarshal(b, &req)
		return req, err
	}
}

func ParseCommitRequest(r *http.Request) (CommitRequest, error) {
	b, err := golang_commons.ParseReqByte(r)
	if err != nil {
		return NewCommitRequestStruct(), err
	}

	return CommitRequestFromBytes(b, golang_commons.GetFormatByRequest(r))
}
