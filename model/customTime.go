package model

import (
	"strconv"
	"strings"
	"time"
)

const ctLayout = "2006-01-02 15:04:05.999999"

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	ct.Time, err = time.Parse(ctLayout, s)
	return err
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(ct.Format(ctLayout))), nil
}
