package main

import (
	"image/jpeg"
	"log"
	"os"

	cover "github.com/builderscon/multimedia-tools/video-cover"
	"github.com/pkg/errors"
)

func main() {
	if err := _main(); err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
}

func _main() error {
	img, err := cover.LoadImage(os.Args[1])
	if err != nil {
		return errors.Wrap(err, `failed to get image`)
	}

	var txt cover.Text
	txt.Title = "Opening"
	txt.Conference = "builderscon tokyo 2016"
	txt.Date = "Dec 3, 2016"
	txt.Speaker = "Daisuke Maki"

	if err := cover.WriteText(img, &txt); err != nil {
		return errors.Wrap(err, `failed to write text`)
	}

	f, err := os.Create(`out.jpeg`)
	if err != nil {
		return errors.Wrap(err, `failed to create file`)
	}
	defer f.Close()
	return jpeg.Encode(f, img, nil)
}
