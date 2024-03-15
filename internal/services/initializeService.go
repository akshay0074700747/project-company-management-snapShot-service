package services

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

type SnapShotEngine struct {
	Srv *SnapShotService
}

func NewSnapShotEngine(srv *SnapShotService) *SnapShotEngine {
	return &SnapShotEngine{
		Srv: srv,
	}
}
func (engine *SnapShotEngine) Start(addr string) {

	r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	r.Get("/project/snapshots", engine.Srv.GetSnapshots)
	r.Get("/project/task/stages", engine.Srv.getStages)
	r.Get("/project/task/stages/count", engine.Srv.getStagesCount)
	r.Get("/snapshots/pull",engine.Srv.getSnapshotbyID)

	fmt.Println("snapShot service is starting")
	http.ListenAndServe(addr, handler)
}
