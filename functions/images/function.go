package images

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
)

type Request struct {
	ImageURL string `json:"image_url"`
	PetID    string `json:"pet_id"`
}

func CollectImages(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")

	if projectID == "" {
		http.Error(w, "PROJECT ID was not set", 500)
		return
	}

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer client.Close()

	// Sets the name of the image file to annotate.

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	defer r.Body.Close()

	var req Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	uri := vision.NewImageFromURI(req.ImageURL)

	web, err := client.DetectWeb(ctx, uri, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	matches := append(web.FullMatchingImages, web.PartialMatchingImages...)
	matches = append(matches, web.VisuallySimilarImages...)

	var imagesURL []string

	for _, match := range matches {
		imagesURL = append(imagesURL, match.Url)
	}

	savedImages, err := saveToBucket(fmt.Sprintf("%s-images", projectID), imagesURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Could not create datastore client: %v", err)
	}

	defer func() {
		_ = firestoreClient.Close()
	}()
	_, err = firestoreClient.Collection("pets").Doc(req.PetID).Update(ctx, []firestore.Update{
		{
			Path:  "more_images",
			Value: savedImages,
		},
	})
	if err != nil {
		log.Fatalf("Could not update pet: %v", err)
	}

	marshal, err := json.Marshal(savedImages)
	_, err = io.Copy(w, bytes.NewReader(marshal))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func saveToBucket(bucketName string, images []string) ([]string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	// Connect to bucket
	bucket := client.Bucket(bucketName)
	// Get the url response

	var urls []string
	for _, image := range images {
		response, err := http.Get(image)
		if err != nil {
			fmt.Printf("ERROR %s", err)
			continue
		}

		if response.StatusCode == http.StatusOK {
			// Setup the GCS object with the filename to write to

			u, err := url.Parse(image)
			if err != nil {
				fmt.Printf("ERROR %s", err)
				response.Body.Close()
				continue
			}

			ext := filepath.Ext(u.Path)
			obj := bucket.Object(fmt.Sprintf("%s%s", uuid.NewString(), ext))

			// w implements io.Writer.
			w := obj.NewWriter(ctx)

			// Copy file into GCS
			if _, err := io.Copy(w, response.Body); err != nil {
				fmt.Printf("ERROR %s", err)
				response.Body.Close()
				continue
			}

			// Close, just like writing a file. File appears in GCS after
			if err := w.Close(); err != nil {
				fmt.Printf("ERROR %s", err)
				response.Body.Close()
				continue
			}

			urls = append(urls, fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, obj.ObjectName()))
		}

		response.Body.Close()
	}

	return urls, nil
}
