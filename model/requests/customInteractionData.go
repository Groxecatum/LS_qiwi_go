package requests

import (
	"encoding/json"
	"encoding/xml"
)

type CustomRequest struct {
	XMLName  xml.Name `json:",omitempty" xml:"request"`
	Type     string   `json:"type"       xml:"type,attr"`
	FrontEnd string   `json:"frontEnd"   xml:"frontEnd,attr"`
}

type CustomResponse struct {
	XMLName     xml.Name `json:",omitempty"   xml:"response"`
	Result      int      `json:"result"       xml:"result,attr"`
	ResultDescr string   `json:"resultDescr"  xml:"resultDescr,attr"`
}

func CustomRequestFromBytes(b []byte, format string) (CustomRequest, error) {
	req := CustomRequest{}
	switch format {
	case "json":
		err := json.Unmarshal(b, &req)
		return req, err
	default:
		err := xml.Unmarshal(b, &req)
		return req, err
	}
}
