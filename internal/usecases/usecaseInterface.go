package usecases

import "github.com/akshay0074700747/projectandCompany_management_snapShot-service/entities"

type SnapShotUsecaseInterfaces interface {
	InsertSnapShot(string, []byte, string, string) error
	InsertMetaData(entities.SnapMessage, bool, string) error
	GetSnapshotData(string, string) (entities.SnapMessage, error)
	GetStages(string, string) ([]entities.StagesDetails, error)
	GetStagesCount(string)([]entities.UserProgress,error)
}
