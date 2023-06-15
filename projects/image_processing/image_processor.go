package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

// ImageProcessor object hashing
type ImageProcessor struct {
	img      []byte
	scale    float64
	uniqueID string
}

// New Initialize new image processor
func (p *ImageProcessor) New(scale float64) {
	p.scale = scale
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

// bytesToGrid create grid of pixels
func (p *ImageProcessor) imageToGrid(img image.Image) [][]color.Color {
	var grid [][]color.Color
	size := img.Bounds().Size()
	for i := 0; i < size.X; i++ {
		var y []color.Color
		for j := 0; j < size.Y; j++ {
			y = append(y, img.At(i, j))
		}
		grid = append(grid, y)
	}
	return grid
}

// DownscaleImage reduce resolution of the Image
// based on https://go-recipes.dev/more-working-with-images-in-go-30b11ab2a9f0
func (p *ImageProcessor) DownscaleImage(data []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	grid := p.imageToGrid(img)
	resized := p.resize(grid)
	bytes, err := p.gridToBytes(resized)
	p.img = bytes
	return bytes, err
}

// Process downscale and generateUniqueID for passed image
func (p *ImageProcessor) Process(img []byte) {
	data, err := p.DownscaleImage(img)
	if err != nil {
		log.Fatal(err)
	}
	p.generateUniqueID(data)
	p.img = data
	fmt.Println(p.uniqueID)
}

func (p *ImageProcessor) readFile(path string) ([]byte, error) {
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

func (p *ImageProcessor) generateUniqueID(bytes []byte) string {
	hash := sha256.New()
	hash.Write(bytes)
	uniqueID := hex.EncodeToString(hash.Sum(nil)[:])
	p.uniqueID = uniqueID
	return uniqueID
}
