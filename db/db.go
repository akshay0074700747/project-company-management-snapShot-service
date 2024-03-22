package db

import (
	"context"
	"fmt"
	"log"

	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMinio(cfg config.Config) *minio.Client {

	minioClient, err := minio.New(cfg.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AcessKey, cfg.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
		fmt.Println(err)
	}

	return minioClient
}

func ConnectMongo(cfg config.Config) *mongo.Database {

	clientOptions := options.Client().ApplyURI(cfg.MongoUrl)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database("snapShotDB")

}
