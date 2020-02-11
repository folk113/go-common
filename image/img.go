package image

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

func Circle(newFileName ,rawFileName string) {
	file, err := os.Create(newFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	imageFile, err := os.Open(rawFileName)
	if err != nil {
		fmt.Println(err)
	}
	defer imageFile.Close()

	srcImg, _ := png.Decode(imageFile)
	w := srcImg.Bounds().Max.X - srcImg.Bounds().Min.X
	h := srcImg.Bounds().Max.Y - srcImg.Bounds().Min.Y

	d := w
	if w > h {
		d = h
	}
	maskImg := circleMask(d)
	dstImg := image.NewRGBA(image.Rect(0, 0, d, d))
	draw.DrawMask(dstImg, srcImg.Bounds().Add(image.Pt(0, 0)), srcImg, image.Pt((w-d)/2, (h-d)/2), maskImg, image.Pt(0, 0), draw.Src)
	png.Encode(file, dstImg)
}

func circleMask(d int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, d, d))
	for x := 0; x < d; x++ {
		for y := 0; y < d; y++ {
			dis := math.Sqrt(math.Pow(float64(x-d/2), 2) + math.Pow(float64(y-d/2), 2))
			if dis > float64(d)/2 {
				img.Set(x, y, color.RGBA{255, 255, 255, 0})
			} else {
				img.Set(x, y, color.RGBA{0, 0, 255, 255})
			}
		}
	}
	return img
}
