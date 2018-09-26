package requests

import (
	"encoding/json"
	"git.kopilka.kz/BACKEND/golang_commons"
	"git.kopilka.kz/BACKEND/golang_commons/errors"
	"git.kopilka.kz/BACKEND/golang_commons/model/entities"
	"net/http"
)

func CheckAuth(sessionId, login, psw, authCacheHost, authCachePath string) (bool, error) {
	loginReq := ActorCheckSessionRequest{SessionId: sessionId}

	code, resp, err := golang_commons.SendObjectJSON(loginReq, authCacheHost, authCachePath)
	if err != nil {
		return false, err
	}

	if code != http.StatusOK {
		return false, errors.AuthErr
	}

	resposeStruct := CustomBoolResponse{}
	err = json.Unmarshal(resp, &resposeStruct)
	if err != nil {
		return false, err
	}

	return resposeStruct.Result, nil
}

func GetAuthorizedActor(sessionId, login, psw, authCacheHost, authCachePath string) (entities.Actor, error) {
	loginReq := ActorCheckSessionRequest{SessionId: sessionId}
	responseStruct := entities.Actor{}

	code, resp, err := golang_commons.SendObjectJSON(loginReq, authCacheHost, authCachePath)
	if err != nil {
		return responseStruct, err
	}

	if code != http.StatusOK {
		return responseStruct, errors.AuthErr
	}

	err = json.Unmarshal(resp, &responseStruct)
	if err != nil {
		return responseStruct, err
	}

	return responseStruct, nil
}

func CheckBalance(cardNum string, bonuses int, balanceCacheHost, balanceCachePath string) (bool, error) {
	balanceReq := ActorCheckSessionRequest{SessionId: ""} //BalanceCheckRequest{CardNum: cardNum, Bonuses: bonuses}

	code, resp, err := golang_commons.SendObjectJSON(balanceReq, balanceCacheHost, balanceCachePath)
	if err != nil {
		return false, err
	}

	if code != http.StatusOK {
		return false, errors.BalanceErr
	}

	resposeStruct := CustomBoolResponse{}
	err = json.Unmarshal(resp, &resposeStruct)
	if err != nil {
		return false, err
	}

	return resposeStruct.Result, nil

}
