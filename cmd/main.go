package main

import (
	"log"

	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/config"
	injectdependency "github.com/akshay0074700747/projectandCompany_management_snapShot-service/injectDependency"
)

func main() {

	config, err := config.LoadConfigurations()
	if err != nil {
		log.Fatal("cannot load configurations", err)
	}

	engine := injectdependency.Initialize(config)

	engine.Start(":50005")
}
