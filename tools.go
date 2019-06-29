package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func conv_priDT_tT(dt primitive.DateTime) time.Time {
	return time.Unix(int64(dt)/1000, 0)
}
func conv_tT_priDT(t time.Time) primitive.DateTime {
	return primitive.DateTime(t.UnixNano() / 1000000)
}
func conv_priDe_dD(de primitive.Decimal128) decimal.Decimal {
	tr, err := decimal.NewFromString(de.String())
	if err != nil {
		fmt.Println(err)
	}
	return tr
}
func conv_dD_priDe(dd decimal.Decimal) primitive.Decimal128 {
	tr, err := primitive.ParseDecimal128(dd.String())
	if err != nil {
		fmt.Println(err)
	}
	return tr
}
func conv_HouseDetailJ_D(hj HouseDetailJ) HouseDetailD {
	return HouseDetailD{
		hj.UserID,
		hj.Token,
		hj.HostID,
		hj.Icon,
		hj.HouseID,
		conv_tT_priDT(hj.Time),
		conv_dD_priDe(hj.Price),
		conv_dD_priDe(hj.Square),
		hj.Shiting,
		hj.Title,
		hj.Description,
		hj.Location,
		hj.Pictures,
		hj.Others,
		0,
	}
}
func conv_HouseDetailD_J(hd HouseDetailD) HouseDetailJ {
	return HouseDetailJ{
		hd.UserID,
		hd.Token,
		hd.HostID,
		hd.Icon,
		hd.HouseID,
		conv_priDT_tT(hd.Time),
		conv_priDe_dD(hd.Price),
		conv_priDe_dD(hd.Square),
		hd.Shiting,
		hd.Title,
		hd.Description,
		hd.Location,
		hd.Pictures,
		hd.Others,
	}
}
func conv_HouseListItemD_J(hld HouseListItemD) HouseListItemJ {
	return HouseListItemJ{
		hld.UserID,
		hld.HostID,
		hld.Icon,
		hld.HouseID,
		conv_priDT_tT(hld.Time),
		conv_priDe_dD(hld.Price),
		conv_priDe_dD(hld.Square),
		hld.Shiting,
		hld.Title,
		hld.Location,
		hld.Picture,
		hld.Others,
	}
}
func conv_PcJ_D(j PaychannelJ) PaychannelD {
	return PaychannelD{
		conv_dD_priDe(j.AliPay),
		conv_dD_priDe(j.WechatPay),
		conv_dD_priDe(j.Balance),
	}
}
func conv_PcD_J(d PaychannelD) PaychannelJ {
	return PaychannelJ{
		conv_priDe_dD(d.AliPay),
		conv_priDe_dD(d.WechatPay),
		conv_priDe_dD(d.Balance),
	}
}
func conv_POJ_D(j PayOrderJ) PayOrderD {
	return PayOrderD{
		j.UserID,
		j.Token,
		j.HouseID,
		j.HostID,
		j.OrderID,
		j.DiscountID,
		conv_PcJ_D(j.Pay),
		conv_tT_priDT(j.Time),
		conv_tT_priDT(j.Start),
		conv_tT_priDT(j.Stop),
		j.Result,
	}
}
func conv_POD_J(d PayOrderD) PayOrderJ {
	return PayOrderJ{
		d.UserID,
		d.Token,
		d.HouseID,
		d.HostID,
		d.OrderID,
		d.DiscountID,
		conv_PcD_J(d.Pay),
		conv_priDT_tT(d.Time),
		conv_priDT_tT(d.Start),
		conv_priDT_tT(d.Stop),
		d.Result,
	}
}
func conv_DcdJ_D(j DiscountdetailJ) DiscountdetailD {
	return DiscountdetailD{
		j.UserID,
		j.DiscountID,
		conv_dD_priDe(j.Reduce),
		j.Type,
		j.Description,
		conv_tT_priDT(j.Outdate),
		j.Useable,
	}
}
func conv_DcdD_J(d DiscountdetailD) DiscountdetailJ {
	return DiscountdetailJ{
		d.UserID,
		d.DiscountID,
		conv_priDe_dD(d.Reduce),
		d.Type,
		d.Description,
		conv_priDT_tT(d.Outdate),
		d.Useable,
	}
}
func conv_list_DcdJ_D(j []DiscountdetailJ) []DiscountdetailD {
	var d []DiscountdetailD
	for _, item := range j {
		d = append(d, conv_DcdJ_D(item))
	}
	return d
}
func conv_list_DcdD_J(d []DiscountdetailD) []DiscountdetailJ {
	var j []DiscountdetailJ
	for _, item := range d {
		j = append(j, conv_DcdD_J(item))
	}
	return j
}
func conv_list_POJ_D(j []PayOrderJ) []PayOrderD {
	var d []PayOrderD
	for _, item := range j {
		d = append(d, conv_POJ_D(item))
	}
	return d
}
func conv_list_POD_J(d []PayOrderD) []PayOrderJ {
	var j []PayOrderJ
	for _, item := range d {
		j = append(j, conv_POD_J(item))
	}
	return j
}
func conv_WSJ_D(j WalletStructJ) WalletStructD {
	return WalletStructD{
		j.UserID,
		j.Score,
		conv_dD_priDe(j.Balance),
		conv_list_DcdJ_D(j.DiscountList),
		conv_list_POJ_D(j.PayOrderList),
	}
}
func conv_WSD_J(d WalletStructD) WalletStructJ {
	return WalletStructJ{
		d.UserID,
		d.Score,
		conv_priDe_dD(d.Balance),
		conv_list_DcdD_J(d.DiscountList),
		conv_list_POD_J(d.PayOrderList),
	}
}
