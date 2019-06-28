package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var (
	client *mongo.Client
	err    error

	Collection [12] *mongo.Collection
)

const ProjectName = "houserent"
const UserIdentify = "useridentify"
const UserInfo = "userinfo"
const HouseInfo = "houseinfo"
const HouseListInfo = "houselistinfo"
const ExpireHouse = "expirehouse"
const Wallet = "wallet"
const Discount = "discount"
const Order = "order"
const HouseCheckPool="housecheckpool"
const Checker="checker"
const Houseids="houseids"
const Test = "test"

const (
	USERIDENTIFY  = 0
	USERINFO      = 1
	HOUSEINFO     = 2
	HOUSELISTINFO = 3
	EXPIREHOUSE   = 4
	WALLET        = 5
	DISCOUNT      = 6
	ORDER         = 7
	HOUSECHECKPOOL=8
	CHECKER=9
	HOUSEIDS=10
	TEST          = 11
)

func initDB() {
	log.SetOutput(os.Stdout)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	Collection[USERIDENTIFY] = client.Database(ProjectName).Collection(UserIdentify)
	Collection[USERINFO] = client.Database(ProjectName).Collection(UserInfo)
	Collection[HOUSEINFO] = client.Database(ProjectName).Collection(HouseInfo)
	Collection[HOUSELISTINFO] = client.Database(ProjectName).Collection(HouseListInfo)
	//TODO move rentedhouse to expirehouse
	Collection[EXPIREHOUSE] = client.Database(ProjectName).Collection(ExpireHouse)
	Collection[WALLET] = client.Database(ProjectName).Collection(Wallet)
	Collection[DISCOUNT] = client.Database(ProjectName).Collection(Discount)
	Collection[ORDER] = client.Database(ProjectName).Collection(Order)
	Collection[HOUSECHECKPOOL]=client.Database(ProjectName).Collection(HouseCheckPool)
	Collection[CHECKER]=client.Database(ProjectName).Collection(Checker)
	Collection[HOUSEIDS]=client.Database(ProjectName).Collection(Houseids)
	Collection[TEST] = client.Database(ProjectName).Collection(Test)

}
