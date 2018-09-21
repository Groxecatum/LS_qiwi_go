package model

type GetBalanceRequest struct {
	UserId int `json:"user_id"`
}

type GetBalanceResponse struct {
	Balance int `json:"balance"`
}

type UpdateBalanceRequest struct {
	UserId int `json:"user_id"`
}
