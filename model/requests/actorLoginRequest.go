package requests

import (
	"encoding/json"
	"git.kopilka.kz/BACKEND/golang_commons"
	"net/http"
)

/*  <request type=”mrct_OpenSalesPointSession”>
      <login>[логин]</login>
	  <psw>[пароль]</psw>
      <deviceId>[идентификатор устройства]</deviceId>
      <pushtoken>[пуш токен]</pushtoken>
    </request>
 Ответ:
    <response result=”[код результата. 0 - Ok]” resultDescr=”[описание результата]”>
      <sessionId>[открытая сессия актора]</sessionId>
*   </response>*/

type ActorLoginRequest struct {
	Login     string `json:"login"`
	Password  string `json:"psw"`
	DeviceId  string `json:"deviceId"`
	PushToken string `json:"pushtoken"`
}

/*  <request type=””>
      <sessionId>[Идентификатор сессии точки продаж]</sessionId>
    </request>
 Ответ:
    <response result=”[код результата. 0 - Ok]” resultDescr=”[описание результата]”>
      <sessionId>[открытая сессия актора]</sessionId>
*   </response>*/

type ActorCheckSessionRequest struct {
	SessionId string `json:"sessionId"`
	Login     string `json:"login"`
	Password  string `json:"psw"`
}

type ActorBySessionRequest struct {
	Session string `json:"sessionId" xml:"sessionId"`
}

func ParseActorLoginRequest(r *http.Request) (ActorLoginRequest, error) {
	b, err := golang_commons.ParseReqByte(r)
	var req ActorLoginRequest
	if err != nil {
		return req, err
	}

	err = json.Unmarshal(b, &req)
	return req, err
}

func ParseActorCheckSessionRequest(r *http.Request) (ActorCheckSessionRequest, error) {
	b, err := golang_commons.ParseReqByte(r)
	var req ActorCheckSessionRequest
	if err != nil {
		return req, err
	}

	err = json.Unmarshal(b, &req)
	return req, err
}

func ParseActorBySessionRequest(r *http.Request) (ActorBySessionRequest, error) {
	b, err := golang_commons.ParseReqByte(r)
	var req ActorBySessionRequest
	if err != nil {
		return req, err
	}

	err = json.Unmarshal(b, &req)
	return req, err
}
