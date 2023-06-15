package main

import (
	"bytes"
	"image/color"
	"image/png"
	"testing"
)

func TestFileRead(t *testing.T) {
	p := ImageProcessor{}
	expected := 45190
	actual, err := p.readFile("images/test_image.png")
	if err != nil {
		t.Fatal(err)
	}
	if expected != len(actual) {
		t.Fatalf("%v != %v", expected, actual)
	}
}

func TestUniqueID(t *testing.T) {
	p := ImageProcessor{}
	bytes, err := p.readFile("images/test_image.png")
	if err != nil {
		t.Fatal(err)
	}
	expected := "cfd47f07013ec51cc534771496e961c86a2e53d93cc3528ebbabde7630b07566"
	actual := p.generateUniqueID(bytes)
	if expected != actual {
		t.Fatalf("%v != %v", expected, actual)
	}
}

func TestImgToGrid(t *testing.T) {
	p := ImageProcessor{}
	data, err := p.readFile("images/test_image.png")
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	expected := 652
	actual := len(p.imageToGrid(img))
	if expected != actual {
		t.Fatalf("expected len: %v != actual len: %v", expected, actual)
	}
	expectedColor := color.Color(color.RGBA{199, 198, 202, 255})
	e1, e2, e3, e4 := expectedColor.RGBA()
	actualColor := p.imageToGrid(img)[300][100]
	a1, a2, a3, a4 := actualColor.RGBA()

	if e1 != a1 || e2 != a2 || e3 != a3 || e4 != a4 {
		t.Fatalf("expected color: %v != actual color: %v", expectedColor, actualColor)
	}
}

func TestDownscale(t *testing.T) {
	p := ImageProcessor{scale: 0.9}
	expected := 25330
	input, err := p.readFile("images/test_image.png")
	if err != nil {
		t.Fatal(err)
	}
	actual, err := p.DownscaleImage(input)
	if err != nil {
		t.Fatal(err)
	}
	if expected != len(actual) {
		t.Fatalf("%v != %v", expected, len(actual))
	}
}