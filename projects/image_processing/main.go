package main

import (
	"log"
)

func main() {
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
		// out, _ := os.Create("./images/output.png")
		// defer out.Close()
		// img, err := png.Decode(bytes.NewReader(imgProcessor.img))
		// err = png.Encode(out, img)
		// if err != nil {
		// 	log.Println(err)
		// }
	}
}
