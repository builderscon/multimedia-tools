package cover

import (
	"context"
	"encoding/csv"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
)

type Text struct {
	Speaker       string
	Title         string
	TitleFontSize float64
	Conference    string
	Date          string
}

const fontFile = "/Users/daisuke/Library/Fonts/migmix-1m-bold.ttf"

var typefont *truetype.Font

func init() {
	var err error
	typefont, err = readFontFile(fontFile)
	if err != nil {
		panic(err.Error())
	}
}

func LoadImage(path string) (draw.Image, error) {
	var rdr io.Reader
	if strings.HasPrefix(path, "http") {
		res, err := http.Get(path)
		if err != nil {
			return nil, errors.Wrap(err, `failed to fetch image`)
		}
		defer res.Body.Close()

		rdr = res.Body
	} else {

		f, err := os.Open(path)
		if err != nil {
			return nil, errors.Wrap(err, `failed to load image`)
		}
		defer f.Close()
		rdr = f
	}

	img, err := jpeg.Decode(rdr)
	if err != nil {
		return nil, errors.Wrap(err, `failed to decode image`)
	}

	bounds := img.Bounds()
	if bounds.Dx() > 1280 || bounds.Dy() > 720 {
		img = imaging.Fill(img, 1280, 720, imaging.Center, imaging.Lanczos)
		bounds = image.Rect(0, 0, 1280, 720)
	}

	rgba := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.DrawMask(rgba, bounds, img, image.ZP, &gradmask{w: bounds.Dx(), h: bounds.Dy()}, image.ZP, draw.Src)
	return rgba, nil
}

func GetThumbnail(id string) (*image.RGBA, error) {
	u := `https://i.ytimg.com/vi/` + id + `/maxresdefault.jpg`
	res, err := http.Get(u)
	if err != nil {
		return nil, errors.Wrap(err, `failed to fetch thumbnail`)
	}

	defer res.Body.Close()
	img, err := jpeg.Decode(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, `failed to decode image`)
	}
	bounds := img.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.DrawMask(rgba, bounds, img, image.ZP, &gradmask{w: bounds.Dx(), h: bounds.Dy()}, image.ZP, draw.Src)
	return rgba, nil
}

func WriteText(base draw.Image, txt *Text) error {
	bounds := base.Bounds()
	fc := freetype.NewContext()
	fc.SetDPI(300)
	fc.SetFont(typefont)
	fc.SetClip(bounds)
	fc.SetDst(base)
	fc.SetSrc(image.White)

	pt := freetype.Pt(10, 100)

	fc.SetFontSize(txt.TitleFontSize)
	for _, l := range strings.Split(txt.Title, "\n") {
		fc.DrawString(l, pt)
		pt.Y += fc.PointToFixed(11 * 1.8)
	}

	// For these, start from botoom
	pt = freetype.Pt(10, bounds.Dy())
	pt.Y -= fc.PointToFixed(6)

	// the very bottom
	fc.SetFontSize(6)
	fc.DrawString(txt.Date, pt)

	pt.Y -= fc.PointToFixed(6)
	fc.SetFontSize(6)
	fc.DrawString(txt.Conference, pt)
	pt.Y -= fc.PointToFixed(10)
	fc.SetFontSize(10)
	fc.DrawString(txt.Speaker, pt)
	return nil
}

func Run(ctx context.Context, src io.Reader) error {
	rdr := csv.NewReader(src)
	rdr.Comma = '\t'

	for {
		record, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, `failed to read TSV`)
		}

		if len(record[0]) == 0 {
			continue
		}

		img, err := LoadImage(record[0])
		if err != nil {
			return errors.Wrap(err, `failed to get image`)
		}

		// left 40% of the image is going to be used for our text
		textarea := image.Rect(0, 0, int(float64(img.Bounds().Dx())*0.4), img.Bounds().Dy())

		// Create a new canvas to write the text
		canvas := image.NewAlpha(textarea)

		var txt Text
		txt.Title = strings.Replace(record[1], "\\n", "\n", -1)
		log.Printf("%s", txt.Title)
		ftsize, err := strconv.Atoi(record[2])
		if err != nil {
			return errors.Wrap(err, `failed to parse font size`)
		}
		txt.TitleFontSize = float64(ftsize)
		txt.Conference = record[3]
		txt.Date = record[4]
		txt.Speaker = record[5]

		if err := WriteText(canvas, &txt); err != nil {
			return errors.Wrap(err, `failed to write text`)
		}

		// Draw the canvas to the image
		fitimg := imaging.Fit(canvas, textarea.Dx(), textarea.Dy(), imaging.Lanczos)
		draw.DrawMask(img, textarea, fitimg, image.ZP, nil, image.ZP, draw.Over)

		f, err := os.Create(record[6])
		if err != nil {
			return errors.Wrap(err, `failed to create file`)
		}
		defer f.Close()
		if err := jpeg.Encode(f, img, nil); err != nil {
			return errors.Wrap(err, `failed to encode jpeg`)
		}
	}
	return nil
}
