package main

import (
	"encoding/gob"
	"flag"
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const(
	SecondsInADay=86400
	MaxTokenProximity=2
)
var(
	e = echo.New()

	searcher=riot.Engine{}

	//weiboData=flag.String()
	dictFile=flag.String("dict_file",
		"../../pkg/mod/github.com/go-ego/riot@v0.0.0-20190307162011-3d971d90bc83/data/dict/dictionary.txt", "词典文件")
	stopTokenFile = flag.String("stop_token_file",
		"../../pkg/mod/github.com/go-ego/riot@v0.0.0-20190307162011-3d971d90bc83/data/dict/stop_tokens.txt", "停用词文件")
	staticFolder = flag.String("static_folder", "static", "静态文件目录")

)

func main() {
	initDB()
	gob.Register(HouseScoringCriteria{})
	searcher.Init(types.EngineOpts{
		Using:1,
		GseDict:*dictFile,
		StopTokenFile:*stopTokenFile,
		UseStore:true,
		StoreFolder:"../storage",
		//StoreShards: 8,
		//StoreEngine:"bg",
	})
	defer searcher.Close()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.POST("/testconnection", testconnection)
	e.POST("/testdate", testDate)
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
	e.POST("/joinchecker",joinchecher)
	e.POST("/getcheckerinfo",getcheckerinfo)
	e.POST("/gettocheckhouse",gettocheckhouse)
	e.POST("/putcheckresult",putcheckresult)
	e.Logger.Fatal(e.Start(":1323"))
}
