package cover

import (
	"image"
	"image/draw"
	"image/jpeg"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
)

type Text struct {
	Speaker    string
	Title      string
	Conference string
	Date       string
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

func LoadImage(path string) (*image.RGBA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, `failed to load image`)
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
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

func WriteText(base *image.RGBA, txt *Text) error {
	fc := freetype.NewContext()
	fc.SetDPI(300)
	fc.SetFont(typefont)
	fc.SetFontSize(21)
	fc.SetClip(base.Bounds())
	fc.SetDst(base)
	fc.SetSrc(image.White)

	pt := freetype.Pt(10, 100)

	fc.DrawString(txt.Title, pt)
	pt.Y += fc.PointToFixed(11 * 1.8)

	// For these, start from botoom
	bounds := base.Bounds()
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
