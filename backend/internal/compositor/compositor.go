package compositor

import (
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png" // register PNG decoder
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Compositor struct {
	BaseAssetPath string
}

func NewCompositor(basePath string) *Compositor {
	return &Compositor{BaseAssetPath: basePath}
}

func (c *Compositor) Compose(imageURLs []string, layout string) (string, error) {
	if len(imageURLs) == 0 {
		return "", os.ErrNotExist
	}

	var imgs []image.Image
	for _, u := range imageURLs {
		resp, err := http.Get(u)
		if err != nil {
			return "", err
		}
		img, _, err := image.Decode(resp.Body)
		resp.Body.Close()
		if err != nil {
			// fallback: create an empty white square if decode fails
			img = image.NewRGBA(image.Rect(0, 0, 512, 512))
		}
		imgs = append(imgs, img)
	}

	// In a real system, we'd scale the images to be uniform, but we assume 512x512 uniformly generated
	size := imgs[0].Bounds().Max

	var finalImg draw.Image
	if layout == "1x4" {
		finalImg = image.NewRGBA(image.Rect(0, 0, size.X, size.Y*len(imgs)))
		for i, img := range imgs {
			draw.Draw(finalImg, image.Rect(0, i*size.Y, size.X, (i+1)*size.Y), img, image.Point{0, 0}, draw.Src)
		}
	} else {
		// Default to 2x2. Assumes exactly 4 images for standard Yon Koma
		finalImg = image.NewRGBA(image.Rect(0, 0, size.X*2, size.Y*2))
		if len(imgs) >= 4 {
			draw.Draw(finalImg, image.Rect(0, 0, size.X, size.Y), imgs[0], image.Point{0, 0}, draw.Src)
			draw.Draw(finalImg, image.Rect(size.X, 0, size.X*2, size.Y), imgs[1], image.Point{0, 0}, draw.Src)
			draw.Draw(finalImg, image.Rect(0, size.Y, size.X, size.Y*2), imgs[2], image.Point{0, 0}, draw.Src)
			draw.Draw(finalImg, image.Rect(size.X, size.Y, size.X*2, size.Y*2), imgs[3], image.Point{0, 0}, draw.Src)
		}
	}

	fileName := "comic_" + uuid.New().String() + ".jpg"
	saveDir := filepath.Join(c.BaseAssetPath, "comics")
	_ = os.MkdirAll(saveDir, 0755)

	fullPath := filepath.Join(saveDir, fileName)
	outFile, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if err := jpeg.Encode(outFile, finalImg, &jpeg.Options{Quality: 90}); err != nil {
		return "", err
	}

	return "/assets/comics/" + fileName, nil
}
