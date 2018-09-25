package requests

import "encoding/xml"

type CustomRequest struct {
	XMLName  xml.Name `xml:"request"`
	Type     string   `json:"type"       xml:"type,attr"`
	FrontEnd string   `json:"frontEnd"   xml:"frontEnd,attr"`
}

type CustomResponse struct {
	XMLName     xml.Name `xml:"response"`
	Result      string   `json:"result"            xml:"result,attr"`
	ResultDescr string   `json:"resultDescr"       xml:"resultDescr,attr"`
}
