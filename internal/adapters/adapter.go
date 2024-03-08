package adapters

import (
	"context"
	"fmt"
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

func (snap *SnapShotAdapter) InsertSnapshotMetaDatas(req entities.SnapMessage, isStaged bool, key string) error {

	var update bson.M
	if isStaged {
		req.Snapshots[len(req.Snapshots)-1].IsStaged = true
		fmt.Println(req.Snapshots[len(req.Snapshots)-1].IsStaged, "===here-=====")
		update = bson.M{
			"$set": bson.M{
				"user_id":    req.UserID,
				"project_id": req.ProjectID,
				"key":        key,
			},
			"$addToSet": bson.M{
				"snapshots": bson.D{{Key: "$each", Value: req.Snapshots}},
			},
		}
	} else {
		update = bson.M{
			"$set": bson.M{
				"user_id":    req.UserID,
				"project_id": req.ProjectID,
			},
			"$addToSet": bson.M{
				"snapshots": bson.D{{Key: "$each", Value: req.Snapshots}},
			},
		}
	}

	coll := snap.MongoDb.Collection("metaDatas")

	filter := bson.M{"project_id": req.ProjectID, "user_id": req.UserID}
	opts := options.Update().SetUpsert(true)

	_, err := coll.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (snap *SnapShotAdapter) GetSnapshotData(userID, projectID string) (entities.SnapMessage, error) {

	coll := snap.MongoDb.Collection("metaDatas")

	filter := bson.M{"project_id": projectID, "user_id": userID}

	var res entities.SnapMessage

	if err := coll.FindOne(context.TODO(), filter).Decode(&res); err != nil {
		return entities.SnapMessage{}, err
	}

	return res, nil
}

func (snap *SnapShotAdapter) GetStages(userID, projectID string) ([]entities.StagesDetails, error) {

	coll := snap.MongoDb.Collection("metaDatas")

	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"project_id", projectID}, {"user_id", userID}}}},
		bson.D{{"$unwind", "$snapshots"}},
		bson.D{{"$match", bson.D{{"snapshots.is_staged", true}}}},
		bson.D{{"$project", bson.D{
			{"_id", 0},
			{"filename", "$snapshots.filename"},
			{"key", "$snapshots.key"},
			{"description", "$snapshots.description"},
		}}},
	}

	cursor, err := coll.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var stages []entities.StagesDetails
	for cursor.Next(context.Background()) {
		var stage entities.StagesDetails
		if err := cursor.Decode(&stage); err != nil {
			return nil, err
		}
		stages = append(stages, stage)
	}

	return stages, nil
}

func (snap *SnapShotAdapter) GetStagesCount(projectID string) ([]entities.UserProgress, error) {

	coll := snap.MongoDb.Collection("metaDatas")

	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"project_id", projectID}}}},
		bson.D{{"$unwind", "$snapshots"}},
		bson.D{{"$match", bson.D{{"snapshots.is_staged", true}}}},
		bson.D{{"$group", bson.D{
			{"_id", "$user_id"},
			{"stages", bson.D{{"$sum", 1}}},
		}}},
		bson.D{{"$project", bson.D{
			{"_id", 0},
			{"user_id", "$_id"},
			{"stages", 1},
		}}},
	}

	cursor, err := coll.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var stages []entities.UserProgress
	for cursor.Next(context.Background()) {
		var stage entities.UserProgress
		if err := cursor.Decode(&stage); err != nil {
			return nil, err
		}
		stages = append(stages, stage)
	}

	fmt.Println(stages)

	return stages, nil
}
