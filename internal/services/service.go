package services

import (
	"encoding/json"
	"net/http"

	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/userpb"
	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/entities"
	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/helpers"
	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/internal/usecases"
)

type SnapShotService struct {
	Usecase  usecases.SnapShotUsecaseInterfaces
	UserConn userpb.UserServiceClient
}

func NewSnapShotService(usecase usecases.SnapShotUsecaseInterfaces, usrAddr string) *SnapShotService {
	userConn, _ := helpers.DialGrpc(usrAddr)
	return &SnapShotService{
		Usecase:  usecase,
		UserConn: userpb.NewUserServiceClient(userConn),
	}
}

func (snap *SnapShotService) SendProjectSnapshot(w http.ResponseWriter, r *http.Request) {

}

func (snap *SnapShotService) SendSnapShot(w http.ResponseWriter, r *http.Request) {

}

func (snap *SnapShotService) GetTotalProjectVersionsSnapshots(w http.ResponseWriter, r *http.Request) {

}

func (snap *SnapShotService) GetTotalProjectSnapshots(w http.ResponseWriter, r *http.Request) {

	var res = make(map[string]string)
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		helpers.PrintErr(err, "error happened decoding the json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res["versionID"] == "" {
		//show the total snapshots for the specified version
	} else {
		//show the total snapshots for the current version
	}
}

func (snap *SnapShotService) GetSnapshots(w http.ResponseWriter, r *http.Request) {

	var res = make(map[string]string)
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		helpers.PrintErr(err, "error happened decoding the json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res["userID"] == "" || res["projectID"] == "" {
		http.Error(w, "the userID and projectID cannot be empty", http.StatusInternalServerError)
		return
	}

	ress, err := snap.Usecase.GetSnapshotData(res["usrID"], res["projectID"])
	if err != nil {
		helpers.PrintErr(err, "error happened at GetSnapshotData")
		http.Error(w, "the userID and projectID cannot be empty", http.StatusInternalServerError)
		return
	}

	jsonDta, err := json.Marshal(ress)
	if err != nil {
		helpers.PrintErr(err, "error happened at marshalling to json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonDta)
}

func (snap *SnapShotService) getStages(w http.ResponseWriter, r *http.Request) {

	queries := r.URL.Query()

	ress, err := snap.Usecase.GetStages(queries.Get("userID"), queries.Get("projectID"))
	if err != nil {
		helpers.PrintErr(err, "error happened at GetStages usecase")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rees entities.StageRes

	rees.Stages = len(ress)
	rees.Details = ress

	jsonDta, err := json.Marshal(rees)
	if err != nil {
		helpers.PrintErr(err, "error happened at marshaling to json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonDta)
}

func (snap *SnapShotService) getStagesCount(w http.ResponseWriter, r *http.Request) {

	queries := r.URL.Query()

	usersProgress, err := snap.Usecase.GetStagesCount(queries.Get("projectID"))
	if err != nil {
		helpers.PrintErr(err, "error happened at GetStagesCount usecase")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var res entities.ListofUserProgress
	res.UserAndProgress = usersProgress

	jsonDta, err := json.Marshal(res)
	if err != nil {
		helpers.PrintErr(err, "error happened at marshaling to json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonDta)

}
