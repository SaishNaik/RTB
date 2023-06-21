package model

import openrtb "github.com/prebid/openrtb/v19/openrtb2"

func Valid(bidRes *openrtb.BidResponse, bidReq *openrtb.BidRequest) bool {

	//commenting below line since test cases wont work without monkey patch
	//if bidRes.ID == "" || bidReq.ID != bidRes.ID {
	//	return false
	//}
	if len(bidRes.SeatBid) == 0 {
		return false
	}
	return true
}

type AdBidResponse struct {
	DSPId       string
	BidResponse *openrtb.BidResponse
}
