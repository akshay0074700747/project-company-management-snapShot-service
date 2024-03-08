package injectdependency

import (
	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/config"
	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/db"
	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/internal/adapters"
	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/internal/services"
	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/internal/usecases"
)

func Initialize(cfg config.Config) *services.SnapShotEngine {

	minioDB := db.ConnectMinio(cfg)
	mongoDB := db.ConnectMongo(cfg)
	adapter := adapters.NewSnapShotAdapter(minioDB, mongoDB)
	usecase := usecases.NewSnapShotUseCases(adapter)
	service := services.NewSnapShotService(usecase, ":50001")

	go service.StartConsumerGroup()

	return services.NewSnapShotEngine(service)
}
