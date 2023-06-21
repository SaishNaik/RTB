package app

import (
	"fmt"
	"github.com/ip2location/ip2location-go/v9"
	"log"
	"net/http"
	"ssp/config"
	"ssp/dicontainer"
	"ssp/lib/mongo"
	manager2 "ssp/manager"
	"ssp/router"
)

type ISSPSvcApp interface {
	Init(string)
	Start() error
	GetMainConfig() *config.MainConfig
}

type SSPSvcApp struct {
	mainConfig *config.MainConfig
	httpServer *http.Server
	router     router.IRouter
	//db         da
}

func New(filePath string) ISSPSvcApp {
	app := new(SSPSvcApp)
	app.mainConfig = config.LoadMainConfig(filePath)
	//todo logging
	return app
}

func (app *SSPSvcApp) GetMainConfig() *config.MainConfig {
	return app.mainConfig
}

func (app *SSPSvcApp) Init(filePath string) {
	mainConfig := app.GetMainConfig()

	//db setup
	mongoClient := &mongo.Client{}
	var mongoInitialized bool
	var err error
	if mainConfig.MongoDB.Enabled {
		mongoInitialized, err = mongoClient.Auth(mainConfig.MongoDB.Uri, mainConfig.MongoDB.Timeout)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	if !mongoInitialized {
		panic("Mongo initialisation failed")
		return
	}
	ipClient, err := ip2location.OpenDB(mainConfig.IP2LocationDBPATH)
	if err != nil {
		panic("ipClient initialisation failed")
		return
	}
	manager := manager2.Manager{Client: mongoClient, IPClient: ipClient}

	//router setup

	//dependency injection
	di := dicontainer.NewDiContainer()
	di.StartDependenciesInjection(manager)

	app.router = router.NewRouter(app.mainConfig.GinMode)
	app.router.InitRoutes(di)

	app.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", app.mainConfig.Application.Port),
		Handler: app.router.GetMux(),
	}
}

func (app *SSPSvcApp) Start() error {
	log.Println("########## Server starting ##########")
	err := app.httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("unable to start server, %v", err)
		return err
	}
	return nil
}
