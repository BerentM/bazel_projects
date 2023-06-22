package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

func greet(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func processLocalFiles(w http.ResponseWriter, _ *http.Request) {
	var images [][]byte

	paths := []string{"projects/image_processing/images/test_image.png"}
	imgProcessor := NewImageProcessor(1)
	S3Client := NewS3Client("dtmx-images")

	for _, path := range paths {
		data, err := imgProcessor.readFile(path)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), 400)
			return
		}
		images = append(images, data)
	}

	for _, img := range images {
		imgProcessor.Process(img)
		S3Client.Upload(imgProcessor.img, imgProcessor.uniqueID)
	}
	fmt.Fprintf(w, "Last processed img UniqueID: %v", imgProcessor.uniqueID)
}

func loadImage(_ http.ResponseWriter, r *http.Request) (multipart.File, error) {
	if r.Method == "POST" {
		file, _, err := r.FormFile("file")
		defer file.Close()
		return file, err
	}
	return nil, fmt.Errorf("Method not allowed")
}

// processHTTPFile it requires that request Method is POST and file key is "file"
func processHTTPFile(w http.ResponseWriter, r *http.Request) {
	imgProcessor := NewImageProcessor(1)
	S3Client := NewS3Client("dtmx-images")

	file, err := loadImage(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 400)
		return
	}
	img, _ := io.ReadAll(file)
	imgProcessor.Process(img)
	S3Client.Upload(imgProcessor.img, imgProcessor.uniqueID)
	fmt.Fprintf(w, "Last processed img UniqueID: %v", imgProcessor.uniqueID)
}

func main() {
	http.HandleFunc("/", greet)
	http.HandleFunc("/local", processLocalFiles)
	http.HandleFunc("/image", processHTTPFile)
	http.ListenAndServe(":8080", nil)
}
