package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var e = echo.New()

func main() {
	initDB()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.GET("/testconnection", testconnection)
	e.GET("/testdate", testDate)
	e.POST("/testjson", testjson)
	e.POST("/register", register)
	e.POST("/login", login)
	e.POST("/resetpassword", resetPassword)
	e.POST("/userinfo", updateUserInfo)
	e.POST("/getuserinfo", getUserinfo)
	e.POST("/logout", logout)
	e.POST("/puthouse", puthouse)
	e.POST("/gethouse", gethouse)
	e.POST("/gethouselist", gethouselist)
	e.POST("/gethouselist/:queryparam", gethouselist)
	e.POST("/getmyputs", getMyPuts)
	e.POST("/getmyrented", getmyrented)
	e.POST("/getdiscountlist", getDiscountList)
	e.POST("/pay", pay)
	e.Logger.Fatal(e.Start(":1323"))
}
