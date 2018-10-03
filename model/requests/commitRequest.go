package requests

import "time"

type CommitRequest struct {
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
	TransactionId string `json:"transactionId"     xml:"transactionId"`
}
