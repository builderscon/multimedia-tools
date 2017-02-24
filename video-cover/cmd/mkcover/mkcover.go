package main

import (
	"context"
	"io"
	"log"
	"os"

	cover "github.com/builderscon/multimedia-tools/video-cover"
	tty "github.com/builderscon/multimedia-tools/video-cover/internal/tty"
	"github.com/pkg/errors"
)

func main() {
	if err := _main(); err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
}

func _main() error {
	var src io.Reader
	if len(os.Args) > 1 {
log.Printf("Using file %s", os.Args[1])
		f, err := os.Open(os.Args[1])
		if err != nil {
			return errors.Wrapf(err, `failed to open file for reading: %s`, os.Args[1])
		}
		defer f.Close()
		src = f
	} else if tty.IsTty(os.Stdin) {
log.Printf("Using os.Stdin")
		src = os.Stdin
	} else {
		return errors.New("No input data available")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return cover.Run(ctx, src)
}
