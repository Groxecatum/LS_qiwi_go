package model

import "encoding/xml"

type VerifyData struct {
	XMLName     xml.Name `xml:"request"`
	FrontEnd    string   `xml:"frontEnd,attr"`
	Type        string   `xml:"type,attr"`
	Account     string   `xml:"cellPhoneNum"`
}

type VerifyDataResponse struct {
	XMLName       xml.Name `xml:"response"`
	Result        int      `xml:"result,attr"`
	Description   string   `xml:"resultDescr,attr"`
}