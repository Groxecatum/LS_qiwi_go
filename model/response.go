package model

import "encoding/xml"

//    <? xml version = "1.0"; encoding="UTF-8"?>
//    <response>
//        <osmp_txn_id>1234567</osmp_txn_id>
//        <result>0</result>
//        <comment></comment>
//        <prv_txn>2016</prv_txn>
//        <sum>500.00</sum>
//    </response>

type Response struct {
	XMLName     xml.Name `xml:"response"`
	Osmp_txn_id string   `xml:"osmp_txn_id"`
	Result      int		 `xml:"result"`
	Comment     string	 `xml:"comment"`
	Prv_txn     string	 `xml:"prv_txn"`
	Sum         string	 `xml:"sum"`
}

func NewErrorResponse(code int, comment string) Response {
	return Response{Osmp_txn_id: "0", Result: code, Comment: "Ошибка проведения транзации!", Prv_txn: "0", Sum: "0"}
}