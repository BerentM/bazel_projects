package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func greet(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func processLocalFiles(w http.ResponseWriter, _ *http.Request) {
	var images [][]byte

	aws := AwsHelper{}
	aws.New()
	aws.CheckBuckets()

	paths := []string{"projects/image_processing/images/test_image.png"}
	imgProcessor := ImageProcessor{scale: 0.9}

	for _, path := range paths {
		data, err := imgProcessor.readFile(path)
		if err != nil {
			log.Fatal(err)
		}
		images = append(images, data)
	}

	for _, img := range images {
		imgProcessor.Process(img)
		aws.Upload(imgProcessor.img, imgProcessor.uniqueID)
	}
	fmt.Fprintf(w, "Last processed img UniqueID: %v", imgProcessor.uniqueID)
}

func main() {
	http.HandleFunc("/", greet)
	http.HandleFunc("/local", processLocalFiles)
	http.ListenAndServe(":8080", nil)
}
