package usecases

import "github.com/akshay0074700747/projectandCompany_management_snapShot-service/entities"

type SnapShotUsecaseInterfaces interface {
	InsertSnapShot(string, []byte, string, string) error
	InsertMetaData(entities.SnapMessage) error
	GetSnapshotData(string, string) (entities.SnapMessage, error)
}
