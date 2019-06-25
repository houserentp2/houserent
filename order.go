package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"time"
)

type PayOrder struct {
	UserID string `json:"user_id"`
	Token string `json:"token"`
	HouseID string `json:"house_id"`
	HostID string `json:"host_id"`
	OrderID string `json:"order_id"`
	DiscountID string `json:"discount_id"`
	Pay Paychannel `json:"pay"`
	Time primitive.DateTime `json:"time"`
}
type Paychannel struct {
	AliPay decimal.Decimal `json:"ali_pay"`
	WechatPay decimal.Decimal `json:"wechat_pay"`
	Balance decimal.Decimal `json:"balance"`
}
type PaySucc struct {
	OrderID string `json:"order_id"`
	Payresult string `json:"payresult"`
}
type Discountdetail struct{
	UserID string `json:"user_id"`
	DiscountID string `json:"discount_id"`
	Recduce decimal.Decimal `json:"reduce"`
	Type string `json:"type"`
	Description string `json:"description"`
	Useable int `json:"useable"`
}
type WalletStruct struct {
	UserID string `json:"user_id"`
	Score int32 `json:"score"`
	Balance decimal.Decimal `json:"balance"`
	DiscountList []Discountdetail `json:"discount_list"`
	PayOrderList []PayOrder `json:"pay_order_list"`

}

type GetPrePayInfo LoginSucc
func genDiscountID()string{
	str:= strconv.Itoa(int(time.Now().Unix()%8999999999)+1000000000)
	filter:=bson.D{{"discount_id",str}}
	var result  Discountdetail
	err=Collection[DISCOUNT].FindOne(context.TODO(),filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if result.DiscountID==str{
		return genDiscountID()
	}
	return str
}
func genOrderID()string{
	str:= strconv.Itoa(int(time.Now().Unix()%8999999999)+1000000000)
	filter:=bson.D{{"order_id",str}}
	var result  PayOrder
	err=Collection[ORDER].FindOne(context.TODO(),filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if result.OrderID==str{
		return genOrderID()
	}
	return str
}
func getDiscountList(c echo.Context)error{
	requestbody:=new(GetPrePayInfo)
	err:=c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Wrong Format")
	}
	if !checkToken(requestbody.UserID,requestbody.Token) {
		return c.String(http.StatusOK,"Invalid Token")
	}
	filter:=bson.D{{"user_id",requestbody.UserID}}
	var rec WalletStruct
	err=Collection[WALLET].FindOne(context.TODO(),filter).Decode(rec)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Failed")
	}
	if rec.UserID==""{
		zero,_:=decimal.NewFromString("0")
		newdisv,_:=decimal.NewFromString("100")
		newdist:=Discountdetail{requestbody.UserID, genDiscountID(), newdisv, "-", "新用户优惠",1}
		newuser:=WalletStruct{
			UserID:       requestbody.UserID,
			Score:        0,
			Balance:      zero,
			DiscountList: []Discountdetail{newdist,},
		}
		_,err:=Collection[WALLET].InsertOne(context.TODO(),newuser)
		if err != nil {
			fmt.Println(err)
		}
		_,err=Collection[DISCOUNT].InsertOne(context.TODO(),newdist)
		if err != nil {
			fmt.Println(err)
		}

		return c.JSON(http.StatusOK,newuser.DiscountList)
	}
	return c.JSON(http.StatusOK,rec.DiscountList)
}
func pay(c echo.Context)error{
	//TODO need lock?
	requestbody:=new(PayOrder)
	err:=c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Wrong Format")
	}
	if !checkToken(requestbody.UserID,requestbody.Token) {
		return c.String(http.StatusOK,"Invalid Token")
	}
	//check house
	var house HouseDetail
	filter:=bson.D{{"house_id",requestbody.HouseID}}
	err=Collection[HOUSEINFO].FindOne(context.TODO(),filter).Decode(house)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Failed")
	}
	if house.HouseID==""{
		return c.String(http.StatusOK,"House Not Exist")
	}
	if house.UserID!=requestbody.HostID{
		return c.String(http.StatusOK,"Wrong HostID")
	}
	//check money
	var wallet WalletStruct
	filteru:=bson.D{{"user_id",requestbody.UserID}}
	err=Collection[WALLET].FindOne(context.TODO(),filteru).Decode(wallet)
	if wallet.UserID==""{
		return c.String(http.StatusOK,"Failed2")
	}
	if wallet.Balance.LessThan(requestbody.Pay.Balance){
		return c.String(http.StatusOK,"Lacking Balance")
	}
	var useable=false
	var loc=0
	for index,disc:=range wallet.DiscountList{ //check discount
		if disc.DiscountID==requestbody.DiscountID{
			loc=index
			useable=true
			if disc.Useable==0{
				return c.String(http.StatusOK,"Discount Unuseable")
			}
		}
	}
	if !(requestbody.DiscountID=="" || useable){
		return c.String(http.StatusOK,"Discount Unuseable")
	}
	paysum:=decimal.Sum(requestbody.Pay.AliPay,requestbody.Pay.WechatPay,requestbody.Pay.Balance)
	price,_:=decimal.NewFromString(house.Price)
	if paysum.LessThan(price){
		return c.String(http.StatusOK,"Pay Not Enough")
	}
	//pay
	_,err=Collection[HOUSEINFO].DeleteOne(context.TODO(),filter)
	if err != nil {
		fmt.Println(err)
	}
	_,err=Collection[HOUSELISTINFO].DeleteOne(context.TODO(),filter)
	if err != nil {
		fmt.Println(err)
	}
	if useable {
		_, err = Collection[DISCOUNT].DeleteOne(context.TODO(), bson.D{{"disount_id", requestbody.DiscountID}})
		if err != nil {
			fmt.Println(err)
		}

		wallet.DiscountList[loc].Useable=0
	}
	requestbody.Token=""
	requestbody.Time=primitive.DateTime(time.Now().Unix())
	requestbody.OrderID=genOrderID()
	wallet.PayOrderList=append(wallet.PayOrderList, *requestbody)
	wallet.Balance=wallet.Balance.Sub(requestbody.Pay.Balance)

	_,err=Collection[WALLET].ReplaceOne(context.TODO(),filteru,wallet)
	if err != nil {
		fmt.Println(err)
	}
	result:=PaySucc{requestbody.OrderID,"Success"}
	return c.JSON(http.StatusOK,result)
}
