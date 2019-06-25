package image

import (
	"image"

	"github.com/golang/freetype/raster"
)

type (
	//PainterRGBA paints to the RGBA image using any image as a source
	PainterRGBA struct {
		im *image.RGBA
		si image.Image
	}
)

const m = 1<<16 - 1

//NewImagePainterRGBA creates a new RGBA painter instance
func NewImagePainterRGBA(rgba *image.RGBA, src image.Image) *PainterRGBA {
	return &PainterRGBA{im: rgba, si: src}
}

func (ip *PainterRGBA) Paint(ss []raster.Span, done bool) {
	b := ip.im.Bounds()
	for _, s := range ss {
		if s.Y < b.Min.Y {
			continue
		}
		if s.Y >= b.Max.Y {
			return
		}
		if s.X0 < b.Min.X {
			s.X0 = b.Min.X
		}
		if s.X1 > b.Max.X {
			s.X1 = b.Max.X
		}
		if s.X0 >= s.X1 {
			continue
		}
		y := s.Y - ip.im.Rect.Min.Y
		x0 := s.X0 - ip.im.Rect.Min.X
		// RGBAPainter.Paint() in $GOPATH/src/github.com/golang/freetype/raster/paint.go
		i0 := (s.Y-ip.im.Rect.Min.Y)*ip.im.Stride + (s.X0-ip.im.Rect.Min.X)*4
		i1 := i0 + (s.X1-s.X0)*4
		for i, x := i0, x0; i < i1; i, x = i+4, x+1 {
			ma := s.Alpha
			cr, cg, cb, ca := ip.si.At(x, y).RGBA()
			dr := uint32(ip.im.Pix[i+0])
			dg := uint32(ip.im.Pix[i+1])
			db := uint32(ip.im.Pix[i+2])
			da := uint32(ip.im.Pix[i+3])
			a := (m - (ca * ma / m)) * 0x101
			ip.im.Pix[i+0] = uint8((dr*a + cr*ma) / m >> 8)
			ip.im.Pix[i+1] = uint8((dg*a + cg*ma) / m >> 8)
			ip.im.Pix[i+2] = uint8((db*a + cb*ma) / m >> 8)
			ip.im.Pix[i+3] = uint8((da*a + ca*ma) / m >> 8)
		}
	}
}
