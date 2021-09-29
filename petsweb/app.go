package petsweb

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"takeoff-projects/denys-klymenko/core/pets"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

var projectID string

// Pet model stored in Datastore

// GetPets Returns all pets from datastore ordered by likes in Desc Order
func GetPets() ([]pets.Pet, error) {
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal(`You need to set the environment variable "GOOGLE_CLOUD_PROJECT"`)
	}

	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(context.TODO(), projectID)
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}

	defer func() {
		_ = firestoreClient.Close()
	}()
	// Create a query to fetch all Pet entities".
	all, err := firestoreClient.Collection("pets").OrderBy("likes", firestore.Desc).Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Could not get pets: %v", err)
	}

	var ps []pets.Pet
	for _, snapshot := range all {
		var pet pets.Pet
		err := snapshot.DataTo(&pet)
		if err != nil {
			log.Fatalf("Could not convert document to Pet type: %v", err)
		}

		pet.Name = snapshot.Ref.ID
		ps = append(ps, pet)
	}

	return ps, nil
}

func Add(pet pets.Pet) error {
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal(`You need to set the environment variable "GOOGLE_CLOUD_PROJECT"`)
	}

	ctx := context.Background()
	firestoreClient, err := firestore.NewClient(context.TODO(), projectID)
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}

	defer func() {
		_ = firestoreClient.Close()
	}()
	doc, _, err := firestoreClient.Collection("pets").Add(ctx, pet)
	if err != nil {
		log.Fatalf("Could not add pet: %v", err)
	}

	pet.Name = doc.ID

	emitSimilarImagesTask(projectID, pet)

	return nil

}

const region = "us-central1"

func emitSimilarImagesTask(projectID string, pet pets.Pet) {
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		fmt.Printf("NewClient: %v", err)
	}

	defer client.Close()

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, region, "images-queue")
	workerURL := fmt.Sprintf("https://%s-%s.cloudfunctions.net/collect-images", region, projectID)

	body := struct {
		ImageURL string `json:"image_url"`
		PetID    string `json:"pet_id"`
	}{
		ImageURL: pet.Image,
		PetID:    pet.Name,
	}

	bytes, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("cloudtasks.CreateTask marshal: %v", err)
	}

	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					Url:        workerURL,
					HttpMethod: taskspb.HttpMethod_POST,
					Body:       bytes,
				},
			},
		},
	}

	_, err = client.CreateTask(ctx, req)
	if err != nil {
		fmt.Printf("cloudtasks.CreateTask error: %v", err)
	}
}
