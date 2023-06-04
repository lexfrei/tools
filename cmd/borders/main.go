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

func main() {
	cmd.Execute()

	// Validate input. Additional border can't be less than 0, but can be more than 100.
	if cmd.AdditionalBorder < 0 {
		log.Fatal("additional border can't be less than 0")
	}

	var (
		borderClr = convertStringToColor(cmd.BorderColor)
		err       error
	)

	if cmd.InputFile != "" {
		err = AddBorderToImage(cmd.InputFile, cmd.Prefix, borderClr, cmd.AdditionalBorder)
	} else {
		files, err := getAllJPGFiles(cmd.Directory)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			// skip if file name starts with prefix
			if strings.HasPrefix(filepath.Base(file), cmd.Prefix) {
				continue
			}

			err = AddBorderToImage(file, cmd.Prefix, borderClr, cmd.AdditionalBorder)
			if err != nil {
				log.Printf("failed to add border to %s: %v", file, err)
			}
		}
	}

	if err != nil {
		log.Fatal(err)
	}
}

func AddBorderToImage(filePath, prefix string, borderColor color.Color, percent int) error {
	// Open the file.
	file, err := os.Open(filePath)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer file.Close()

	// Decode the JPEG data.
	img, err := jpeg.Decode(file)
	if err != nil {
		return errors.Wrap(err, "failed to decode JPEG")
	}

	// Calculate the new image size.
	newImage := generateImageWithBorder(img, borderColor, percent)

	// Create a new file name.
	newFileName, err := generateNewFileName(filePath, prefix)
	if err != nil {
		return errors.Wrap(err, "failed to generate new file name")
	}

	// Open a file for writing.
	outputFile, err := os.Create(newFileName)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer outputFile.Close()

	// Encode the new image to JPEG and write it to the file.
	err = jpeg.Encode(outputFile, newImage, &jpeg.Options{Quality: 100})
	if err != nil {
		return errors.Wrap(err, "failed to encode JPEG")
	}

	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// convertStringToColor converts a string to a color.Color.
// It returns white if the string is not recognized.
func convertStringToColor(colorString string) color.Color {
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

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".jpeg" {
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
func generateImageWithBorder(img image.Image, borderColor color.Color, paspartuPercent int) image.Image {
	// Create the new image with a border.
	borderSize := max(img.Bounds().Dx(), img.Bounds().Dy())
	extraSpace := borderSize * paspartuPercent / 100
	newSize := borderSize + 2*extraSpace
	newImage := image.NewRGBA(image.Rect(0, 0, newSize, newSize))

	// Fill the image with the border color.
	for x := 0; x < newSize; x++ {
		for y := 0; y < newSize; y++ {
			newImage.Set(x, y, borderColor)
		}
	}

	// Copy the old image into the center of the new one.
	rect := img.Bounds()
	rect = rect.Add(image.Pt((newSize-rect.Dx())/2, (newSize-rect.Dy())/2))
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

	newFileName := absParantDir + "/" + prefix + "_" + fileName + "_" + cmd.BorderColor + ".jpg"

	// If file with the same name exists, add a number to the end of the file name.
	if _, err := os.Stat(newFileName); err == nil {
		for i := 1; ; i++ {
			newFileName = absParantDir + "/" + prefix + "_" + fileName + "_" + cmd.BorderColor + "_" + strconv.Itoa(i) + ".jpg"
			if _, err := os.Stat(newFileName); os.IsNotExist(err) {
				break
			}
		}
	}

	return newFileName, nil
}