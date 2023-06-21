package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	_ "github.com/jarcoal/httpmock"
	openrtb "github.com/prebid/openrtb/v19/openrtb2"
	"github.com/stretchr/testify/assert"
	"net/http"
	manager2 "ssp/manager"
	"ssp/model"
	"ssp/model/datastore"
	"testing"
)

func TestSSpService(t *testing.T) {
	mockController := gomock.NewController(t)
	dspPubModel := model.NewMockIDSPPubModel(mockController)
	bidRequestStatModel := model.NewMockIBidRequestStatModel(mockController)
	adslotModel := model.NewMockIAdSlotModel(mockController)
	manager := manager2.Manager{Client: nil}
	sspService := NewSSPSvcHandler(dspPubModel, manager, bidRequestStatModel, adslotModel)
	ctx := context.Background()
	pubId := "64562217e9325b4c5bbe5433"
	adSlotId := "6458c5c947743c1a71fa9f7b"
	adType := "Banner"
	adMarkupHigh := "<iframe src=\"https://upload.wikimedia.org/wikipedia/commons/2/2a/New_Logo_AD.jpg\"></iframe>"
	adMarkuplow := `<iframe src="https://content.instructables.com/F7E/7YQE/IR6IMN04/F7E7YQEIR6IMN04.jpg?auto=webp&fit=bounds&frame=1&height=620&width=620"></iframe>`
	highbidResponse := openrtb.BidResponse{
		ID: "1",
		SeatBid: []openrtb.SeatBid{{
			Bid: []openrtb.Bid{{
				Price: 1,
				AdM:   adMarkupHigh,
			}},
			Seat:  "",
			Group: 0,
			Ext:   nil,
		}},
		BidID:      "1",
		Cur:        "",
		CustomData: "",
		NBR:        nil,
		Ext:        nil,
	}
	lowbidResponse := openrtb.BidResponse{
		ID: "1",
		SeatBid: []openrtb.SeatBid{{
			Bid: []openrtb.Bid{{
				Price: 0.5,
				AdM:   adMarkuplow,
			}},
			Seat:  "",
			Group: 0,
			Ext:   nil,
		}},
		BidID:      "1",
		Cur:        "",
		CustomData: "",
		NBR:        nil,
		Ext:        nil,
	}
	siteId := "64562217e9325b4c5bbe5443"
	adslot := &datastore.AdSlot{
		Id: adSlotId,
		Site: datastore.Site{
			SiteId: siteId,
			Domain: "http://localhost:3001",
		},
		PubId: pubId,
	}
	ip := "8.8.8.8"
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8"
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://dsp.com/dsp1",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{})
			return resp, err
		},
	)
	httpmock.RegisterResponder("POST", "https://dsp.com/high",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, highbidResponse)
			return resp, err
		},
	)
	httpmock.RegisterResponder("POST", "https://dsp.com/low",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, lowbidResponse)
			return resp, err
		},
	)

	t.Run("GetAdMarkup", func(t *testing.T) {
		t.Run("should throw error when ad slot id is missing", func(t *testing.T) {
			adRequest := &model.AdRequest{
				AdType: adType,
				//AdSlotId:adSlotId,
			}
			markup, err := sspService.GetAdMarkup(ctx, adRequest)
			assert.Empty(t, markup)
			assert.Equal(t, errors.New("adslotid  missing in GetAdMarkupPayload"), err)
		})

		t.Run("should throw error when adtype is missing", func(t *testing.T) {
			adRequest := &model.AdRequest{
				AdType:   "",
				AdSlotId: adSlotId,
			}
			markup, err := sspService.GetAdMarkup(ctx, adRequest)
			assert.Empty(t, markup)
			assert.Equal(t, errors.New("adtype missing in GetAdMarkupPayload"), err)
		})

		t.Run("should return markup as empty when no dsps are available", func(t *testing.T) {
			adRequest := &model.AdRequest{
				AdType:   adType,
				AdSlotId: adSlotId,
				Width:    "43",
				Height:   "43",
				IP:       ip,
				UA:       ua,
			}
			adslotModel.EXPECT().GetAdSlotDetails(ctx, manager.GetMongoClient(), adSlotId).Return(adslot, nil).Times(1)
			dspPubModel.EXPECT().
				GetAllDSPForPub(ctx, manager.GetMongoClient(), pubId).
				Return(nil, nil).
				Times(1)
			markup, err := sspService.GetAdMarkup(ctx, adRequest)
			assert.Empty(t, markup)
			assert.Equal(t, errors.New("no dsps are connected to this pub"), err)
		})

		t.Run("should return markup as empty when dsps retun empty response", func(t *testing.T) {

			adRequest := &model.AdRequest{
				AdType:   adType,
				AdSlotId: adSlotId,
				Width:    "43",
				Height:   "43",
				IP:       ip,
				UA:       ua,
			}
			dsps := []datastore.DSP{{
				Id:   "64562217e9325b4c5bbe5433",
				Name: "Dsp1",
				Url:  "https://dsp.com/dsp1",
			}}

			adslotModel.EXPECT().GetAdSlotDetails(ctx, manager.GetMongoClient(), adSlotId).Return(adslot, nil).Times(1)
			dspPubModel.EXPECT().
				GetAllDSPForPub(ctx, manager.GetMongoClient(), pubId).
				Return(dsps, nil).
				Times(1)

			bidRequestStatModel.EXPECT().
				SaveBidRequestStat(ctx, manager.GetMongoClient(), gomock.Any()).Return(nil).AnyTimes()
			markup, err := sspService.GetAdMarkup(ctx, adRequest)
			assert.Empty(t, markup)
			assert.Equal(t, nil, err)
		})

		t.Run("markup should contain highest bid and impression url", func(t *testing.T) {

			adRequest := &model.AdRequest{
				AdType:   adType,
				AdSlotId: adSlotId,
				Width:    "43",
				Height:   "43",
				IP:       ip,
				UA:       ua,
			}
			dsps := []datastore.DSP{
				{
					Id:   "64562217e9325b4c5bbe5434",
					Name: "Dsp2",
					Url:  "https://dsp.com/high",
				},
				{
					Id:   "64562217e9325b4c5bbe5435",
					Name: "Dsp3",
					Url:  "https://dsp.com/low",
				},
			}

			adslotModel.EXPECT().GetAdSlotDetails(ctx, manager.GetMongoClient(), adSlotId).Return(adslot, nil).Times(1)
			dspPubModel.EXPECT().
				GetAllDSPForPub(ctx, manager.GetMongoClient(), pubId).
				Return(dsps, nil).
				Times(1)

			bidRequestStatModel.EXPECT().
				SaveBidRequestStat(ctx, manager.GetMongoClient(), gomock.Any()).Return(nil).AnyTimes()

			markup, err := sspService.GetAdMarkup(ctx, adRequest)

			//fmt.Println(markup)
			assert.NotEmpty(t, markup)
			assert.Equal(t, nil, err)
			assert.Contains(t, markup, adMarkupHigh)
			assert.Contains(t, markup, "/imp")
		})
	})

	t.Run("RecordImpression", func(t *testing.T) {
		t.Run("bidrequestid not present", func(t *testing.T) {
			expectedErr := errors.New("bidReqId not set")
			err := sspService.RecordImpression(ctx, "")
			assert.Equal(t, expectedErr, err)
		})

		t.Run("bidrequest id not present in db", func(t *testing.T) {
			expectedErr := errors.New("no bid request available with this id")
			bidReqId := "1"
			bidRequestStatModel.EXPECT().GetBidRequestStat(ctx, manager.GetMongoClient(), bidReqId).Return(nil, nil).Times(1)
			err := sspService.RecordImpression(ctx, bidReqId)
			assert.Equal(t, expectedErr, err)
		})

		t.Run("should return success", func(t *testing.T) {
			bidReqId := "1"
			stat := &datastore.BidRequestStat{
				BidReqId:        bidReqId,
				Country:         "",
				OS:              "",
				PubId:           pubId,
				ExpectedRevenue: 0.5,
				Revenue:         0,
				ExpectedProfit:  0.3,
				Profit:          0,
				Impression:      0,
				DSPId:           "1",
			}
			bidRequestStatModel.EXPECT().
				GetBidRequestStat(ctx, manager.GetMongoClient(), bidReqId).
				Return(stat, nil).
				Times(1)

			bidRequestStatModel.EXPECT().
				UpdateImpression(ctx, manager.GetMongoClient(), bidReqId, stat).
				Return(nil).
				Times(1)
			err := sspService.RecordImpression(ctx, bidReqId)
			assert.Equal(t, nil, err)
		})

	})
	//todo test cases for win url
}
