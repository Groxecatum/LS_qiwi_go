package golang_commons

import "encoding/xml"

type LSResponse struct {
	XMLName     xml.Name `xml:"response"`
	Result      int      `xml:"result,attr"`
	ResultDescr int      `xml:"resultDescr,attr"`
}
