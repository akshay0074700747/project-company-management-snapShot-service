package adapters

import (
	"context"
	"io"

	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/entities"
	"github.com/minio/minio-go/v7"
)

type SnapShotAdapterInterfaces interface {
	InsertSnapShot(context.Context, string, io.Reader, int64, minio.PutObjectOptions) error
	InsertSnapshotMetaDatas(entities.SnapMessage) error
	GetSnapshotData(string, string) (entities.SnapMessage, error)
}
