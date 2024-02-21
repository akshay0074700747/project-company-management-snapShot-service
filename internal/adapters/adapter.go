package adapters

import (
	"context"
	"io"

	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/entities"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SnapShotAdapter struct {
	MinioDB *minio.Client
	MongoDb *mongo.Database
}

func NewSnapShotAdapter(miniodb *minio.Client, mongodb *mongo.Database) *SnapShotAdapter {
	return &SnapShotAdapter{
		MinioDB: miniodb,
		MongoDb: mongodb,
	}
}

func (snap *SnapShotAdapter) InsertSnapShot(ctx context.Context, fileName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) error {

	_, err := snap.MinioDB.PutObject(ctx, "snapshots-storage-bucket", fileName, reader, objectSize, opts)
	if err != nil {
		return err
	}

	return nil
}

func (snap *SnapShotAdapter) InsertSnapshotMetaDatas(req entities.SnapMessage) error {

	coll := snap.MongoDb.Collection("metaDatas")

	filter := bson.M{"project_id": req.ProjectID, "email": req.Email}
	update := bson.M{
		"$set": bson.M{
			"email":      req.Email,
			"project_id": req.ProjectID,
			"snapshots":  req.Snapshots,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := coll.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (snap *SnapShotAdapter) GetSnapshotData(email, projectID string) (entities.SnapMessage, error) {

	coll := snap.MongoDb.Collection("metaDatas")

	filter := bson.M{"project_id": projectID, "email": email}

	var res entities.SnapMessage

	if err := coll.FindOne(context.TODO(), filter).Decode(&res); err != nil {
		return entities.SnapMessage{}, err
	}

	return res, nil
}
