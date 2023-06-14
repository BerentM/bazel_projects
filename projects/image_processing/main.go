package main

import (
	"fmt"
	"log"
)

func main() {
	imgProcessor := ImageProcessor{objectPath: "images/test_image.png", scale: 0.9}
	data, err := imgProcessor.readFile()
	if err != nil {
		log.Fatal(err)
	}
	data, err = imgProcessor.DownscaleImage(data)
	if err != nil {
		log.Fatal(err)
	}
	imgProcessor.generateUniqueID(data)
	fmt.Println(imgProcessor.uniqueID)

	// out, _ := os.Create("./images/output.png")
	// defer out.Close()
	// img, err := png.Decode(bytes.NewReader(data))
	// err = png.Encode(out, img)
	// if err != nil {
	// 	log.Println(err)
	// }

	aws := AwsHelper{}
	aws.New()
	aws.CheckBuckets()
	aws.Upload(data, imgProcessor.uniqueID)
	// aws.Download(imgProcessor.uniqueID)
}
