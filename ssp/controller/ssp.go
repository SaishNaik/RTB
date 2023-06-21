package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ssp/model"
	"ssp/service"
)

type ISSPController interface {
	AdRequest(*gin.Context)
	Impression(*gin.Context)
}

type SSPController struct {
	service service.ISSPSvcHandler
}

func NewSSPController(service service.ISSPSvcHandler) ISSPController {
	return &SSPController{
		service: service,
	}
}

func (ssp *SSPController) AdRequest(ctx *gin.Context) {
	//markup := `<img src="https://upload.wikimedia.org/wikipedia/commons/2/2a/New_Logo_AD.jpg">`
	//ctx.Writer.Write([]byte(markup))

	context := ctx.Request.Context()
	//todo requestid

	var adRequest *model.AdRequest
	err := ctx.BindJSON(&adRequest)
	if err != nil {
		//todo error response
		fmt.Println(err)
		return
	}

	adRequest.UA = ctx.Request.UserAgent()
	adRequest.IP = ctx.ClientIP()

	//validate Request
	err = adRequest.Validate()
	if err != nil {
		//todo error response
		return
	}

	// call service
	markup, err := ssp.service.GetAdMarkup(context, adRequest)
	if err != nil {
		//todo error response
		return
	}

	//todo error handling
	ctx.Writer.Write([]byte(markup))
}

func (ssp *SSPController) Impression(ctx *gin.Context) {
	context := ctx.Request.Context()
	//todo requestid

	bidReqId := ctx.Param("bidReqId")
	if bidReqId == "" {
		//todo error
	}

	// call service to record impression
	err := ssp.service.RecordImpression(context, bidReqId)
	if err != nil {
		//todo error response
		return
	}

	//todo error handling
	//todo send correct response
	ctx.Writer.Write([]byte("recorded"))
}
