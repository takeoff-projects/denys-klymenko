package pets

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"log"
	"takeoff-projects/denys-klymenko/core/pets"
)

type PetsDAL struct {
	Client    *firestore.Client
	ProjectID string
}

func (d *PetsDAL) Create(ctx context.Context, pet pets.Create) error {
	_, _, err := d.Client.Collection("pets").Add(ctx, pet)
	if err != nil {
		log.Fatalf("Could not create pet: %v", err)
	}

	return nil
}

func (d *PetsDAL) Get(ctx context.Context, petID string) (pets.Pet, error) {
	var pet pets.Pet
	doc, err := d.Client.Collection("pets").Doc(petID).Get(ctx)
	if err != nil {
		return pet, fmt.Errorf("could get pet: %w", err)
	}

	err = doc.DataTo(&pet)
	if err != nil {
		return pet, fmt.Errorf("could map pet: %w", err)
	}

	return pet, nil
}

func (d *PetsDAL) Update(ctx context.Context, petID string, update pets.Update) error {
	_, err := d.Client.Collection("pets").Doc(petID).Update(ctx, Updates(update))
	if err != nil {
		log.Fatalf("Could not update pet: %v", err)
	}

	return nil
}

func (d *PetsDAL) List(ctx context.Context) ([]pets.Pet, error) {
	all, err := d.Client.Collection("pets").OrderBy("likes", firestore.Desc).Documents(ctx).GetAll()
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

func (d *PetsDAL) Delete(ctx context.Context, petID string) error {
	_, err := d.Client.Collection("pets").Doc(petID).Delete(ctx)
	if err != nil {
		log.Fatalf("Could not delete pet: %v", err)
	}

	return nil
}

func Updates(update pets.Update) []firestore.Update {
	var updates []firestore.Update

	if update.Email != nil {
		updates = append(updates, firestore.Update{
			Path:  "email",
			Value: *update.Email,
		})
	}

	if update.Added != nil {
		updates = append(updates, firestore.Update{
			Path:  "added",
			Value: *update.Added,
		})
	}

	if update.Caption != nil {
		updates = append(updates, firestore.Update{
			Path:  "caption",
			Value: *update.Caption,
		})
	}

	if update.Image != nil {
		updates = append(updates, firestore.Update{
			Path:  "image",
			Value: *update.Image,
		})
	}

	if update.Likes != nil {
		updates = append(updates, firestore.Update{
			Path:  "likes",
			Value: update.Likes,
		})
	}

	if update.Owner != nil {
		updates = append(updates, firestore.Update{
			Path:  "owner",
			Value: update.Owner,
		})
	}

	if update.Petname != nil {
		updates = append(updates, firestore.Update{
			Path:  "petname",
			Value: update.Petname,
		})
	}

	return updates
}
