package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"time"
)

type PhoneReg struct {
	PhoneNum string `json:"phone_num"`
	Code string `json:"code"`
	Time primitive.DateTime `json:"time"`
	Password string `json:"password"`
}
type UserIdentifyStruct struct {
	PhoneNum string `json:"phone_num"`
	Password string `json:"password"`
	Time primitive.DateTime `json:"time"`
	UserID string `json:"user_id"`
	Token string `json:"token"`
}
type RegSucc struct {
	UserID string `json:"user_id"`
	Token string `json:"token"`
}
type UserDetInfo struct {
	UserID string `json:"user_id"`
	Token string `json:"token"`
	Nickname string `json:"nickname"`
	ID string `json:"id"`
	Resident Resident `json:"resident"`
}
type Resident struct {
	Province string `json:"province"`
	City string `json:"city"`
	Zone string `json:"zone"`
	Path string `json:"path"`
}
type PhoneLogin struct {
	PhoneNum string `json:"phone_num"`
	UserID string `json:"user_id"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}
type LoginSucc struct {
	UserID string `json:"user_id"`
	Token string `json:"token"`
}
type Userminid LoginSucc

func register(c echo.Context)error{ // 手机号注册
	// 获取请求体检查格式
	user:=new(PhoneReg)
	err:=c.Bind(user)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Wrong Format")
	}
	// 检查手机号是否已注册
	filter:=bson.D{{"phone_num",user.PhoneNum}}
	var result  UserIdentifyStruct
	err=Collection[USERIDENTIFY].FindOne(context.TODO(),filter).Decode(&result)
	if err != nil { // 查询数据库问题
		fmt.Println(err)
		return c.String(http.StatusOK,"Reg Failed 0")
	}
	if result.PhoneNum !=""{ // 手机号已注册
		fmt.Println(result.PhoneNum)
		return c.String(http.StatusOK,"Phone Number Existed")
	}
	//生成UserID
	result.UserID=genUserID()
	if result.UserID==""{ // 生成ID失败
		return c.String(http.StatusOK,"Reg Failed 1")
	}
	result.PhoneNum=user.PhoneNum
	result.Password=user.Password
	result.Token=genToken(result.UserID)
	var regres=RegSucc{result.UserID,result.Token} // 返回体
	// 注册结果插入数据库
	insert,err:=Collection[USERIDENTIFY].InsertOne(context.TODO(),result)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Reg Failed 2")
	}
	fmt.Println(insert.InsertedID)
	// 注册时自动生成空的UserDetailInfo
	userinfo:=new(UserDetInfo)
	userinfo.UserID=result.UserID
	insertuserinfo,err:=Collection[USERINFO].InsertOne(context.TODO(),userinfo)
	if err != nil { // 生成失败时清理注册
		fmt.Print(err)
		del,err:=Collection[USERIDENTIFY].DeleteOne(context.TODO(),filter)
		if err!=nil {
			fmt.Println(err)
		}
		fmt.Println(del)
		return c.String(http.StatusOK,"Reg Failed 3")
	}
	fmt.Println(insertuserinfo)

	return c.JSON(http.StatusOK,regres)
}
func genUserID() string{ // 生成UserID
	str:=strconv.Itoa(int(time.Now().Unix()%89999999+10000000))
	// 检查重复，重复重新生成
	filter:=bson.D{{"user_id",str}}
	var result  UserIdentifyStruct
	err=Collection[USERIDENTIFY].FindOne(context.TODO(),filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if result.UserID==str{
		return genUserID()
	}
	return str
}
func genToken(str string)string{ // 生成Token
	//TODO 重写Token生成算法
	ctx:=md5.New()
	ctx.Write([]byte(str))
	return hex.EncodeToString(ctx.Sum(nil))
}
func login(c echo.Context)error{ // 登录
	// 检查格式
	requestbody:=new(PhoneLogin)
	err:=c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Wrong Format")
	}
	// 检查登录形式，暂支持手机号和UserID登录
	//TODO？ 手机号和验证码直接登录
	//TODO 重构
	var filter bson.D
	var flag int32
	if requestbody.UserID!=""{
		filter=bson.D{{"user_id",requestbody.UserID}}
		flag=0
	}else if requestbody.PhoneNum!=""{
		filter=bson.D{{"phone_num",requestbody.PhoneNum}}
		flag=1
	}else{
		return c.String(http.StatusOK,"Lack Info")
	}
	// 查询数据库
	var result  UserIdentifyStruct
	err=Collection[0].FindOne(context.TODO(),filter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Login Failed 0")
	}
	if result.UserID==""{
		return c.String(http.StatusOK,"Account Not Exist")
	}
	// 验证
	if flag==0{
		if requestbody.UserID==result.UserID && requestbody.Password==result.Password{
			logsucc:=LoginSucc{result.UserID,genToken(result.UserID)}
			updateRes,err:=Collection[USERIDENTIFY].UpdateOne(context.TODO(),filter, bson.M{"$set": bson.M{"token":logsucc.Token}})
			if err!=nil{
				fmt.Println(err)
				return c.String(http.StatusOK,"Login Failed 1")
			}
			fmt.Println(updateRes)
			return c.JSON(http.StatusOK,logsucc)
		}
		return c.String(http.StatusOK,"Wrong Password")
	}else {
		if requestbody.PhoneNum==result.PhoneNum && requestbody.Password==result.Password{
			logsucc:=LoginSucc{result.UserID,genToken(result.UserID)}
			updateRes,err:=Collection[USERIDENTIFY].UpdateOne(context.TODO(),filter, bson.M{"$set": bson.M{"token":logsucc.Token}})
			if err!=nil{
				fmt.Println(err)
				return c.String(http.StatusOK,"Login Failed 2")
			}
			fmt.Println(updateRes)
			return c.JSON(http.StatusOK,logsucc)
		}
		return c.String(http.StatusOK,"Wrong Password")
	}
}
func updateUserInfo(c echo.Context)error{ // 更新用户信息
	requestbody:=new(UserDetInfo)
	olddata:=new(UserDetInfo)
	err:=c.Bind(requestbody)// 检查格式
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Wrong Format")
	}
	//check token
	if !checkToken(requestbody.UserID,requestbody.Token){
		return c.String(http.StatusOK,"Invalid Token")
	}
	//update
	filter:=bson.D{{"user_id",requestbody.UserID}}
	err=Collection[USERINFO].FindOneAndReplace(context.TODO(),filter,requestbody).Decode(olddata)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Failed to Update")
	}
	return c.String(http.StatusOK,"Update Success")
}
func getUserinfo(c echo.Context)error{
	requestbody:=new(Userminid)
	err:=c.Bind(requestbody)
	if !checkToken(requestbody.UserID,requestbody.Token){
		return c.String(http.StatusOK,"Invalid Token")
	}
	if err!=nil{
		fmt.Println(err)
		return c.String(http.StatusOK,"Wrong Format")
	}
	filter:=bson.D{{"user_id",requestbody.UserID}}
	userinfo:=new(UserDetInfo)
	err=Collection[USERINFO].FindOne(context.TODO(),filter).Decode(userinfo)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Failed to Get")
	}
	return c.JSON(http.StatusOK,userinfo)
}
func logout(c echo.Context)error{
	requestbody:=new(Userminid)
	err:=c.Bind(requestbody)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Wrong Format")
	}
	filter:=bson.D{{"user_id",requestbody.UserID}}
	userindetify:=new(UserIdentifyStruct)
	err=Collection[USERIDENTIFY].FindOne(context.TODO(),filter).Decode(userindetify)
	if err!=nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Invalid UserID")
	}
	if(userindetify.Token==requestbody.Token){
		Updateres,err:=Collection[USERIDENTIFY].UpdateOne(context.TODO(),filter,bson.M{"$set":bson.M{"token":""}})
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusOK,"Invalid UserID")
		}
		fmt.Println(Updateres)
	}
	return c.String(http.StatusOK,"Logout Success")
}
func checkToken(id,token string)bool{
	//当前通过检查token一致性确认有效性
	// TODO 从token解码有效性
	filter:=bson.D{{"user_id",id}}
	userindetify:=new(UserIdentifyStruct)
	err:=Collection[USERIDENTIFY].FindOne(context.TODO(),filter).Decode(userindetify)
	if err!=nil {
		fmt.Println(err)
		return false
	}
	if token!=userindetify.Token{
		return false
	}
	return true
}
func resetPassword(c echo.Context)error{
	// 获取请求体检查格式
	user:=new(PhoneReg)
	err:=c.Bind(user)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Wrong Format")
	}
	// 检查手机号是否已注册
	filter:=bson.D{{"phone_num",user.PhoneNum}}
	var result  UserIdentifyStruct
	err=Collection[USERIDENTIFY].FindOne(context.TODO(),filter).Decode(&result)
	if err != nil { // 查询数据库问题
		fmt.Println(err)
		return c.String(http.StatusOK,"Reg Failed 0")
	}
	if result.PhoneNum ==""{ // 手机号未注册
		fmt.Println(result.PhoneNum)
		return c.String(http.StatusOK,"Phone Number No Existed")
	}
	//TODO check code
	//修改密码
	Updateres,err:=Collection[USERIDENTIFY].UpdateOne(context.TODO(),filter,bson.M{"$set":bson.M{"password":user.Password,"token":""}})
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusOK,"Failed to Update Password")
	}
	fmt.Println(Updateres)
	return c.String(http.StatusOK,"Update Password Success")
}
