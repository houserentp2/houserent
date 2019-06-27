package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"time"
)

type PayOrderJ struct {
	UserID     string      `json:"userid"`
	Token      string      `json:"token"`
	HouseID    string      `json:"houseid"`
	HostID     string      `json:"hostid"`
	OrderID    string      `json:"orderid"`
	DiscountID string      `json:"discountid"`
	Pay        PaychannelJ `json:"pay"`
	Time       time.Time   `json:"time"`
}
type PayOrderD struct {
	UserID     string             `json:"userid"`
	Token      string             `json:"token"`
	HouseID    string             `json:"houseid"`
	HostID     string             `json:"hostid"`
	OrderID    string             `json:"orderid"`
	DiscountID string             `json:"discountid"`
	Pay        PaychannelD        `json:"pay"`
	Time       primitive.DateTime `json:"time"`
}
type PaychannelJ struct {
	AliPay    decimal.Decimal `json:"alipay"`
	WechatPay decimal.Decimal `json:"wechatpay"`
	Balance   decimal.Decimal `json:"balance"`
}
type PaychannelD struct {
	AliPay    primitive.Decimal128 `json:"alipay"`
	WechatPay primitive.Decimal128 `json:"wechatpay"`
	Balance   primitive.Decimal128 `json:"balance"`
}
type PaySucc struct {
	OrderID   string `json:"orderid"`
	Payresult string `json:"payresult"`
}
type DiscountdetailJ struct {
	UserID      string          `json:"userid"`
	DiscountID  string          `json:"discountid"`
	Reduce      decimal.Decimal `json:"reduce"`
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Useable     int             `json:"useable"`
}
type DiscountdetailD struct {
	UserID      string               `json:"userid"`
	DiscountID  string               `json:"discountid"`
	Reduce      primitive.Decimal128 `json:"reduce"`
	Type        string               `json:"type"`
	Description string               `json:"description"`
	Useable     int                  `json:"useable"`
}
type WalletStructJ struct {
	UserID       string            `json:"userid"`
	Score        int32             `json:"score"`
	Balance      decimal.Decimal   `json:"balance"`
	DiscountList []DiscountdetailJ `json:"discountlist"`
	PayOrderList []PayOrderJ       `json:"payorderlist"`
}
type WalletStructD struct {
	UserID       string               `json:"userid"`
	Score        int32                `json:"score"`
	Balance      primitive.Decimal128 `json:"balance"`
	DiscountList []DiscountdetailD    `json:"discountlist"`
	PayOrderList []PayOrderD          `json:"payorderlist"`
}

type GetPrePayInfo LoginSucc

func genDiscountID() string {
	str := strconv.Itoa(int(time.Now().Unix()%8999999999) + 1000000000)
	filter := bson.D{{"discountid", str}}
	var result DiscountdetailD
	err = Collection[DISCOUNT].FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		if err != mongo.ErrNoDocuments {
			return ""
		}
	}
	if result.DiscountID == str {
		return genDiscountID()
	}
	return str
}
func genOrderID() string {
	str := strconv.Itoa(int(time.Now().Unix()%8999999999) + 1000000000)
	filter := bson.D{{"orderid", str}}
	var result PayOrderD
	err = Collection[ORDER].FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		if err != mongo.ErrNoDocuments {
			return ""
		}
	}
	if result.OrderID == str {
		return genOrderID()
	}
	return str
}
func getDiscountList(c echo.Context) error {
	requestbody := new(GetPrePayInfo)
	err := c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Wrong Format")
	}
	if !checkToken(requestbody.UserID, requestbody.Token) {
		return c.String(http.StatusOK, "Invalid Token")
	}
	filter := bson.D{{"userid", requestbody.UserID}}
	var rec WalletStructD
	err = Collection[WALLET].FindOne(context.TODO(), filter).Decode(&rec)
	if err != nil {
		fmt.Println(err)
		if err != mongo.ErrNoDocuments {
			return c.String(http.StatusOK, "Failed")
		}
	}
	if rec.UserID == "" {
		zeroD, _ := primitive.ParseDecimal128("0")
		newdisvD, _ := primitive.ParseDecimal128("100")
		newdist := DiscountdetailD{
			requestbody.UserID,
			genDiscountID(),
			newdisvD,
			"-",
			"新用户优惠",
			1}
		newuser := WalletStructD{
			UserID:       requestbody.UserID,
			Score:        0,
			Balance:      zeroD,
			DiscountList: []DiscountdetailD{newdist,},
		}
		_, err := Collection[WALLET].InsertOne(context.TODO(), newuser)
		if err != nil {
			fmt.Println(err)
		}
		_, err = Collection[DISCOUNT].InsertOne(context.TODO(), newdist)
		if err != nil {
			fmt.Println(err)
		}
		resultJ := conv_WSD_J(newuser)
		return c.JSON(http.StatusOK, resultJ.DiscountList)
	}
	resultJ := conv_WSD_J(rec)
	return c.JSON(http.StatusOK, resultJ.DiscountList)
}
func pay(c echo.Context) error {
	//TODO need lock?
	requestbodyJ := new(PayOrderJ)
	err := c.Bind(requestbodyJ)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Wrong Format")
	}
	if !checkToken(requestbodyJ.UserID, requestbodyJ.Token) {
		return c.String(http.StatusOK, "Invalid Token")
	}
	//check house
	var houseD HouseDetailD
	filter := bson.D{{"houseid", requestbodyJ.HouseID}}
	err = Collection[HOUSEINFO].FindOne(context.TODO(), filter).Decode(&houseD)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Failed")
	}
	if houseD.HouseID == "" {
		return c.String(http.StatusOK, "House Not Exist")
	}
	if houseD.UserID != requestbodyJ.HostID {
		return c.String(http.StatusOK, "Wrong HostID")
	}
	//check money
	var walletD WalletStructD
	filteru := bson.D{{"userid", requestbodyJ.UserID}}
	err = Collection[WALLET].FindOne(context.TODO(), filteru).Decode(&walletD)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Failed3")
	}
	if walletD.UserID == "" {
		return c.String(http.StatusOK, "Failed2")
	}
	if conv_priDe_dD(walletD.Balance).LessThan(requestbodyJ.Pay.Balance) {
		return c.String(http.StatusOK, "Lacking Balance")
	}
	var useable = false
	var loc = 0
	for index, disc := range walletD.DiscountList { //check discount
		if disc.DiscountID == requestbodyJ.DiscountID {
			loc = index
			useable = true
			if disc.Useable == 0 {
				return c.String(http.StatusOK, "Discount Unuseable")
			}
		}
	}
	if !(requestbodyJ.DiscountID == "" || useable) {
		return c.String(http.StatusOK, "Discount Unuseable")
	}
	dismoney:=conv_priDe_dD(walletD.DiscountList[loc].Reduce)
	paysum := decimal.Sum(requestbodyJ.Pay.AliPay, requestbodyJ.Pay.WechatPay, requestbodyJ.Pay.Balance,dismoney)
	price := conv_priDe_dD(houseD.Price)
	if paysum.LessThan(price) {
		return c.String(http.StatusOK, "Pay Not Enough")
	}
	//pay
	_, err = Collection[HOUSEINFO].DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
	}
	_, err = Collection[HOUSELISTINFO].DeleteOne(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
	}
	if useable {
		_, err = Collection[DISCOUNT].DeleteOne(context.TODO(), bson.D{{"disountid", requestbodyJ.DiscountID}})
		if err != nil {
			fmt.Println(err)
		}

		walletD.DiscountList[loc].Useable = 0
	}
	requestbodyJ.Token = ""
	requestbodyJ.Time = time.Now()
	requestbodyJ.OrderID = genOrderID()
	requestbodyD := conv_POJ_D(*requestbodyJ)
	walletD.PayOrderList = append(walletD.PayOrderList, requestbodyD)
	walletD.Balance = conv_dD_priDe(conv_priDe_dD(walletD.Balance).Sub(requestbodyJ.Pay.Balance))

	_, err = Collection[WALLET].ReplaceOne(context.TODO(), filteru, walletD)
	if err != nil {
		fmt.Println(err)
	}
	result := PaySucc{requestbodyJ.OrderID, "Success"}
	return c.JSON(http.StatusOK, result)
}
func getmyrented(c echo.Context) error {
	userm := new(GetPrePayInfo)
	err := c.Bind(userm)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Wrong Format")
	}
	if !checkToken(userm.UserID, userm.Token) {
		return c.String(http.StatusOK, "Invalid Token")
	}
	filter := bson.D{{"userid", userm.UserID}}
	var walletD WalletStructD
	err = Collection[WALLET].FindOne(context.TODO(), filter).Decode(&walletD)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "err 01")
	}
	return c.JSON(http.StatusOK, conv_list_POD_J(walletD.PayOrderList))
}
