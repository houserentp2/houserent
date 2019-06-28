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
	HostID   string `json:"hostid"`
	Icon     string  `json:"icon"`
	HouseID  string               `json:"houseid"`
	Time     primitive.DateTime   `json:"time"`
	Price    primitive.Decimal128 `json:"price"`
	Square   primitive.Decimal128 `json:"square"`
	Shiting  Shiting              `json:"shiting"`
	Title    string               `json:"title"`
	Location Resident             `json:"location"`
	Picture  string               `json:"picture"`
	Others      OtherHouseDetail     `json:"others"`
}
type HouseListItemJ struct {
	UserID   string          `json:"userid"`
	HostID   string `json:"hostid"`
	Icon     string  `json:"icon"`
	HouseID  string          `json:"houseid"`
	Time     time.Time       `json:"time"`
	Price    decimal.Decimal `json:"price"`
	Square   decimal.Decimal `json:"square"`
	Shiting  Shiting         `json:"shiting"`
	Title    string          `json:"title"`
	Location Resident        `json:"location"`
	Picture  string          `json:"picture"`
	Others      OtherHouseDetail     `json:"others"`
}
type Shiting struct {
	Shi  int32 `json:"shi"`
	Ting int32 `json:"ting"`
}
type OtherHouseDetail struct {
	Water    string   `json:"water"`
	Power    string   `json:"power"`
	Net      string   `json:"net"`
	Hot      string   `json:"hot"`
	Aircon   string   `json:"aircon"`
	Bus      string   `json:"bus"`
	Short    int      `json:"short"`
	Long     int      `json:"long"`
	Capacity int32    `json:"capacity"`
	Comments []string `json:"comments"`
	Status   Status `json:"status"`
}
type Status struct {
	Tolive   int    `json:"tolive"`
	Living   int      `json:"living"`
	Lived   int `json:"lived"`
}
type GetHouseIDStruct struct {
	UserID  string `json:"userid"`
	Token   string `json:"token"`
	HouseID string `json:"houseid"`
}
type HouseDetailD struct {
	UserID      string               `json:"userid"`
	Token       string               `json:"token"`
	HostID   string `json:"hostid"`
	Icon     string  `json:"icon"`
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
	HostID   string `json:"hostid"`
	Icon     string  `json:"icon"`
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
type HouseID struct{
	Houseid string `json:"houseid"`
}
func genHouseID() string {
	str := strconv.Itoa(int(time.Now().Unix()%8999999999) + 1000000000)
	filter := bson.D{{"houseid", str}}
	var result HouseDetailD
	err = Collection[HOUSEIDS].FindOne(context.TODO(), filter).Decode(&result)
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
	return HouseListItemD{detail.UserID, detail.HostID,detail.Icon,detail.HouseID, detail.Time, detail.Price, detail.Square, detail.Shiting, detail.Title, detail.Location, detail.Pictures[0],detail.Others}
}
func genMiniDetailJ(detail HouseDetailJ) HouseListItemJ {
	return HouseListItemJ{detail.UserID, detail.HostID,detail.Icon,detail.HouseID, detail.Time, detail.Price, detail.Square, detail.Shiting, detail.Title, detail.Location, detail.Pictures[0],detail.Others}
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
	requestbody.Time = time.Now()
	requestbodyD := conv_HouseDetailJ_D(*requestbody)
	insertRes, err := Collection[HOUSECHECKPOOL].InsertOne(context.TODO(), requestbodyD)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Failed to Create")
	}
	fmt.Println(insertRes.InsertedID)
	houseid:=HouseID{requestbody.HouseID}
	insertRes, err = Collection[HOUSEIDS].InsertOne(context.TODO(), houseid)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Failed to Create2")
	}
	fmt.Println(insertRes.InsertedID)
	return c.String(http.StatusOK, requestbody.HouseID)
}
func updatehouse(c echo.Context) error {
	requestbody := new(HouseDetailJ)
	err := c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Wrong Format")
	}
	if !checkToken(requestbody.UserID, requestbody.Token) {
		return c.String(http.StatusOK, "Invalid Token")
	}
	if requestbody.HouseID == "" {
		return c.String(http.StatusOK, "Lack HouseID")
	}
	filter := bson.D{{"houseid", requestbody.HouseID}}
	idinfo:=new(HouseID)
	err=Collection[HOUSEIDS].FindOne(context.TODO(),filter).Decode(idinfo)
	if err != nil {
		if err==mongo.ErrNoDocuments{
			return c.String(http.StatusOK, "Invalid HouseID")
		}
		fmt.Println(err)
		return c.String(http.StatusOK, "ERR 00")
	}
	data := new(HouseDetailD)
	storeflag:=0
	err = Collection[HOUSEINFO].FindOne(context.TODO(), filter).Decode(data)
	if err != nil {
		if err!=mongo.ErrNoDocuments{
			fmt.Println(err)
			return c.String(http.StatusOK, "ERR 01")
		}else {
			err = Collection[HOUSECHECKPOOL].FindOne(context.TODO(), filter).Decode(data)
			if err != nil {
				if err!=mongo.ErrNoDocuments{
					fmt.Println(err)
					return c.String(http.StatusOK, "ERR 02")
				}else {
					return c.String(http.StatusOK, "ERR 03")
				}
			}else {
				_,err=Collection[HOUSECHECKPOOL].DeleteOne(context.TODO(),filter)
				if err != nil {
					fmt.Println(err)
					return c.String(http.StatusOK, "ERR 04")
				}
			}
		}
	}else{
		_,err=Collection[HOUSEINFO].DeleteOne(context.TODO(),filter)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK, "ERR 05")
		}
		_,err=Collection[HOUSELISTINFO].DeleteOne(context.TODO(),filter)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK, "ERR 06")
		}
	}
	requestbody.Token = ""
	requestbody.Time = time.Now()
	requestdata := conv_HouseDetailJ_D(*requestbody)
	dbres := Collection[HOUSECHECKPOOL].FindOneAndReplace(context.TODO(), filter, requestdata)
	fmt.Println(dbres)
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
	nowtime := conv_tT_priDT(time.Now().AddDate(0, -1, 0))
	//fmt.Println(nowtime)
	//money,_:=primitive.ParseDecimal128("20")
	if queryParam == "" {
		curser, err := Collection[HOUSELISTINFO].Find(context.TODO(), bson.M{"time": bson.M{"$gte": nowtime}}, options.Find().SetLimit(10), options.Find().SetSort(bson.M{"price": -1}))
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
	return c.String(http.StatusOK, "aaa")
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
	var resultD []HouseListItemD
	itemx:=new(HouseDetailD)
	curser,err:=Collection[HOUSECHECKPOOL].Find(context.TODO(),filter, options.Find().SetSort(bson.M{"time": -1}))
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Failed 0")
	}
	for curser.Next(context.Background()) {
		err = curser.Decode(itemx)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK, "Failed 1")
		}
		resultD = append(resultD, genMiniDetailD(*itemx))
	}
	curser, err = Collection[HOUSELISTINFO].Find(context.TODO(), filter, options.Find().SetSort(bson.M{"time": -1}))
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Failed 2")
	}
	item := new(HouseListItemD)
	for curser.Next(context.Background()) {
		err = curser.Decode(item)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK, "Failed 3")
		}
		resultD = append(resultD, *item)
	}
	curser, err = Collection[EXPIREHOUSE].Find(context.TODO(), filter, options.Find().SetSort(bson.M{"time": -1}))
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Failed 4")
	}
	for curser.Next(context.Background()) {
		err = curser.Decode(item)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK, "Failed 5")
		}
		resultD = append(resultD, *item)
	}
	var resultJ []HouseListItemJ
	for _, item := range resultD {
		resultJ = append(resultJ, conv_HouseListItemD_J(item))
	}
	defer curser.Close(context.Background())
	return c.JSON(http.StatusOK, resultJ)
}
type CommentStruct struct{
	UserID string `json:"userid"`
	Token  string `json:"token"`
	HouseID string `json:"houseid"`
	Comment string `json:"comment"`
}
func putcomment(c echo.Context)error{
	requestbody:=new(CommentStruct)
	err:=c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Wrong Format")
	}
	if !checkToken(requestbody.UserID,requestbody.Token){
		return c.String(http.StatusOK, "Invalid Token")
	}
	filter:=bson.D{{"houseid",requestbody.HouseID}}
	doc:=new(HouseDetailD)
	err=Collection[HOUSEINFO].FindOne(context.TODO(),filter).Decode(doc)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "ERR 01")
	}
	doc.Others.Comments=append(doc.Others.Comments,requestbody.Comment)
	inres:=Collection[HOUSEINFO].FindOneAndReplace(context.TODO(),filter,doc)
	fmt.Println(inres)
	return  c.String(http.StatusOK, "Put Comment Success")
}
func checkinvite(c CheckerStruct)bool{
	return true
}
type CheckerStruct struct {
	UserID string `json:"userid"`
	Token  string `json:"token"`
	Invite string `json:"invite"`
}
func joinchecher(c echo.Context)error{
	requestbody:=new(CheckerStruct)
	err:=c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Wrong Format")
	}
	if !checkToken(requestbody.UserID,requestbody.Token) {
		return c.String(http.StatusOK, "Invalid Token")
	}
	filter:=bson.D{{"userid",requestbody.UserID}}
	doc:=new(CheckerStruct)
	err=Collection[CHECKER].FindOne(context.TODO(),filter).Decode(doc)
	if err != nil {
		if err!=mongo.ErrNoDocuments {
			fmt.Println(err)
			return c.String(http.StatusOK, "ERR 01")
		}
	}
	if doc.UserID!=""{
		return c.String(http.StatusOK, "Permission Existed")
	}
	if checkinvite(*requestbody){
		_,err=Collection[CHECKER].InsertOne(context.TODO(),requestbody)
		if err != nil {
			fmt.Println(err)
		}
		return c.String(http.StatusOK, "Permission Get")
	}else {
		return c.String(http.StatusOK, "Invalid Invite")
	}
}
func getcheckerinfo(c echo.Context)error{
	requestbody:=new(LoginSucc)
	err:=c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Wrong Format")
	}
	if !checkToken(requestbody.UserID,requestbody.Token) {
		return c.String(http.StatusOK, "Invalid Token")
	}
	filter:=bson.D{{"userid",requestbody.UserID}}
	doc:=new(CheckerStruct)
	err=Collection[CHECKER].FindOne(context.TODO(),filter).Decode(doc)
	if err != nil {
		if err!=mongo.ErrNoDocuments {
			fmt.Println(err)
			return c.String(http.StatusOK, "ERR 01")
		}
		return c.String(http.StatusOK, "NOPERM")
	}
	if doc.UserID==requestbody.UserID{
		return c.String(http.StatusOK, "Permission Existed")
	}
	return c.String(http.StatusOK, "ERR 02")
}
func gettocheckhouse(c echo.Context)error{
	requestbody:=new(LoginSucc)
	err:=c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Wrong Format")
	}
	if !checkToken(requestbody.UserID,requestbody.Token) {
		return c.String(http.StatusOK, "Invalid Token")
	}
	filter:=bson.D{{"userid",requestbody.UserID}}
	doc:=new(CheckerStruct)
	err=Collection[CHECKER].FindOne(context.TODO(),filter).Decode(doc)
	if err != nil {
		if err!=mongo.ErrNoDocuments {
			fmt.Println(err)
			return c.String(http.StatusOK, "ERR 01")
		}
		return c.String(http.StatusOK, "NOPERM")
	}
	house:=new(HouseDetailD)
	err=Collection[HOUSECHECKPOOL].FindOne(context.TODO(),bson.M{"price":bson.M{"$gte":"0"}}).Decode(house)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "ERR 01")
	}
	return c.JSON(http.StatusOK,conv_HouseDetailD_J(*house))
}
type Checkresult struct {
	UserID string `json:"userid"`
	Token  string `json:"token"`
	HouseID string `json:"houseid"`
	Result int `json:"result"`
}
func putcheckresult(c echo.Context)error{
	requestbody:=new(Checkresult)
	err:=c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK, "Wrong Format")
	}
	if !checkToken(requestbody.UserID,requestbody.Token) {
		return c.String(http.StatusOK, "Invalid Token")
	}
	filter:=bson.D{{"houseid",requestbody.HouseID}}
	if requestbody.Result==1{
		house:=new(HouseDetailD)
		err=Collection[HOUSECHECKPOOL].FindOne(context.TODO(),filter).Decode(house)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK, "ERR 01")
		}
		_,err=Collection[HOUSEINFO].InsertOne(context.TODO(),house)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK, "ERR 02")
		}
		_,err=Collection[HOUSELISTINFO].InsertOne(context.TODO(),genMiniDetailD(*house))
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK, "ERR 03")
		}
		return c.String(http.StatusOK, "Success")
	}else if requestbody.Result==0{
		_,err=Collection[HOUSECHECKPOOL].DeleteOne(context.TODO(),filter)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK, "ERR 04")
		}
	}
	return c.String(http.StatusOK, "ERR 05")
}