package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

// ImageProcessor object hashing
type ImageProcessor struct {
	data       []byte
	scale      float64
	objectPath string
	uniqueID   string
}

func (p *ImageProcessor) resize(grid [][]color.Color) [][]color.Color {
	xlen, ylen := int(float64(len(grid))*p.scale), int(float64(len(grid[0]))*p.scale)
	resized := make([][]color.Color, xlen)
	for i := 0; i < len(resized); i++ {
		resized[i] = make([]color.Color, ylen)
	}
	for x := 0; x < xlen; x++ {
		for y := 0; y < ylen; y++ {
			xp := int(math.Floor(float64(x) / p.scale))
			yp := int(math.Floor(float64(y) / p.scale))
			resized[x][y] = grid[xp][yp]
		}
	}
	return resized
}

func (p *ImageProcessor) gridToBytes(grid [][]color.Color) ([]byte, error) {
	buf := new(bytes.Buffer)
	// create an image and set the pixels using the grid
	xlen, ylen := len(grid), len(grid[0])
	rect := image.Rect(0, 0, xlen, ylen)
	img := image.NewNRGBA(rect)
	for x := 0; x < xlen; x++ {
		for y := 0; y < ylen; y++ {
			img.Set(x, y, grid[x][y])
		}
	}
	err := png.Encode(buf, img)
	return buf.Bytes(), err
}

// DownscaleImage reduce resolution of the Image
// based on https://go-recipes.dev/more-working-with-images-in-go-30b11ab2a9f0
func (p *ImageProcessor) DownscaleImage(data []byte) ([]byte, error) {
	var grid [][]color.Color
	img, err := png.Decode(bytes.NewReader(data))

	// create grid of pixels
	size := img.Bounds().Size()
	for i := 0; i < size.X; i++ {
		var y []color.Color
		for j := 0; j < size.Y; j++ {
			y = append(y, img.At(i, j))
		}
		grid = append(grid, y)
	}
	if err != nil {
		log.Fatal(err)
	}

	resized := p.resize(grid)
	bytes, err := p.gridToBytes(resized)
	p.data = bytes
	return bytes, err
}

func (p *ImageProcessor) readFile() ([]byte, error) {
	file, err := os.Open(p.objectPath)
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

func (p *ImageProcessor) generateUniqueID(bytes []byte) string {
	hash := sha256.New()
	hash.Write(bytes)
	uniqueID := hex.EncodeToString(hash.Sum(nil)[:])
	p.uniqueID = uniqueID
	return uniqueID
}
