package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/akshay0074700747/projectandCompany_management_protofiles/pb/userpb"
	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/entities"
	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/helpers"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (snap *SnapShotService) StartConsumerGroup() {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        "host.docker.internal:9092",
		"group.id":                 "snapshotConsumers",
		"auto.offset.reset":        "earliest",
		"enable.auto.commit":       "false",
		"allow.auto.create.topics": true})
	if err != nil {
		helpers.PrintErr(err, "error occured at creating a kafka consumer")
		return
	}

	topic := "SnapshotTopic"

	err = consumer.Assign([]kafka.TopicPartition{
		{
			Topic:     &topic,
			Partition: 0,
			Offset:    kafka.OffsetStored,
		},
	})
	if err != nil {
		helpers.PrintErr(err, "Error assigning partitions")
		return
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(sigchan)
		consumer.Close()
	}()

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Received signal: %v\n", sig)
			run = false

		default:
			ev := consumer.Poll(1)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("Received message ")

				var msg entities.SnapMsg
				var meta entities.SnapMessage

				err := json.Unmarshal(e.Value, &msg)
				if err != nil {
					fmt.Printf("Error unmarshalling message value: %v\n", err)
					return
				}

				userID, err := snap.UserConn.GetByEmail(context.TODO(), &userpb.GetByEmailReq{
					Email: msg.Email,
				})
				if err != nil {
					fmt.Printf("Error retrieving userID from userservice", err)
					return
				}

				meta.UserID = userID.UserID
				meta.ProjectID = msg.ProjectID

				for _, snapShot := range msg.Snapshots {
					if err = snap.Usecase.InsertSnapShot(snapShot.Filename, snapShot.File, msg.Email, msg.ProjectID); err != nil {
						helpers.PrintErr(err, "Error happened at InsertSnapShot")
					}
					meta.Snapshots = append(meta.Snapshots, entities.Snapshot{
						Filename:    snapShot.Filename,
						Key:         snapShot.Key,
						Description: snapShot.Description,
						Progress:    snapShot.Progress,
					})
				}

				fmt.Println("=====hree=====", msg.IsStaged)

				if err = snap.Usecase.InsertMetaData(meta, msg.IsStaged, msg.Key); err != nil {
					helpers.PrintErr(err, "Error occured on InsertMetaData usecase")
					continue
				}

				_, err = consumer.CommitOffsets([]kafka.TopicPartition{e.TopicPartition})
				if err != nil {
					helpers.PrintErr(err, "Error committing offset")
				}
			case (kafka.Error):
				helpers.PrintErr(e, "errror occured at consumer")
			default:
				fmt.Printf("Ignored event: %v\n", e)

			}
		}
	}

	fmt.Println("Consumer shutting down...")

}
