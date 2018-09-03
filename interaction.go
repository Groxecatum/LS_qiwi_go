package golang_commons

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const ClientBySessionServlet = "msrv_GetClientBySession"
const ActorBySessionServlet = "msrv_GetActorBySession"

//Данные клиента по сессии
func ClientFromSession(sessionId, LSUrl, LSPath string) (Client, error) {
	hc := http.Client{}

	body := "<request frontEnd=\"web\" type=\"" + ClientBySessionServlet + "\"><sessionId>" + sessionId + "</sessionId></request>"
	req, err := http.NewRequest("POST", LSUrl+LSPath+ClientBySessionServlet, strings.NewReader(body))
	if err != nil {
		log.Println("Error getting user for session " + sessionId)
	}

	req.Header.Add("Content-Type", "application/xml")

	resp, err := hc.Do(req)
	if err != nil {
		return Client{}, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	log.Println(string(respBody))

	cli := Client{}
	err = xml.Unmarshal(respBody, &cli)

	return cli, err
}

//Данные клиента по сессии
func ActorFromSession(sessionId, LSUrl, LSPath string) (Actor, error) {
	hc := http.Client{}

	body := "<request frontEnd=\"web\" type=\"" + ActorBySessionServlet + "\"><sessionId>" + sessionId + "</sessionId></request>"
	req, err := http.NewRequest("POST", LSUrl+LSPath+ActorBySessionServlet,
		strings.NewReader(body))
	if err != nil {
		log.Println("Error getting user for session " + sessionId)
	}

	req.Header.Add("Content-Type", "application/xml")

	resp, err := hc.Do(req)
	if err != nil {
		return Actor{}, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	log.Println(string(respBody))

	cli := Actor{}
	err = xml.Unmarshal(respBody, &cli)

	return cli, err
}
