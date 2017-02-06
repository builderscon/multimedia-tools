package cover

import (
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
)

func readFontFile(s string) (*truetype.Font, error) {
	buf, err := ioutil.ReadFile(s)
	if err != nil {
		return nil, errors.Wrap(err, `failed to read file`)
	}

	f, err := freetype.ParseFont(buf)
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse font`)
	}

	return f, nil
}
