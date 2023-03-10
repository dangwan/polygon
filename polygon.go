package polygon

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

// Point 判断点是否位于区域内，此方法只适用于矩形，无法应用于所有情况
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"Y"`
}

// Polygon 目标图片多边形信息
type Polygon struct {
	//X      int     `json:"x"`
	//Y      int     `json:"y"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Area   []Point `json:"area"` // 顺时针坐标
}

// SaveImage save image to local
func SaveImage(img image.Image, path string) error {
	out := bytes.NewBuffer(make([]byte, 0))
	if err := jpeg.Encode(out, img, &jpeg.Options{Quality: 100}); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	if err = png.Encode(b, img); err != nil {
		return err
	}
	err = b.Flush()
	return err
}

func imageToRGBA(src image.Image) *image.RGBA {
	// No conversion needed if image is an *image.RGBA.
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}
	// Use the image/draw package to convert to *image.RGBA.
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}

// ClipForPolygon clip image to target polygon
func ClipForPolygon(t *Polygon, img image.Image) (*image.RGBA, error) {
	if t == nil {
		return nil, fmt.Errorf("clip image is null")
	}
	if len(t.Area) == 0 {
		return imageToRGBA(img), nil
	}
	clipImg := image.NewRGBA(image.Rect(0, 0, t.Width, t.Height))
	draw.DrawMask(clipImg, clipImg.Bounds(), img, image.Point{}, t, image.Point{}, draw.Over)
	return clipImg, nil
}

// ColorModel set color model
func (c *Polygon) ColorModel() color.Model {
	return color.AlphaModel
}

// Bounds set image bounds
func (c *Polygon) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.Width, c.Height)
}

// At judge point alpha
func (c *Polygon) At(x, y int) color.Color {
	if c.PointInArea(float64(x), float64(y)) {
		return color.Alpha{A: 255} // save pixel
	}
	return color.Alpha{}
}

// PointInArea 判断一个点是否在一个面内
func (c *Polygon) PointInArea(x, y float64) bool {
	pointNum := len(c.Area) //点个数
	intersectCount := 0     //cross points count of x
	precision := 2e-10      //浮点类型计算时候与0比较时候的容差
	p1 := Point{}           //neighbour bound vertices
	p2 := Point{}
	p := Point{x, y} //测试点

	p1 = c.Area[0] //left vertex
	for i := 0; i < pointNum; i++ {
		if p.X == p1.X && p.Y == p1.Y {
			return true
		}
		p2 = c.Area[i%pointNum]
		if p.Y < math.Min(p1.Y, p2.Y) || p.Y > math.Max(p1.Y, p2.Y) {
			p1 = p2
			continue //next ray left point
		}

		if p.Y > math.Min(p1.Y, p2.Y) && p.Y < math.Max(p1.Y, p2.Y) {
			if p.X <= math.Max(p1.X, p2.X) { //x is before of ray
				if p1.Y == p2.Y && p.X >= math.Min(p1.X, p2.X) {
					return true
				}

				if p1.X == p2.X { //ray is vertical
					if p1.X == p.X { //overlies on a vertical ray
						return true
					} else { //before ray
						intersectCount++
					}
				} else { //cross point on the left side
					xinters := (p.Y-p1.Y)*(p2.X-p1.X)/(p2.Y-p1.Y) + p1.X
					if math.Abs(p.X-xinters) < precision {
						return true
					}

					if p.X < xinters { //before ray
						intersectCount++
					}
				}
			}
		} else { //special case when ray is crossing through the vertex
			if p.Y == p2.Y && p.X <= p2.X { //p crossing over p2
				p3 := c.Area[(i+1)%pointNum]
				if p.Y >= math.Min(p1.Y, p3.Y) && p.Y <= math.Max(p1.Y, p3.Y) {
					intersectCount++
				} else {
					intersectCount += 2
				}
			}
		}
		p1 = p2 //next ray left point
	}
	if intersectCount%2 == 0 { //偶数在多边形外
		return false
	} else { //奇数在多边形内
		return true
	}
}
