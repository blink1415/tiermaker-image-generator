package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/image/font/opentype"
	"github.com/blink1415/tiermaker-image-generator/generate_image"
)

func main() {
	fontBytes, err := os.ReadFile("font.ttf")
	if err != nil {
		fmt.Println("Error loading font:", err)
		return
	}

	fontVar, err := opentype.Parse(fontBytes)
	if err != nil {
		fmt.Println("Error parsing font:", err)
		return
	}

	wordsFile, err := os.ReadFile("text.txt")
	if err != nil {
		fmt.Println("Error reading words file:", err)
		return
	}

	wordList := strings.Split(string(wordsFile), "\n")

	if _, err := os.Stat("img"); os.IsNotExist(err) {
		os.Mkdir("img", os.ModePerm)
	}

	for _, word := range wordList {
		if word == "" {
			continue
		}

		if err = generate_image.GenerateImage(word, fontVar); err != nil {
			fmt.Printf("Error generating image: %s", err)
		}
	}
}

