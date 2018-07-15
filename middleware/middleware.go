package middleware

import (
	"net/http"
	"time"
	"../config"
	"log"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func Method(m string) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			f(w, r)
		}
	}
}

func Logging() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() { log.Println(r.URL.Path, time.Since(start)) }()
			f(w, r)
		}
	}
}

func checkField(r *http.Request, w http.ResponseWriter, name string) {
	_, ok := r.URL.Query()[name]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func Verify() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			addr := r.Header.Get("X-Forwarded-For");
			log.Println("ReqAddress", "ReqAddress" + addr);
			checkField(r, w,"command");
			checkField(r, w,"txn_id");
			checkField(r, w,"account");
			checkField(r, w,"sum");
			if (config.Config.IP_CHECK) {
				if !((addr == config.Config.QIWI_IP_1) || (addr == config.Config.QIWI_IP_2)) {
					http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				}
			}
			f(w, r)
		}
	}
}