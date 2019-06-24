package main

import (
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
)

type PayOrder struct {
	UserID string `json:"user_id"`
	Token string `json:"token"`
	HouseID string `json:"house_id"`
	DiscountID string `json:"discount_id"`
	Pay decimal.Decimal `json:"pay"`
}
type PaySucc struct {
	OrderID string `json:"order_id"`
	Payresult string `json:"payresult"`
}
func pay(c echo.Context)error{
	return nil
}
