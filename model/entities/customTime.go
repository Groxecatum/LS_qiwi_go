package entities

import (
	"encoding/xml"
	"strconv"
	"strings"
	"time"
)

const ctLayoutTime = "2006-01-02 15:04:05.999999"

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	ct.Time, err = time.Parse(ctLayoutTime, s)
	return err
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(ct.Format(ctLayoutTime))), nil
}

func (ct *CustomTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	var v string
	d.DecodeElement(&v, &start)
	parsed, err := time.Parse(ctLayoutTime, v)
	if err != nil {
		return err
	}
	ct.Time = parsed
	return nil
}

func (ct CustomTime) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	v := ct.Format(ctLayoutTime)
	return e.EncodeElement(&v, start)
}
