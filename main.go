package main

import (
	"net/http"
	"./middleware"
	"./service"
	"./config"
)

func main() {
	config.LoadConfig()
	// Для понятности порядка middleware(он обратный) -
	// представляйте, как будто я одеваете в кучу шкурок апельсин. А запрос их снимает
	http.HandleFunc("/transaction", middleware.Chain(service.NewTransaction, middleware.Verify(), middleware.Logging(), middleware.Method("GET")))

	http.ListenAndServe(":" + config.Config.APP_PORT, nil)
}
