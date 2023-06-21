package main

import (
	"flag"
	"log"
	"ssp/app"
)

func main() {
	//todo add comment for flags
	var configFilePath, serverPort string
	flag.StringVar(&configFilePath, "config", "config.json", "absolute path to the configuration file")
	flag.StringVar(&serverPort, "server_port", "3000", "port on which server runs")
	flag.Parse()

	//create new application
	sspApp := app.New(configFilePath)
	// Initiate App
	sspApp.Init(configFilePath)
	//Start App
	err := sspApp.Start()
	if err != nil {
		log.Fatalf("unable to start application, %v", err)
	}

	//	http.HandleFunc("/", router.GetAdMarkup)
	//http.ListenAndServe(":3000", nil)
}
