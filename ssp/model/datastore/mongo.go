package datastore

//todo change DSP model, add more fields
type DSP struct {
	Id   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Url  string `json:"url" bson:"url"`
}

type DSPPub struct {
	PubId   string `json:"pub_id" bson:"pub_id"`
	DSPList []DSP  `json:"dsps" bson:"dsps"`
}

//todo current assumption one bid request one impression
type BidRequestStat struct {
	BidReqId        string  `json:"bid_req_id" bson:"bid_req_id"`
	Country         string  `json:"country" bson:"country"`
	OS              string  `json:"os" bson:"os"`
	ExpectedRevenue float64 `json:"expected_revenue" bson:"expected_revenue"`
	Revenue         float64 `json:"revenue" bson:"revenue"`
	ExpectedProfit  float64 `json:"expected_profit" bson:"expected_profit"`
	Profit          float64 `json:"profit" bson:"profit"`
	Impression      int32   `json:"impression" bson:"impression"`
	DSPId           string  `json:"dsp_id" bson:"dsp_id"`
	PubId           string  `json:"pub_id" bson:"pub_id"`
}

type AdSlot struct {
	Id    string `json:"_id" bson:"_id"`
	Site  Site   `json:"site" bson:"site"`
	PubId string `json:"pub_id" bson:"pub_id"`
}

type Site struct {
	SiteId string `json:"id" bson:"id"`
	Domain string `json:"domain" bson:"domain"`
}

//type AdSlot struct {
//	Id     primitive.ObjectID `json:"_id" bson:"_id"`
//	PubId  string             `json:"pub_id" bson:"pub_id"`
//	Adtype string             `json:"ad_type" bson:"ad_type"`
//}
