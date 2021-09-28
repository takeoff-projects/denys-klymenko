package petsdb

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
)

var projectID string

// Pet model stored in Datastore
type Pet struct {
	Added   time.Time `firestore:"added" json:"added"`
	Caption string    `firestore:"caption" json:"caption"`
	Email   string    `firestore:"email" json:"email"`
	Image   string    `firestore:"image" json:"image"`
	Likes   int       `firestore:"likes" json:"likes"`
	Owner   string    `firestore:"owner" json:"owner"`
	Petname string    `firestore:"petname" json:"petname"`
	Name    string    // The ID used in the datastore.
}

// GetPets Returns all pets from datastore ordered by likes in Desc Order
func GetPets() ([]Pet, error) {
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

	var pets []Pet
	for _, snapshot := range all {
		var pet Pet
		err := snapshot.DataTo(&pet)
		if err != nil {
			log.Fatalf("Could not convert document to Pet type: %v", err)
		}

		pet.Name = snapshot.Ref.ID
		pets = append(pets, pet)
	}

	return pets, nil
}

func Add(pet Pet) error {
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
	_, _, err = firestoreClient.Collection("pets").Add(ctx, pet)
	if err != nil {
		log.Fatalf("Could not add pet: %v", err)
	}

	return nil

}
