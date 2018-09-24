package golang_commons

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const ClientBySessionServlet = "msrv/msrv_GetClientBySession"
const ActorBySessionServlet = "msrv/msrv_GetActorBySession"

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

	log.Println("Client from session answer: " + string(respBody))

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
		return Actor{}, err
	}

	req.Header.Add("Content-Type", "application/xml")

	resp, err := hc.Do(req)
	if err != nil {
		return Actor{}, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	log.Println("Actor from session answer: " + string(respBody))

	cli := Actor{}
	err = xml.Unmarshal(respBody, &cli)

	return cli, err
}

func SendObjectJSON(obj interface{}, url, path string) ([]byte, error) {
	hc := http.Client{}
	b, err := json.Marshal(obj)
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest("POST", url+path, strings.NewReader(string(b)))
	if err != nil {
		return []byte{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := hc.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)

	log.Println("Sending object answer: " + string(respBody))

	return respBody, err
}