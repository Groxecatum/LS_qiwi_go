package requests

type GetBalanceRequest struct {
	UserId int `json:"user_id"`
}

type GetBalanceResponse struct {
	Balance int `json:"balance"`
}

type UpdateBalanceRequest struct {
	UserId int `json:"user_id"`
}

type CheckBalanceRequest struct {
	CardNum string  `json:"card_num"`
	Balance float64 `json:"balance"`
}

type CheckBalanceResponse struct {
	Result bool `json:"result"`
}
