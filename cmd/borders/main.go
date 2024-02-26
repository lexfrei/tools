/*
Copyright © 2023 Aleksei Sviridkin <f@lex.la>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

 1. Redistributions of source code must retain the above copyright notice,
    this list of conditions and the following disclaimer.

 2. Redistributions in binary form must reproduce the above copyright notice,
    this list of conditions and the following disclaimer in the documentation
    and/or other materials provided with the distribution.

 3. Neither the name of the copyright holder nor the names of its contributors
    may be used to endorse or promote products derived from this software
    without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/lexfrei/tools/cmd/borders/cmd"
)

const (
	jpgExtension = ".jpg"
)

//nolint:funlen,gocyclo,cyclop // it's main func
func main() {
	cmd.Execute()

	if cmd.AdditionalBorder < 0 {
		log.Fatal("additional border can't be less than 0")
	}

	var err error

	//nolint:nestif // it's ok
	if cmd.InputFile != "" {
		img, err := openImageFile(cmd.InputFile)
		if err != nil {
			log.Fatalf("can't open image file: %v", err)
		}

		borderColor := convertStringToColor(cmd.BorderColor, img)

		var proportions float32
		if cmd.Minimal {
			proportions = generateOptimalProporions(img)
		} else {
			proportions = 1
		}

		newImage := generateImageWithBorder(img, borderColor, cmd.AdditionalBorder, proportions)

		err = saveImageFile(newImage, cmd.InputFile, cmd.Prefix)
		if err != nil {
			log.Fatalf("can't save image file: %v", err)
		}
	} else {
		files, err := getAllJPGFiles(cmd.Directory)
		if err != nil {
			log.Fatal(err)
		}

		for range files {
			img, err := openImageFile(cmd.InputFile)
			if err != nil {
				log.Fatalf("can't open image file: %v", err)
			}

			borderColor := convertStringToColor(cmd.BorderColor, img)

			var proportions float32
			if cmd.Minimal {
				proportions = generateOptimalProporions(img)
			} else {
				proportions = 1
			}

			newImage := generateImageWithBorder(img, borderColor, cmd.AdditionalBorder, proportions)

			err = saveImageFile(newImage, cmd.InputFile, cmd.Prefix)
			if err != nil {
				log.Fatalf("can't save image file: %v", err)
			}
		}
	}

	if err != nil {
		log.Fatal(err)
	}
}

// convertStringToColor converts a string to a color.Color.
// It returns white if the string is not recognized.
func convertStringToColor(colorString string, img image.Image) color.Color {
	if colorString == "avg" {
		return getAverageColor(img)
	}

	colorMap := map[string]color.RGBA{
		"white":     {255, 255, 255, 255},
		"black":     {0, 0, 0, 255},
		"red":       {255, 0, 0, 255},
		"green":     {0, 255, 0, 255},
		"blue":      {0, 0, 255, 255},
		"yellow":    {255, 255, 0, 255},
		"cyan":      {0, 255, 255, 255},
		"magenta":   {255, 0, 255, 255},
		"gray":      {128, 128, 128, 255},
		"darkgray":  {169, 169, 169, 255},
		"lightgray": {211, 211, 211, 255},
		"orange":    {255, 165, 0, 255},
		"pink":      {255, 192, 203, 255},
		"purple":    {128, 0, 128, 255},
		"violet":    {238, 130, 238, 255},
		"brown":     {165, 42, 42, 255},
	}

	// Check if the color string is a hex code
	if strings.HasPrefix(colorString, "#") {
		// Parse the hex code to an integer
		rgbInt, err := strconv.ParseInt(strings.TrimPrefix(colorString, "#"), 16, 32)
		if err != nil {
			return color.RGBA{255, 255, 255, 255}
		}

		// Extract the red, green, and blue components from the integer
		// and return the color
		//nolint:gomnd // 0xFF is a magic number and it's ok
		return color.RGBA{uint8(rgbInt >> 16), uint8((rgbInt >> 8) & 0xFF), uint8(rgbInt & 0xFF), 255}
	}

	// Check if the color string is a named color
	if namedColor, ok := colorMap[colorString]; ok {
		return namedColor
	}

	return colorMap["white"]
}

// getAllJPGFiles returns all jpg files in a directory.
func getAllJPGFiles(dirPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(dirPath, func(path string, _ os.FileInfo, _ error) error {
		if filepath.Ext(path) == jpgExtension || filepath.Ext(path) == ".jpeg" {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to walk directory")
	}

	return files, nil
}

// generateImageWithBorder generates a new image with a border.
func generateImageWithBorder(
	img image.Image, borderColor color.Color, paspartuPercent int, targetProportions float32,
) image.Image {
	var width, height int

	if img.Bounds().Dx() < img.Bounds().Dy() {
		height = img.Bounds().Dy() * (1 + 2*paspartuPercent/100)
		width = int(float32(height) * targetProportions)
	} else {
		width = img.Bounds().Dx() * (1 + 2*paspartuPercent/100)
		height = int(float32(width) * targetProportions)
	}

	newImage := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill the image with the border color.
	draw.Draw(newImage, newImage.Bounds(), &image.Uniform{borderColor}, image.Point{}, draw.Src)

	// Copy the old image into the center of the new one.
	rect := img.Bounds()
	rect.Min.X = (width - img.Bounds().Dx()) / 2
	rect.Min.Y = (height - img.Bounds().Dy()) / 2
	rect.Max.X = rect.Min.X + img.Bounds().Dx()
	rect.Max.Y = rect.Min.Y + img.Bounds().Dy()
	draw.Draw(newImage, rect, img, image.Point{}, draw.Over)

	return newImage
}

// generateNewFileName generates a new file name for the new image.
func generateNewFileName(filePath, prefix string) (string, error) {
	// Generate file name for the new image.
	absParantDir, err := filepath.Abs(filepath.Dir(filePath))
	if err != nil {
		return "", errors.Wrap(err, "failed to get absolute path")
	}

	fileName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	newFileName := absParantDir + "/" + prefix + "_" + fileName + "_" + cmd.BorderColor + jpgExtension

	// If file with the same name exists, add a number to the end of the file name.
	if _, err := os.Stat(newFileName); err == nil {
		for i := 1; ; i++ {
			newFileName = absParantDir + "/" +
				prefix + "_" + fileName + "_" + cmd.BorderColor + "_" + strconv.Itoa(i) + jpgExtension
			if _, err := os.Stat(newFileName); os.IsNotExist(err) {
				break
			}
		}
	}

	return newFileName, nil
}

// getAverageColor calculates the average color of an image.
//
//nolint:gomnd // 8 is a magic number and it's ok
func getAverageColor(img image.Image) color.Color {
	bounds := img.Bounds()

	// Define variables to store the sum of color channel values
	var sumR, sumG, sumB uint32

	// Iterate over each pixel in the top-left corner of the image
	for y := bounds.Min.Y; y < bounds.Min.Y+10; y++ {
		for x := bounds.Min.X; x < bounds.Min.X+10; x++ {
			pixel := img.At(x, y)
			r, g, b, _ := pixel.RGBA()

			// Add color channel values to the sum
			sumR += r
			sumG += g
			sumB += b
		}
	}

	// Calculate the average color channel values
	pixelCount := 10 * 10
	avgR := uint8(sumR / uint32(pixelCount) >> 8)
	avgG := uint8(sumG / uint32(pixelCount) >> 8)
	avgB := uint8(sumB / uint32(pixelCount) >> 8)

	return color.RGBA{avgR, avgG, avgB, 255}
}

func openImageFile(filePath string) (image.Image, error) {
	// Open the image file
	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	defer imgFile.Close()

	// Decode the image file
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode image")
	}

	return img, nil
}

func saveImageFile(img image.Image, filePath, prefix string) error {
	// Create a new file for the new image
	newFileName, err := generateNewFileName(filePath, prefix)
	if err != nil {
		return errors.Wrap(err, "failed to generate new file name")
	}

	outputFile, err := os.Create(newFileName)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer outputFile.Close()

	// Encode the new image to JPEG and write it to the file.
	err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: 100})
	if err != nil {
		return errors.Wrap(err, "failed to encode image")
	}

	return nil
}

// generateOptimalProporions generates optimal proportions for the new image.
// https://www.adobe.com/express/discover/sizes/instagram
//
//nolint:gomnd // magic proportions for instagram
func generateOptimalProporions(img image.Image) float32 {
	// if vertical image
	if img.Bounds().Dx() < img.Bounds().Dy() {
		return float32(0.8)
	}

	return 1.91
}
