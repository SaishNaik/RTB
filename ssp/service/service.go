package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	openrtb "github.com/prebid/openrtb/v19/openrtb2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"ssp/constant"
	manager2 "ssp/manager"
	"ssp/model"
	"ssp/model/datastore"
	"strconv"
	"strings"
	"sync"
)

//todo remove commission from bid response

//todo check if functions added need to be added in interface
type ISSPSvcHandler interface {
	GetAdMarkup(ctx context.Context, request *model.AdRequest) (markup string, error error)
	RecordImpression(ctx context.Context, bidReqId string) (error error)
	ConductAuction(ctx context.Context, bidRequest *openrtb.BidRequest, pubId string) (*model.AdBidResponse, error)
}

type SSPSvcHandler struct {
	DSPModel        model.IDSPPubModel
	Manager         manager2.Manager
	BidReqStatModel model.IBidRequestStatModel
	AdSlotModel     model.IAdSlotModel
}

func NewSSPSvcHandler(dspModel model.IDSPPubModel, manager manager2.Manager, bidReqStatModel model.IBidRequestStatModel, adSlotModel model.IAdSlotModel) ISSPSvcHandler {
	return &SSPSvcHandler{DSPModel: dspModel, Manager: manager, BidReqStatModel: bidReqStatModel, AdSlotModel: adSlotModel}
}

func (s *SSPSvcHandler) GetAdMarkup(ctx context.Context, request *model.AdRequest) (markup string, error error) {
	//validate payload
	err := s.ValidateGetAdMarkupPayload(ctx, request)
	if err != nil {
		return "", err
	}

	adSlotDetails, err := s.AdSlotModel.GetAdSlotDetails(ctx, s.Manager.GetMongoClient(), request.AdSlotId)

	if err != nil {
		return "", err
	}
	fmt.Println("adslotdetails", adSlotDetails)
	if adSlotDetails.Id == "" {
		return "", err
	}

	// Create a BidRequest from AdRequestObject
	bidRequest, err := s.createBidRequest(ctx, request, adSlotDetails)
	if err != nil {
		return "", err
	}

	pubId := adSlotDetails.PubId
	winAdBidResponse, err := s.ConductAuction(ctx, bidRequest, pubId)
	if err != nil {
		return "", err
	}

	if winAdBidResponse != nil {
		markup, err = s.GetMarkupFromBidResponse(ctx, bidRequest, winAdBidResponse.BidResponse)
		if err != nil {
			fmt.Println(err)
			return "", err
		}

		markup, err = s.setUpImpressionUrlInMarkup(ctx, markup, bidRequest.ID)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
	}

	go func() {
		err := s.saveBidRequestStat(ctx, bidRequest, winAdBidResponse, pubId)
		if err != nil {
			fmt.Println(err)
		}
	}()

	//fmt.Println(markup, err)
	return markup, err

}

func (s *SSPSvcHandler) RecordImpression(ctx context.Context, bidReqId string) (error error) {

	if bidReqId == "" {
		return errors.New("bidReqId not set")
	}
	getStat, err := s.BidReqStatModel.GetBidRequestStat(ctx, s.Manager.GetMongoClient(), bidReqId)
	if err != nil {
		return err
	}
	if getStat == nil {
		return errors.New("no bid request available with this id")
	}

	err = s.BidReqStatModel.UpdateImpression(ctx, s.Manager.GetMongoClient(), bidReqId, getStat)
	return err
}

func (s *SSPSvcHandler) setUpImpressionUrlInMarkup(ctx context.Context, markup, bidRequestId string) (string, error) {
	markup = markup + fmt.Sprintf(`<img width="1" height="1" src="http://localhost:3000/imp/%s">`, bidRequestId)
	return markup, nil
}

func (s *SSPSvcHandler) saveBidRequestStat(ctx context.Context, bidRequest *openrtb.BidRequest, winAdBidResponse *model.AdBidResponse, pubId string) error {
	//todo actually add correct value and do validation
	var bidPrice float64
	var dspId string
	if winAdBidResponse != nil {
		winResponse := winAdBidResponse.BidResponse
		bidPrice = winResponse.SeatBid[0].Bid[0].Price
		dspId = winAdBidResponse.DSPId
	}

	stat := &datastore.BidRequestStat{
		BidReqId:        bidRequest.ID,
		Country:         bidRequest.Device.Geo.Country,
		OS:              bidRequest.Device.OS,
		ExpectedRevenue: bidPrice,
		PubId:           pubId,
		ExpectedProfit:  bidPrice * 0.30, // todo stop hardcoding profit
		DSPId:           dspId,
	}

	return s.BidReqStatModel.SaveBidRequestStat(ctx, s.Manager.GetMongoClient(), stat)
}

func (s *SSPSvcHandler) GetMarkupFromBidResponse(ctx context.Context, bidRequest *openrtb.BidRequest, winResponse *openrtb.BidResponse) (string, error) {
	//check if admarkup is present in bid response
	adm := winResponse.SeatBid[0].Bid[0].AdM
	if adm != "" {
		return adm, nil
	} else {
		winUrl := winResponse.SeatBid[0].Bid[0].NURL
		return s.GetMarkupFromWinUrl(ctx, bidRequest, winResponse, winUrl)
	}
}

//todo for each object in bidrequest, check if other fields from struct are needed
func (s *SSPSvcHandler) createBidRequest(ctx context.Context, request *model.AdRequest, adSlotDetails *datastore.AdSlot) (*openrtb.BidRequest, error) {
	//if request.Typegomock.Any().String()
	//todo banner type
	bidRequest := &openrtb.BidRequest{}
	bidRequest.ID = primitive.NewObjectID().Hex()

	width, err := strconv.Atoi(request.Width)
	if err != nil {
		return nil, err
	}

	height, err := strconv.Atoi(request.Height)
	if err != nil {
		return nil, err
	}

	if request.AdType == constant.Banner {
		//todo add BidRequest other fields
		bidRequest.Imp = []openrtb.Imp{{
			ID: "1", //todo make auto increment
			//todo add Banner other fields
			Banner: &openrtb.Banner{
				Format: []openrtb.Format{{
					W: int64(width),
					H: int64(height),
				}},
			},
		}}
	}

	//model.(adSlotDetails)
	model.CreateSiteObject(ctx, bidRequest, adSlotDetails)
	model.CreateDeviceObject(ctx, s.Manager, bidRequest, request)
	return bidRequest, nil
}

//func (s *SSPSvcHandler) RecordImpression(ctx context.Context, bidRequestId string) error {
//
//}

//func (s *SSPSvcHandler) RecordImpression(ctx context.Context, bidRequestId string) error {
//
//	//check if bid request exists
//	s.BidReqModel
//
//	s.BidReqModel.StoreImpression(ctx)
//}

//todo error processing
func (s *SSPSvcHandler) ConductAuction(ctx context.Context, bidRequest *openrtb.BidRequest, pubId string) (*model.AdBidResponse, error) {

	// get all dsps
	dsps, err := s.DSPModel.GetAllDSPForPub(ctx, s.Manager.GetMongoClient(), pubId)
	if err != nil {
		return nil, err
	}
	if len(dsps) == 0 {
		return nil, errors.New("no dsps are connected to this pub")
	}

	var wg sync.WaitGroup
	ch := make(chan *model.AdBidResponse)
	for _, dsp := range dsps {
		wg.Add(1)
		go reqToDSP(ctx, bidRequest, dsp, &wg, ch)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	var winningAdBidResponse *model.AdBidResponse
	maxPrice := 0.0
	for adBidResponse := range ch {
		//todo validation for bid response from adv
		res := adBidResponse.BidResponse
		if res != nil && model.Valid(res, bidRequest) { //todo this in reqToDSP, send res as nil if some issue with response
			price := res.SeatBid[0].Bid[0].Price
			if price > maxPrice {
				maxPrice = price
				winningAdBidResponse = adBidResponse
			}
		}
	}
	return winningAdBidResponse, nil
}

//todo error processing
func reqToDSP(ctx context.Context, bidRequest *openrtb.BidRequest, dsp datastore.DSP, wg *sync.WaitGroup, ch chan<- *model.AdBidResponse) {
	defer wg.Done()
	url := dsp.Url

	bidRequestJson, err := json.Marshal(bidRequest)
	if err != nil {

	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bidRequestJson))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")

	//todo add timeout for client
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	bidResponse := &openrtb.BidResponse{}
	err = json.NewDecoder(res.Body).Decode(bidResponse)
	if err != nil {
		ch <- nil
	}

	adBidResponse := &model.AdBidResponse{
		DSPId:       dsp.Id,
		BidResponse: bidResponse,
	}

	ch <- adBidResponse
}

func (s *SSPSvcHandler) GetMarkupFromWinUrl(ctx context.Context, bidRequest *openrtb.BidRequest, bidResponse *openrtb.BidResponse, winUrl string) (string, error) {

	winUrl, err := s.substituteMacros(ctx, bidRequest, bidResponse, winUrl)
	if err != nil {
		return "", err
	}
	return s.CallWinUrl(ctx, winUrl)
}

func (s *SSPSvcHandler) CallWinUrl(ctx context.Context, winUrl string) (string, error) {
	r, err := http.NewRequest(http.MethodGet, winUrl, nil)
	if err != nil {
		return "", err
	}

	//todo add timeout for client
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(resBody), nil
}

func (s *SSPSvcHandler) substituteMacros(ctx context.Context, bidRequest *openrtb.BidRequest, bidResponse *openrtb.BidResponse, winUrl string) (string, error) {
	//todo later for all macros, also check each macro if correct value is specified
	winUrl = strings.Replace(winUrl, "${AUCTION_ID}", bidRequest.ID, 1)
	winUrl = strings.Replace(winUrl, "${AUCTION_BID_ID}", bidResponse.BidID, 1)
	winUrl = strings.Replace(winUrl, "${AUCTION_PRICE}", fmt.Sprintf("%f", bidResponse.SeatBid[0].Bid[0].Price), 1)
	//winUrl := strings.Replace(winUrl,"${AUCTION_ID}",bidRequest.ID,1)
	//winUrl := strings.Replace(winUrl,"${AUCTION_ID}",bidRequest.ID,1)
	//winUrl := strings.Replace(winUrl,"${AUCTION_ID}",bidRequest.ID,1)
	return winUrl, nil
}
