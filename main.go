package main

import (
	"fmt"
	"github.com/rawnly/collections-downloader/unsplash"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func handleError(e error) {
	panic(e.Error())
}


func main() {
	const workers = 5

	var pageSize float64 = 5
	var photos []unsplash.Photo

	photosLimit := 15
	numberOfPages := 0
	startTime := time.Now()
	args := os.Args[1:len(os.Args)]

	if len(args) <= 0 {
		return
	} else if len(args) == 2 {
		photosLimit, _ = strconv.Atoi(args[1])
	}


	collectionID := args[0]
	collection, err := unsplash.GetCollection(collectionID)

	if err != nil {
		handleError(err)
	}

	fmt.Printf("\n\nCollection: %s\n", collection.Title)
	fmt.Printf("Downloading %d photos of %d\n\n", photosLimit, collection.TotalPhotos)

	numberOfPages = int(math.Ceil(float64(photosLimit) / pageSize))

	jobs := make(chan int, numberOfPages)
	results := make(chan int, collection.TotalPhotos)

	// Get all the photos
	for i := 1; i <= numberOfPages; i++ {
		page := i

		p, err := unsplash.GetCollectionPhotos(*collection, page, int(pageSize))

		if err != nil {
			handleError(err)
		}

		photos = append(photos, p...)
	}


	// Download the photos
	for w := 1; w <= workers; w++ {
		go func() {
			for idx := range jobs {
				photo := photos[idx]

				if err := unsplash.DownloadPhoto(photo, strings.Replace(collection.Title, " ", "_", strings.Count(collection.Title, " "))); err != nil {
					handleError(err)
				}

				results <- idx + 1
			}
		}()
	}

	for j := 1; j < len(photos); j++ {
		jobs <- j
	}

	close(jobs)

	for a:=1; a< len(photos); a++ {
		<-results
	}

	// Calculate elapsed time
	endTime := time.Now()
	elapsed := endTime.Sub(startTime)

	fmt.Println("Job completed in ", elapsed)
}
