package services

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

const thumbMaxWidth = 300
const thumbQuality = 80

func GenerateThumbnail(srcPath, thumbDir string) (string, error) {
	src, err := os.Open(srcPath)
	if err != nil {
		return "", fmt.Errorf("open source: %w", err)
	}
	defer src.Close()

	img, format, err := image.Decode(src)
	if err != nil {
		return "", fmt.Errorf("decode image: %w", err)
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w <= thumbMaxWidth {
		return srcPath, nil
	}

	ratio := float64(thumbMaxWidth) / float64(w)
	newH := int(float64(h) * ratio)

	dst := image.NewRGBA(image.Rect(0, 0, thumbMaxWidth, newH))
	draw.BiLinear.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)

	ext := filepath.Ext(srcPath)
	base := strings.TrimSuffix(filepath.Base(srcPath), ext)
	thumbPath := filepath.Join(thumbDir, base+"_thumb.jpg")

	out, err := os.Create(thumbPath)
	if err != nil {
		return "", fmt.Errorf("create thumbnail: %w", err)
	}
	defer out.Close()

	if format == "png" {
		if err := png.Encode(out, dst); err != nil {
			return "", fmt.Errorf("encode thumbnail: %w", err)
		}
	} else {
		if err := jpeg.Encode(out, dst, &jpeg.Options{Quality: thumbQuality}); err != nil {
			return "", fmt.Errorf("encode thumbnail: %w", err)
		}
	}

	return thumbPath, nil
}
