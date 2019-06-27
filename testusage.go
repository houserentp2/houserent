package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

//test db and json's usage

func testconnection(c echo.Context) error {
	return c.String(http.StatusOK, "Connect Success")
}

type TestDate struct {
	Name  string               `json:"name"`
	Age   int                  `json:"age"`
	Date  primitive.DateTime   `json:"date"`
	Money primitive.Decimal128 `json:"money"`
}
type TestJ struct {
	Name  string          `json:"name"`
	Age   int             `json:"age"`
	Date  time.Time       `json:"date"`
	Money decimal.Decimal `json:"money"`
}

func testDate(c echo.Context) error {
	money, _ := primitive.ParseDecimal128("123.45")

	var test = TestDate{"doorhong", 18, primitive.DateTime(time.Now().UnixNano() / 1e6), money}
	var ret = TestJ{test.Name, test.Age, conv_priDT_tT(test.Date), conv_priDe_dD(test.Money)}
	in, err := Collection[TEST].InsertOne(context.TODO(), test)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(in.InsertedID)
	return c.JSON(http.StatusOK, ret)
}
func testjson(c echo.Context) error {
	rb := new(TestJ)
	err := c.Bind(rb)
	if err != nil {
		fmt.Println(err)
	}
	var insb = TestDate{rb.Name, rb.Age, conv_tT_priDT(rb.Date), conv_dD_priDe(rb.Money)}
	in, err := Collection[TEST].InsertOne(context.TODO(), insb)
	fmt.Println(in.InsertedID)
	return c.String(http.StatusOK, "ok")
}
