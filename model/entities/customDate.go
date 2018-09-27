package entities

import (
	"encoding/xml"
	"strconv"
	"strings"
	"time"
)

const ctLayoutDate = "2006-01-02"

type CustomDate struct {
	time.Time
}

func (ct *CustomDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	ct.Time, err = time.Parse(ctLayoutDate, s)
	return err
}

func (ct CustomDate) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(ct.Format(ctLayoutDate))), nil
}

func (ct *CustomDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var v string
	d.DecodeElement(&v, &start)
	parsed, err := time.Parse(ctLayoutDate, v)
	if err != nil {
		return err
	}
	ct.Time = parsed
	return nil
}

func (ct CustomDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	v := ct.Format(ctLayoutDate)
	return e.EncodeElement(&v, start)
}
