package main

import (
	"github.com/dangwan/polygon"
	"image"
	"os"
)

func test1() {
	f, _ := os.Open("tmp.jpeg")
	img, _, _ := image.Decode(f)
	c := polygon.Polygon{
		Width:  img.Bounds().Dx(),
		Height: img.Bounds().Dy(),
		Area: []polygon.Point{
			{0, 86},
			{700, 0},
			{808, 480},
			{100, 450},
		},
	}
	rgbaImg, _ := polygon.ClipForPolygon(&c, img)
	polygon.SaveImage(rgbaImg, "tmp-crop.jpeg")
}
func test2() {
	f, _ := os.Open("garden.jpeg")
	img, _, _ := image.Decode(f)
	c := polygon.Polygon{
		Width:  img.Bounds().Dx(),
		Height: img.Bounds().Dy(),
		Area: []polygon.Point{
			{200, 0},
			{600, 300},
			{100, 350},
			{0, 100},
		},
	}
	rgbaImg, _ := polygon.ClipForPolygon(&c, img)
	polygon.SaveImage(rgbaImg, "garden-crop.jpeg")
}
func main() {
	test1()
	test2()

}
