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
	XMLName     xml.Name `xml:"person"`
	Osmp_txn_id string
	Result      int
	Comment     string
	Prv_txn     string
	Sum         string
}