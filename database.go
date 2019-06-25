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
	err error

	Collection[12] *mongo.Collection
)
const ProjectName="houserent"
const UserIdentify ="useridentify"
const UserInfo="userinfo"
const HouseInfo="houseinfo"
const HouseListInfo="houselistinfo"
const ExpireHouse="expirehouse"
const Wallet  = "wallet"
const Discount="discount"
const Order  = "order"

const(
	USERIDENTIFY=0
	USERINFO=1
	HOUSEINFO=2
	HOUSELISTINFO=3
	EXPIREHOUSE=4
	WALLET=5
	DISCOUNT=6
	ORDER=7
)
func initDB(){
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
	Collection[USERINFO]=client.Database(ProjectName).Collection(UserInfo)
	Collection[HOUSEINFO]=client.Database(ProjectName).Collection(HouseInfo)
	Collection[HOUSELISTINFO]=client.Database(ProjectName).Collection(HouseListInfo)
	Collection[EXPIREHOUSE]=client.Database(ProjectName).Collection(ExpireHouse)
	Collection[WALLET]=client.Database(ProjectName).Collection(Wallet)
	Collection[DISCOUNT]=client.Database(ProjectName).Collection(Discount)
	Collection[ORDER]=client.Database(ProjectName).Collection(Order)

}
