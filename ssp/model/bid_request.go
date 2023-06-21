package model

import (
	"context"
	openrtb "github.com/prebid/openrtb/v19/openrtb2"
	"ssp/manager"
	"ssp/model/datastore"
	"ssp/utils"
)

func CreateSiteObject(ctx context.Context, bidRequest *openrtb.BidRequest, adslot *datastore.AdSlot) {
	bidRequest.Site = new(openrtb.Site)
	bidRequest.Site.ID = adslot.Site.SiteId
	bidRequest.Site.Domain = adslot.Site.Domain
}

func CreateDeviceObject(ctx context.Context, manager manager.Manager, bidRequest *openrtb.BidRequest, request *AdRequest) {
	bidRequest.Device = new(openrtb.Device)
	bidRequest.Device.OS = utils.GetOsFromUA(request.UA)
	bidRequest.Device.IP = request.IP
	bidRequest.Device.UA = request.UA
	bidRequest.Device.Geo = new(openrtb.Geo)
	bidRequest.Device.Geo.Country = utils.GetCountryFromIP(request.IP, manager)
}

func CreateUserObject(bidRequest *openrtb.BidRequest) {
	bidRequest.User = new(openrtb.User)
}
