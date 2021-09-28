package petsweb

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
	"os"
	"takeoff-projects/denys-klymenko/core/pets"
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
	_, _, err = firestoreClient.Collection("pets").Add(ctx, pet)
	if err != nil {
		log.Fatalf("Could not add pet: %v", err)
	}

	return nil

}
