// Image thumbnail generator
package images

import (
	"fmt"
	"os"
	. "pkg.deepin.io/dde/api/thumbnails/loader"
	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	ImageTypePng  string = "image/png"
	ImageTypeJpeg        = "image/jpeg"
	ImageTypeGif         = "image/gif"
	ImageTypeBmp         = "image/bmp"
	ImageTypeTiff        = "image/tiff"
	ImageTypeSvg         = "image/svg+xml"
)

func init() {
	for _, ty := range SupportedTypes() {
		switch ty {
		case ImageTypeSvg:
			Register(ty, genSvgThumbnail)
		default:
			Register(ty, genImageThumbnail)
		}
	}
}

func SupportedTypes() []string {
	return []string{
		ImageTypePng,
		ImageTypeJpeg,
		ImageTypeGif,
		ImageTypeBmp,
		ImageTypeTiff,
		ImageTypeSvg,
	}
}

func GenThumbnail(src string, width, height int, force bool) (string, error) {
	if width <= 0 || height <= 0 {
		return "", fmt.Errorf("Invalid width or height")
	}

	ty, err := mime.Query(src)
	if err != nil {
		return "", err
	}

	if !IsStrInList(ty, SupportedTypes()) {
		return "", fmt.Errorf("No supported type: %v", ty)
	}

	switch ty {
	case ImageTypeSvg:
		return genSvgThumbnail(src, "", width, height, force)
	}

	return genImageThumbnail(src, "", width, height, force)
}

func genSvgThumbnail(src, bg string, width, height int, force bool) (string, error) {
	tmp := GetTmpImage()
	err := svgToPng(src, tmp)
	if err != nil {
		return "", err
	}

	defer os.Remove(tmp)
	return genImageThumbnail(tmp, bg, width, height, force)
}

func genImageThumbnail(src, bg string, width, height int, force bool) (string, error) {
	src = dutils.DecodeURI(src)
	dest, err := GetThumbnailDest(src, width, height)
	if err != nil {
		return "", err
	}
	if !force && dutils.IsFileExist(dest) {
		return dest, nil
	}

	err = ThumbnailImage(src, dest, width, height)
	if err != nil {
		return "", err
	}
	return dest, nil
}