package usecases

import (
	"bytes"
	"context"

	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/entities"
	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/helpers"
	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/internal/adapters"
	"github.com/minio/minio-go/v7"
)

type SnapShotUseCases struct {
	Adapter adapters.SnapShotAdapterInterfaces
}

func NewSnapShotUseCases(adapter adapters.SnapShotAdapterInterfaces) *SnapShotUseCases {
	return &SnapShotUseCases{
		Adapter: adapter,
	}
}

func (snap *SnapShotUseCases) InsertSnapShot(fileName string, reader []byte, email, projectId string) error {

	newReader := bytes.NewReader(reader)
	err := snap.Adapter.InsertSnapShot(context.TODO(), fileName, newReader, newReader.Size(), minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"email":     email,
			"projectID": projectId,
		},
	})
	if err != nil {
		helpers.PrintErr(err, "error happened at InsertSnapShot adapter")
		return err
	}

	return nil
}

func (snap *SnapShotUseCases) InsertMetaData(req entities.SnapMessage, isStaged bool, key string) error {

	if err := snap.Adapter.InsertSnapshotMetaDatas(req, isStaged, key); err != nil {
		helpers.PrintErr(err, "error occured at InsertSnapshotMetaDatas adapter")
		return err
	}

	return nil
}

func (snap *SnapShotUseCases) GetSnapshotData(email, projetID string) (entities.SnapMessage, error) {

	res, err := snap.Adapter.GetSnapshotData(email, projetID)
	if err != nil {
		helpers.PrintErr(err, "error occured at GetSnapshotData adapter")
		return entities.SnapMessage{}, err
	}

	return res, err
}

func (snap *SnapShotUseCases) GetStages(userID, projectID string) ([]entities.StagesDetails, error) {

	res, err := snap.Adapter.GetStages(userID, projectID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetStages adapter")
		return nil, err
	}

	return res, nil
}

func (snap *SnapShotUseCases) GetStagesCount(projectID string) ([]entities.UserProgress, error) {

	res, err := snap.Adapter.GetStagesCount(projectID)
	if err != nil {
		helpers.PrintErr(err, "error happened at GetStagesCount adapter")
		return nil, err
	}

	return res, nil
}
