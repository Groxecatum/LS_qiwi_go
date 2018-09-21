package model

import (
	"encoding/json"
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

func ParseActorLoginRequest(r *http.Request) (ActorLoginRequest, error) {
	b, err := ParseReqByte(r)
	var req ActorLoginRequest
	if err != nil {
		return req, err
	}

	err = json.Unmarshal(b, &req)
	return req, err
}
