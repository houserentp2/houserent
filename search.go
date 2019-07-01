package main

import (
	"github.com/go-ego/riot/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)
var houses=map[string]SearchFocus{}
type SearchFocus struct{
	HouseID string `json:"houseid"`
	Description string `json:"description"`
	Time primitive.DateTime `json:"time"`
	ClickCount int64 `json:"clickcount"`
}
type HouseScoringFields struct{
	Timestamp int64
	ClickCount int64
}
type HouseScoringCriteria struct{

}
func (criteria HouseScoringCriteria)Score(doc types.IndexedDoc,fields interface{})[]float32{
	if reflect.TypeOf(fields)!=reflect.TypeOf(HouseScoringFields{}){
		return []float32{}
	}
	wsf := fields.(HouseScoringFields)
	output := make([]float32, 3)
	if doc.TokenProximity > MaxTokenProximity {
		output[0] = 1.0 / float32(doc.TokenProximity)
	} else {
		output[0] = 1.0
	}
	output[1] = float32(wsf.Timestamp / (SecondsInADay * 3))
	output[2] = float32(doc.BM25 * (1 + float32(wsf.ClickCount)/10000))
	return output
}
func indexHouseinfo(d HouseDetailD){
	description:=d.Title+" "+d.Description+" "+d.Location.Province+""+d.Location.City+""+d.Location.Zone+""+d.Location.Path
	houses[d.HouseID]=SearchFocus{
		HouseID:d.HouseID,
		Description:description,
		Time:d.Time,
		ClickCount:d.ClickCount,
	}

	searcher.Index(d.HouseID,types.DocData{
		Content:description,
		Fields:HouseScoringFields{
			Timestamp:int64(d.Time),
			ClickCount:d.ClickCount,
		},
	})
	searcher.Flush()
}
func searchfor(param string)[]string{
	output:=searcher.SearchDoc(types.SearchReq{
		Text:param,
		RankOpts:&types.RankOpts{
			ScoringCriteria: &HouseScoringCriteria{},
			OutputOffset:0,
			MaxOutputs:100,
		},
	})
	docs:=[]string{}
	for _,doc:=range output.Docs{
		docs=append(docs,doc.DocId )
	}
	return docs
}