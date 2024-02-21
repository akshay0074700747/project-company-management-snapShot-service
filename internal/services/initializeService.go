package services

import (
	"net/http"

	"github.com/go-chi/chi"
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

	r.Get("/project/snapshots",engine.Srv.GetSnapshots)

	http.ListenAndServe(addr, r)
}
