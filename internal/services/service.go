package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/entities"
	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/helpers"
	"github.com/akshay0074700747/projectandCompany_management_snapShot-service/internal/usecases"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type SnapShotService struct {
	Usecase usecases.SnapShotUsecaseInterfaces
}

func NewSnapShotService(usecase usecases.SnapShotUsecaseInterfaces) *SnapShotService {
	return &SnapShotService{
		Usecase: usecase,
	}
}

func (snap *SnapShotService) SendProjectSnapshot(w http.ResponseWriter, r *http.Request) {
	
}

func (snap *SnapShotService) SendSnapShot(w http.ResponseWriter, r *http.Request)  {
	
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
	}else {
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

	if res["email"] == "" || res["projectID"] == "" {
		http.Error(w, "the email and projectID cannot be empty", http.StatusInternalServerError)
		return
	}

	ress, err := snap.Usecase.GetSnapshotData(res["email"], res["projectID"])
	if err != nil {
		helpers.PrintErr(err, "error happened at GetSnapshotData")
		http.Error(w, "the email and projectID cannot be empty", http.StatusInternalServerError)
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

func (snap *SnapShotService) StartConsumerGroup() {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9092",
		"group.id":           "snapshotConsumers",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": "false"})
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

				meta.Email = msg.Email
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

				if err = snap.Usecase.InsertMetaData(meta); err != nil {
					helpers.PrintErr(err, "Error occured on InsertMetaData usecase")
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
