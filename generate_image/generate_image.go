package generate_image

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

const (
	width  = 400
	height = 400
)

func findLongestElement(strList []string, getLengthFunc func(line string) int) string {
	if len(strList) == 0 {
		return ""
	}

	longest := strList[0]

	for _, str := range strList {
		if getLengthFunc(str) > getLengthFunc(longest) {
			longest = str
		}
	}

	return longest
}

func getLargestPossibleFontSize(drawCtx font.Drawer, lines []string, fontVar *sfnt.Font) (int, error) {
	size := 20
	longestLine := findLongestElement(lines, func(line string) int {
		return drawCtx.MeasureString(line).Ceil()
	})
	for {
		face, err := opentype.NewFace(fontVar, &opentype.FaceOptions{
			Size: float64(size),
			DPI:  72,
		})
		if err != nil {
			return 0, errors.New("unable to get face")
		}
		drawCtx.Face = face
		textWidth := drawCtx.MeasureString(longestLine).Ceil()

		// Check if text is wider than the image
		if textWidth > width {
			return size - 10, nil
		}

		// Check if text is taller than the image
		if face.Metrics().Height.Ceil()*len(lines) > height {
			return size - 20, nil
		}

		size += 10
	}
}

func GenerateImage(text string, fontVar *sfnt.Font) error {
	textForFilename := text
	textForFilename = strings.ReplaceAll(textForFilename, "\n", "-")
	textForFilename = strings.ReplaceAll(textForFilename, " ", "-")
	fileName := fmt.Sprintf("img/%s.png", textForFilename)
	outputFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	face, err := opentype.NewFace(fontVar, &opentype.FaceOptions{
		Size: 72,
		DPI:  72,
	})
	if err != nil {
		return err
	}

	text = strings.ReplaceAll(text, " ", "\n")
	lines := strings.Split(text, "\n")

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	drawContext := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.White),
		Face: face,
	}

	maxFontSize, err := getLargestPossibleFontSize(*drawContext, lines, fontVar)
	if err != nil {
		return err
	}

	face, err = opentype.NewFace(fontVar, &opentype.FaceOptions{
		Size: float64(maxFontSize),
		DPI:  72,
	})
	if err != nil {
		return err
	}

	drawContext.Face = face

	y := (height - (len(lines)/2)*int(face.Metrics().Height.Ceil())) / 2

	for _, line := range lines {
		x := (width - drawContext.MeasureString(line).Ceil()) / 2
		drawContext.Dot = fixed.P(x, y)
		drawContext.DrawString(line)
		y += int(face.Metrics().Height.Ceil())
	}

	err = png.Encode(outputFile, drawContext.Dst)
	if err != nil {
		return err
	}

	return nil
}
