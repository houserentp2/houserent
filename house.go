package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"time"
)

type GetHLStruct LoginSucc
type HouseListItemD struct {
	UserID   string               `json:"userid"`
	HouseID  string               `json:"houseid"`
	Time     primitive.DateTime   `json:"time"`
	Price    primitive.Decimal128 `json:"price"`
	Square   primitive.Decimal128 `json:"square"`
	Shiting  Shiting              `json:"shiting"`
	Title    string               `json:"title"`
	Location Resident             `json:"location"`
	Picture  string               `json:"picture"`
}
type HouseListItemJ struct {
	UserID   string          `json:"userid"`
	HouseID  string          `json:"houseid"`
	Time     time.Time       `json:"time"`
	Price    decimal.Decimal `json:"price"`
	Square   decimal.Decimal `json:"square"`
	Shiting  Shiting         `json:"shiting"`
	Title    string          `json:"title"`
	Location Resident        `json:"location"`
	Picture  string          `json:"picture"`
}
type Shiting struct {
	Shi  int32 `json:"shi"`
	Ting int32 `json:"ting"`
}
type OtherHouseDetail struct {
	Water  string `json:"water"`
	Power  string `json:"power"`
	Net    string `json:"net"`
	Hot    string `json:"hot"`
	Aircon string `json:"aircon"`
	Bus    string `json:"bus"`
}
type GetHouseIDStruct struct {
	UserID  string `json:"userid"`
	Token   string `json:"token"`
	HouseID string `json:"houseid"`
}
type HouseDetailD struct {
	UserID      string               `json:"userid"`
	Token       string               `json:"token"`
	HouseID     string               `json:"houseid"`
	Time        primitive.DateTime   `json:"time"`
	Price       primitive.Decimal128 `json:"price"`
	Square      primitive.Decimal128 `json:"square"`
	Shiting     Shiting              `json:"shiting"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Location    Resident             `json:"location"`
	Pictures    []string             `json:"pictures"`
	Others      OtherHouseDetail     `json:"others"`
}
type HouseDetailJ struct {
	UserID      string           `json:"userid"`
	Token       string           `json:"token"`
	HouseID     string           `json:"houseid"`
	Time        time.Time        `json:"time"`
	Price       decimal.Decimal  `json:"price"`
	Square      decimal.Decimal  `json:"square"`
	Shiting     Shiting          `json:"shiting"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Location    Resident         `json:"location"`
	Pictures    []string         `json:"pictures"`
	Others      OtherHouseDetail `json:"others"`
}

func genHouseID() string {
	str := strconv.Itoa(int(time.Now().Unix()%8999999999) + 1000000000)
	filter := bson.D{{"houseid", str}}
	var result HouseDetailD
	err = Collection[HOUSEINFO].FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		if err != mongo.ErrNoDocuments {
			return ""
		}
	}
	if result.UserID == str {
		return genUserID()
	}
	return str
}
func genMiniDetailD(detail HouseDetailD) HouseListItemD {
	return HouseListItemD{detail.UserID, detail.HouseID, detail.Time, detail.Price, detail.Square, detail.Shiting, detail.Title, detail.Location, detail.Pictures[0]}
}
func genMiniDetailJ(detail HouseDetailJ) HouseListItemJ {
	return HouseListItemJ{detail.UserID, detail.HouseID, detail.Time, detail.Price, detail.Square, detail.Shiting, detail.Title, detail.Location, detail.Pictures[0]}
}
func puthouse(c echo.Context) error {
	requestbody := new(HouseDetailJ)
	err := c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Wrong Format")
	}
	if !checkToken(requestbody.UserID, requestbody.Token) {
		return c.String(http.StatusOK, "Invalid Token")
	}
	requestbody.HouseID = genHouseID()
	requestbody.Token = ""
	requestbodyD := conv_HouseDetailJ_D(*requestbody)
	insertRes, err := Collection[HOUSEINFO].InsertOne(context.TODO(), requestbodyD)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Failed to Create")
	}
	fmt.Println(insertRes.InsertedID)
	insertRes, err = Collection[HOUSELISTINFO].InsertOne(context.TODO(), genMiniDetailD(requestbodyD))
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Failed to Create2")
	}
	fmt.Println(insertRes.InsertedID)
	return c.String(http.StatusOK, requestbody.HouseID)
}
func gethouse(c echo.Context) error {
	requestbody := new(GetHouseIDStruct)
	err := c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Wrong Format")
	}
	if !checkToken(requestbody.UserID, requestbody.Token) {
		return c.String(http.StatusOK, "Invalid Token")
	}
	resultD := new(HouseDetailD)
	filter := bson.D{{"houseid", requestbody.HouseID}}
	err = Collection[HOUSEINFO].FindOne(context.TODO(), filter).Decode(resultD)
	if err != nil || resultD.HouseID == "" {
		return c.String(http.StatusOK, "Cannot Find")
	}
	resultJ := conv_HouseDetailD_J(*resultD)
	return c.JSON(http.StatusOK, resultJ)
}
func gethouselist(c echo.Context) error {
	//TODO FUCK
	requestbody := new(GetHLStruct)
	err := c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Wrong Format")
	}
	if !checkToken(requestbody.UserID, requestbody.Token) {
		return c.String(http.StatusOK, "Invalid Token")
	}
	queryParam := c.Param("queryparam")
	nowtime := conv_tT_priDT(time.Now().AddDate(0,-1,0))
	fmt.Println(nowtime)
	money,_:=primitive.ParseDecimal128("20")
	if queryParam == "" {
		curser, err := Collection[HOUSEINFO].Find(context.TODO(), bson.M{"price": bson.M{"gte": money}}, options.Find().SetLimit(10), options.Find().SetSort(bson.M{"price": -1}))
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK, "Failed")
		}
		defer curser.Close(context.Background())
		var resultD []HouseListItemD
		item := new(HouseListItemD)
		for curser.Next(context.Background()) {
			err = curser.Decode(item)
			if err != nil {
				fmt.Println(err)
			}
			resultD = append(resultD, *item)
		}
		var resultJ []HouseListItemJ
		for _, item := range resultD {
			resultJ = append(resultJ, conv_HouseListItemD_J(item))
		}
		return c.JSON(http.StatusOK, resultJ)
	}
	return c.String(http.StatusOK,"aaa")
}
func getMyPuts(c echo.Context) error {
	requestbody := new(GetHLStruct)
	err := c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Wrong Format")
	}
	if !checkToken(requestbody.UserID, requestbody.Token) {
		return c.String(http.StatusOK, "Invalid Token")
	}
	filter := bson.D{{"userid", requestbody.UserID}}
	curser, err := Collection[HOUSELISTINFO].Find(context.TODO(), filter, options.Find().SetSort(bson.M{"time": -1}))
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Failed")
	}
	defer curser.Close(context.Background())
	var resultD []HouseListItemD
	item := new(HouseListItemD)
	for curser.Next(context.Background()) {
		err = curser.Decode(item)
		if err != nil {
			fmt.Println(err)
		}
		resultD = append(resultD, *item)
	}
	var resultJ []HouseListItemJ
	for _, item := range resultD {
		resultJ = append(resultJ, conv_HouseListItemD_J(item))
	}
	return c.JSON(http.StatusOK, resultJ)
}
