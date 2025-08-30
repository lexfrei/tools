package main

import (
	"fmt"
	"log"
	"os"

	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: exiftool <image_file>")
	}

	rawExif, err := exif.SearchFileAndExtractExif(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	im, err := exifcommon.NewIfdMappingWithStandard()
	if err != nil {
		log.Fatal(err)
	}

	ti := exif.NewTagIndex()

	_, index, err := exif.Collect(im, ti, rawExif)
	if err != nil {
		log.Fatal(err)
	}

	// Получаем корневой IFD
	rootIfd := index.RootIfd

	// Извлекаем нужные теги
	tags := make(map[string]string)

	err = rootIfd.EnumerateTagsRecursively(func(ifd *exif.Ifd, ite *exif.IfdTagEntry) error {
		tagName := ite.TagName()
		value, err := ite.Format()
		if err == nil {
			tags[tagName] = value
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Форматируем вывод
	fmt.Printf("📷 %s %s\n", tags["Make"], tags["Model"])
	if lens, ok := tags["LensModel"]; ok {
		fmt.Printf("🔭 %s\n", lens)
	}

	fmt.Printf("\n%s | %s | %s | ISO %s\n",
		tags["FocalLength"],
		tags["FNumber"],
		tags["ExposureTime"],
		tags["ISOSpeedRatings"])
}
