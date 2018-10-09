package testing

import "git.kopilka.kz/BACKEND/golang_commons/model/entities"

func TestActorStub(login string) entities.Actor {
	switch login {
	case "CCFS01":
		return entities.Actor{Id: 13, MerchantId: 5, Title: "Kopilka"}
	case "sp_taxi5353_01":
		return entities.Actor{Id: 166, MerchantId: 24, Title: "Kopilka"}
	case "sp_ramstore_100":
		return entities.Actor{Id: 180, MerchantId: 25, Title: "Kopilka"}
	case "sp_f_bsb_01":
		return entities.Actor{Id: 77, MerchantId: 2, Title: "Kopilka"}
	default:
		return entities.Actor{Id: 10, MerchantId: 1, Title: "Kopilka"}
	}
}
