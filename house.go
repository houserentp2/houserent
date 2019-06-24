package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"time"
)

type GetHLStruct LoginSucc
type HouseListItem struct {
	HouseID string `json:"house_id"`
	Time primitive.DateTime `json:"time"`
	Price string `json:"price"`
	Square string `json:"square"`
	Shiting Shiting `json:"shiting"`
	Title string `json:"title"`
	Location Resident `json:"location"`
	Picture string `json:"picture"`
}
type Shiting struct {
	Shi int32 `json:"shi"`
	Ting int32 `json:"ting"`
}
type OtherHouseDetail struct {
	Water string `json:"water"`
	Power string `json:"power"`
	Net string `json:"net"`
	Hot string `json:"hot"`
	Aircon string `json:"aircon"`
	Bus string `json:"bus"`
}
type GetHouseIDStruct struct {
	UserID string `json:"user_id"`
	Token string `json:"token"`
	HouseID string `json:"house_id"`
}
type HouseDetail struct{
	UserID string `json:"user_id"`
	Token string `json:"token"`
	HouseID string `json:"house_id"`
	Time primitive.DateTime `json:"time"`
	Price string `json:"price"`
	Square string `json:"square"`
	Shiting Shiting `json:"shiting"`
	Title string `json:"title"`
	Description string `json:"description"`
	Location Resident `json:"location"`
	Pictures []string `json:"picture"`
	Others OtherHouseDetail `json:"others"`
}

func genHouseID()string{
	str:= strconv.Itoa(int(time.Now().Unix()%8999999999)+1000000000)
	filter:=bson.D{{"house_id",str}}
	var result  HouseDetail
	err=Collection[HOUSEINFO].FindOne(context.TODO(),filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if result.UserID==str{
		return genUserID()
	}
	return str
}
func genMiniDetail(detail HouseDetail)HouseListItem{
	return HouseListItem{detail.HouseID, detail.Time,detail.Price,detail.Square,detail.Shiting,detail.Title,detail.Location,detail.Pictures[0]}
}
func puthouse(c echo.Context)error{
	requestbody:=new(HouseDetail)
	err:=c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Wrong Format")
	}
	if !checkToken(requestbody.UserID,requestbody.Token){
		return c.String(http.StatusOK,"Invalid Token")
	}
	requestbody.HouseID=genHouseID()
	requestbody.Token=""
	insertRes,err:=Collection[HOUSEINFO].InsertOne(context.TODO(),requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Failed to Create")
	}
	fmt.Println(insertRes.InsertedID)
	insertRes,err=Collection[HOUSELISTINFO].InsertOne(context.TODO(),genMiniDetail(*requestbody))
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Failed to Create2")
	}
	fmt.Println(insertRes.InsertedID)
	return c.String(http.StatusOK,requestbody.HouseID)
}
func gethouse(c echo.Context)error{
	requestbody:=new(GetHouseIDStruct)
	err:=c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Wrong Format")
	}
	if !checkToken(requestbody.UserID,requestbody.Token){
		return c.String(http.StatusOK,"Invalid Token")
	}
	result:=new(HouseDetail)
	filter:=bson.D{{"house_id",requestbody.HouseID}}
	err=Collection[HOUSEINFO].FindOne(context.TODO(),filter).Decode(result)
	if err != nil || result.HouseID=="" {
		return c.String(http.StatusOK,"Cannot Find")
	}
	return c.JSON(http.StatusOK,result)
}
func gethouselist(c echo.Context)error{
	//TODO FUCK
	requestbody:=new(GetHLStruct)
	err:=c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Wrong Format")
	}
	if !checkToken(requestbody.UserID,requestbody.Token) {
		return c.String(http.StatusOK,"Invalid Token")
	}
	queryParam:=c.Param("queryparam")
	nowtime:=time.Now().Unix()
	if(queryParam==""){
		curser,err:=Collection[HOUSEINFO].Find(context.TODO(),bson.M{"time":bson.M{"gte":nowtime-100000000}},options.Find().SetLimit(10),options.Find().SetSort(bson.M{"time": -1}))
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK,"Failed")
		}
		defer curser.Close(context.Background())
		var result []HouseListItem
		item:=new(HouseListItem)
		for curser.Next(context.Background()){
			err=curser.Decode(item)
			if err != nil {
				fmt.Println(err)
			}
			result= append(result, *item)
		}
		return c.JSON(http.StatusOK,result)
	}
	return nil
}
func getMyPuts(c echo.Context)error{
	requestbody:=new(GetHLStruct)
	err:=c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Wrong Format")
	}
	if !checkToken(requestbody.UserID,requestbody.Token) {
		return c.String(http.StatusOK,"Invalid Token")
	}
	filter:=bson.D{{"user_id",requestbody.UserID}}
	curser,err:=Collection[HOUSELISTINFO].Find(context.TODO(),filter,options.Find().SetSort(bson.M{"time": -1}))
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Failed")
	}
	defer curser.Close(context.Background())
	var result []HouseListItem
	item:=new(HouseListItem)
	for curser.Next(context.Background()){
		err=curser.Decode(item)
		if err != nil {
			fmt.Println(err)
		}
		result= append(result, *item)
	}
	return c.JSON(http.StatusOK,result)
}