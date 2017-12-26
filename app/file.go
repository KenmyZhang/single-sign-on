package app

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	l4g "github.com/alecthomas/log4go"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	_ "golang.org/x/image/bmp"

	"github.com/KenmyZhang/single-sign-on/model"
	"github.com/KenmyZhang/single-sign-on/utils"
)

const (
	/*
	  EXIF Image Orientations
	  1        2       3      4         5            6           7          8

	  888888  888888      88  88      8888888888  88                  88  8888888888
	  88          88      88  88      88  88      88  88          88  88      88  88
	  8888      8888    8888  8888    88          8888888888  8888888888          88
	  88          88      88  88
	  88          88  888888  888888
	*/
	Upright            = 1
	UprightMirrored    = 2
	UpsideDown         = 3
	UpsideDownMirrored = 4
	RotatedCWMirrored  = 5
	RotatedCCW         = 6
	RotatedCCWMirrored = 7
	RotatedCW          = 8

	MaxImageSize                 = 6048 * 4032 // 24 megapixels, roughly 36MB as a raw image
	IMAGE_THUMBNAIL_PIXEL_WIDTH  = 120
	IMAGE_THUMBNAIL_PIXEL_HEIGHT = 100
	IMAGE_PREVIEW_PIXEL_WIDTH    = 1024
)

func prepareImage(fileData []byte) (*image.Image, int, int) {
	// Decode image bytes into Image object
	img, imgType, err := image.Decode(bytes.NewReader(fileData))
	if err != nil {
		l4g.Error(utils.T("api.file.handle_images_forget.decode.error"), err)
		return nil, 0, 0
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// Fill in the background of a potentially-transparent png file as white
	if imgType == "png" {
		dst := image.NewRGBA(img.Bounds())
		draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
		draw.Draw(dst, dst.Bounds(), img, img.Bounds().Min, draw.Over)
		img = dst
	}

	// Flip the image to be upright
	orientation, _ := getImageOrientation(bytes.NewReader(fileData))
	img = makeImageUpright(img, orientation)

	return &img, width, height
}

func makeImageUpright(img image.Image, orientation int) image.Image {
	switch orientation {
	case UprightMirrored:
		return imaging.FlipH(img)
	case UpsideDown:
		return imaging.Rotate180(img)
	case UpsideDownMirrored:
		return imaging.FlipV(img)
	case RotatedCWMirrored:
		return imaging.Transpose(img)
	case RotatedCCW:
		return imaging.Rotate270(img)
	case RotatedCCWMirrored:
		return imaging.Transverse(img)
	case RotatedCW:
		return imaging.Rotate90(img)
	default:
		return img
	}
}

func getImageOrientation(input io.Reader) (int, error) {
	if exifData, err := exif.Decode(input); err != nil {
		return Upright, err
	} else {
		if tag, err := exifData.Get("Orientation"); err != nil {
			return Upright, err
		} else {
			orientation, err := tag.Int(0)
			if err != nil {
				return Upright, err
			} else {
				return orientation, nil
			}
		}
	}
}

func writeFileLocally(f []byte, path string) *model.AppError {
	if err := os.MkdirAll(filepath.Dir(path), 0774); err != nil {
		directory, _ := filepath.Abs(filepath.Dir(path))
		return model.NewLocAppError("WriteFile", "api.file.write_file_locally.create_dir.app_error", nil, "directory="+directory+", err="+err.Error())
	}

	if err := ioutil.WriteFile(path, f, 0644); err != nil {
		return model.NewLocAppError("WriteFile", "api.file.write_file_locally.writing.app_error", nil, err.Error())
	}

	return nil
}

func ReadFile(path string) ([]byte, *model.AppError) {
	if f, err := ioutil.ReadFile(*utils.Cfg.FileSettings.Directory + path); err != nil {
		return nil, model.NewLocAppError("ReadFile", "api.file.read_file.reading_local.app_error", nil, err.Error())
	} else {
		return f, nil
	}
}