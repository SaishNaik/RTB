package dicontainer

import (
	"ssp/controller"
	"ssp/datastore"
	manager2 "ssp/manager"
	"ssp/model"
	"ssp/service"
)

//dependency injection container
type IDiContainer interface {
	StartDependenciesInjection(manager manager2.Manager)
	GetDiContainer() *DiContainer
}

type DiContainer struct {
	SSPController controller.ISSPController
}

func NewDiContainer() IDiContainer {
	return &DiContainer{}
}

func (d *DiContainer) StartDependenciesInjection(manager manager2.Manager) {
	//datastore
	mongostore := datastore.InitMongoStore()

	// todo collection name from config
	dspPubModel := model.InitDSPPubModel(mongostore, "dsp_pub")
	adSlotModel := model.InitAdSlotModel(mongostore, "ad_slot")
	bidRequestStatModel := model.InitBidRequestStatModel(mongostore, "bid_request_stat")
	sspService := service.NewSSPSvcHandler(dspPubModel, manager, bidRequestStatModel, adSlotModel)
	sspController := controller.NewSSPController(sspService)
	d.SSPController = sspController
}

func (d *DiContainer) GetDiContainer() *DiContainer {
	return d
}
