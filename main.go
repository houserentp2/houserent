package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const(
	SecondsInADay=86400
	MaxTokenProximity=2
)
var(
	e = echo.New()
	/*
	searcher=riot.Engine{}
	wbs=map[string]HouseDetailD{}
	//weiboData=flag.String()
	dictFile=flag.String("dict_file",
		"X:/Users/huang/go/pkg/mod/github.com/go-ego/riot@v0.0.0-20190307162011-3d971d90bc83/data/dict/dictionary.txt", "词典文件")
	stopTokenFile = flag.String("stop_token_file",
		"X:/Users/huang/go/pkg/mod/github.com/go-ego/riot@v0.0.0-20190307162011-3d971d90bc83/data/dict/stop_tokens.txt", "停用词文件")
	staticFolder = flag.String("static_folder", "static", "静态文件目录")
	*/
)

func main() {
	initDB()
	/*
	searcher.Init(types.EngineOpts{
		Using:1,
		GseDict:*dictFile,
		StopTokenFile:"X:/Users/huang/go/pkg/mod/github.com/go-ego/riot@v0.0.0-20190307162011-3d971d90bc83/data/dict/stop_tokens.txt",
	})
	defer searcher.Close()
	*/
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
	e.POST("/putcomment",putcomment)
	e.POST("/updatehouse",updatehouse)
	e.POST("/gethouse", gethouse)
	e.POST("/gethouselist", gethouselist)
	e.POST("/gethouselist/:queryparam", gethouselist)
	e.POST("/getmyputs", getMyPuts)
	e.POST("/getmyrented", getmyrented)
	e.POST("/getdiscountlist", getDiscountList)
	e.POST("/pay", pay)
	e.Logger.Fatal(e.Start(":1323"))
}
