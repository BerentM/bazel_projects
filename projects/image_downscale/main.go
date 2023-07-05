package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
)

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	size := stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err
}

func main() {
	inputFilePath := "images/test_image.png"
	outputFilePath := "images/output.png"
	vips.LoggingSettings(nil, vips.LogLevelWarning)

	bytes, err := readFile(inputFilePath)

	image, err := vips.NewImageFromBuffer(bytes)
	if err != nil {
		panic(err)
	}
	if err := image.Thumbnail(128, 128, vips.InterestingNone); err != nil {
		panic(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	ep := vips.NewDefaultPNGExportParams()
	image1bytes, _, err := image.Export(ep)
	err = ioutil.WriteFile(outputFilePath, image1bytes, 0644)
}
