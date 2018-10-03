package golang_commons

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	TYPE_XML  = "xml"
	TYPE_JSON = "json"
)

func GetFormatByRequest(r *http.Request) string {
	format := TYPE_XML
	if r.Header.Get("Content-Type") == "application/json" {
		format = TYPE_JSON
	}

	return format
}

func ParseReqByte(r *http.Request) ([]byte, error) {
	b, err := ioutil.ReadAll(r.Body)
	return b, err
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func Invert(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func GetFullReference(date time.Time, ref string, actorId int) string {
	return date.Format("yyMMddHHmmss") + ref + strconv.Itoa(actorId)
}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}
