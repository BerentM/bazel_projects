package main

import (
	"testing"
)

func TestByteConversion(t *testing.T) {
	oh := ImageProcessor{objectPath: "images/test_image.png"}
	expected := 45190
	actual, err := oh.readFile()
	if err != nil {
		t.Fatal(err)
	}
	if expected != len(actual) {
		t.Fatalf("%q != %q", expected, actual)
	}
}

func TestUniqueID(t *testing.T) {
	oh := ImageProcessor{objectPath: "images/test_image.png"}
	bytes, err := oh.readFile()
	if err != nil {
		t.Fatal(err)
	}
	expected := "cfd47f07013ec51cc534771496e961c86a2e53d93cc3528ebbabde7630b07566"
	actual := oh.generateUniqueID(bytes)
	if expected != actual {
		t.Fatalf("%q != %q", expected, actual)
	}
}
