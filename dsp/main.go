package main

import (
	"github.com/gin-gonic/gin"
	openrtb "github.com/prebid/openrtb/v19/openrtb2"
	"math/rand"
	"net/http"
)

func main() {
	router := gin.Default()
	router.POST("/dsp1", func(c *gin.Context) {

		var bidReq *openrtb.BidRequest
		err := c.BindJSON(&bidReq)
		if err != nil {
			return
		}
		bidResponse := openrtb.BidResponse{
			ID: bidReq.ID,
			SeatBid: []openrtb.SeatBid{{
				Bid: []openrtb.Bid{{
					Price: rand.Float64(),
					AdM:   "<iframe src=\"https://upload.wikimedia.org/wikipedia/commons/2/2a/New_Logo_AD.jpg\"></iframe>",
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
		c.IndentedJSON(http.StatusOK, bidResponse)
	})

	router.POST("/dsp2", func(c *gin.Context) {
		var bidReq *openrtb.BidRequest
		err := c.BindJSON(&bidReq)
		if err != nil {
			return
		}
		bidResponse := openrtb.BidResponse{
			ID: bidReq.ID,
			SeatBid: []openrtb.SeatBid{{
				Bid: []openrtb.Bid{{
					Price: rand.Float64(),
					NURL:  "http://dsp:3002/win_url?bid_request_id=${AUCTION_ID}",
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
		//jsonResp, err := json.Marshal(bidResponse)
		//if err != nil {
		//	fmt.Println(err)
		//}
		c.IndentedJSON(http.StatusOK, bidResponse)
	})

	router.GET("/win_url", func(c *gin.Context) {
		c.Writer.Write([]byte(`<iframe src="https://content.instructables.com/F7E/7YQE/IR6IMN04/F7E7YQEIR6IMN04.jpg?auto=webp&fit=bounds&frame=1&height=620&width=620"></iframe>`))
	})

	http.ListenAndServe(":3002", router)
}
