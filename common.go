package golang_commons

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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
