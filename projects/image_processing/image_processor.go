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
	img       []byte
	thumbnail []byte
	scale     float64
	uniqueID  string
}

// NewImageProcessor Initialize new image processor
func NewImageProcessor(scale float64) *ImageProcessor {
	return &ImageProcessor{
		scale: scale,
	}
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

func (p *ImageProcessor) miniaturization(maxLength float64, grid [][]color.Color) [][]color.Color {
	scale := maxLength / math.Max(float64(len(grid)), float64(len(grid[0])))
	xlen, ylen := int(float64(len(grid))*scale), int(float64(len(grid[0]))*scale)
	resized := make([][]color.Color, xlen)
	for i := 0; i < len(resized); i++ {
		resized[i] = make([]color.Color, ylen)
	}
	for x := 0; x < xlen; x++ {
		for y := 0; y < ylen; y++ {
			xp := int(math.Floor(float64(x) / scale))
			yp := int(math.Floor(float64(y) / scale))
			resized[x][y] = grid[xp][yp]
		}
	}
	fmt.Println("xlen: ", xlen, " ylen: ", ylen)
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

// DownscaleImage reduce resolution of the image
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

// GenerateThumbnail create thumbnail from image
func (p *ImageProcessor) GenerateThumbnail(data []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	grid := p.imageToGrid(img)
	resized := p.miniaturization(128, grid)
	bytes, err := p.gridToBytes(resized)
	p.thumbnail = bytes
	return bytes, err
}

// Process downscale and generateUniqueID for passed image
func (p *ImageProcessor) Process(img []byte) {
	p.scale = p.calculateDownscaleRate(img, 2000000) // 2MB limit
	_, err := p.GenerateThumbnail(img)
	if err != nil {
		log.Fatal(err)
	}
	_, err = p.DownscaleImage(img)
	if err != nil {
		log.Fatal(err)
	}
	p.generateUniqueID(p.img)
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

// calculateDownscaleRate calculate how much []byte exceed byteSizeLimit
func (p *ImageProcessor) calculateDownscaleRate(bytes []byte, byteSizeLimit float64) float64 {
	c := float64(cap(bytes))
	if c <= byteSizeLimit {
		// Don't downscale image if its size is smaller than limit
		return float64(1)
	}
	scale := 1 / (c / byteSizeLimit)
	fmt.Printf("%.8f\n", scale)
	return scale
}
