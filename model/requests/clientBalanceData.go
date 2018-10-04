package requests

type GetBalanceRequest struct {
	CardNum string `json:"card_num"`
}

type GetBalanceResponse struct {
	Balance float64 `json:"balance"`
}

type CheckBalanceRequest struct {
	CardNum string  `json:"card_num"`
	Balance float64 `json:"balance"`
}

type UpdateCardBalance struct {
	CardNum string `json:"card_num"`
}
