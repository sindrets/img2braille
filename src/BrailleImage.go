package main

import (
	"image"
	"image/color"
	"strings"
	"unicode/utf16"
)

type BrailleImage struct {
    width uint
    height uint
    image *image.RGBA
    dots *[]uint8
}

func NewBrailleImage(width uint, height uint, image *image.RGBA) *BrailleImage {
    data := make([]uint8, width * height)

    return &BrailleImage{ width: width, height: height, image: image, dots: &data }
}

func (img *BrailleImage) GetLuminance(x uint, y uint) uint8 {
    c := (*img.image).RGBAAt(int(x), int(y))
    return uint8(Clamp((
        int(c.R) * 299 +
        int(c.G) * 587 +
        int(c.B) * 114) / 1000, 0, 255));
}

func (img *BrailleImage) FillDotsData(threshold uint, invert bool) {
    bounds := (*img.image).Bounds()
    var w uint = uint(bounds.Max.X - bounds.Min.X)
    var h uint = uint(bounds.Max.Y - bounds.Min.Y)
    var xd uint = 0
    var yd uint = 0
    var lum uint8
    var imgLum *image.RGBA
    if invert {
        threshold = 255 - threshold
    }

    if debugMode {
        imgLum = image.NewRGBA(img.image.Rect)
    }

    for y := uint(0); y < h && yd < img.height; y++ {
        for x := uint(0); x < w && xd < img.width; x++ {
            lum = img.GetLuminance(x - uint(bounds.Min.X), y - uint(bounds.Min.Y))

            if debugMode {
                imgLum.SetRGBA(int(x), int(y), color.RGBA{lum, lum, lum, 255})
            }

            if invert {
                lum = 255 - lum
            }

            if uint(lum) >= threshold {
                (*img.dots)[yd * img.width + xd] |= ((1 << (x % 2)) << ((y % 4) << 1))
            }

            if (x + 1) % 2 == 0 {
                xd++;
            }
        }

        xd = 0
        if (y + 1) % 4 == 0 {
            yd++;
        }
    }

    if debugMode {
        WritePngImage(imgLum, "lum.png")
    }
}

func IntToBrailleRune(dots uint8) rune {
    var mask uint16 = 0

    if dots & 0b1 != 0 { // 0,0
        mask |= 0b1
    }
    if dots & 0b10 != 0 { // 1,0
        mask |= 0b1000
    }
    if dots & 0b100 != 0 { // 0,1
        mask |= 0b10
    }
    if dots & 0b1000 != 0 { // 1,1
        mask |= 0b10000
    }
    if dots & 0b10000 != 0 { // 0,2
        mask |= 0b100
    }
    if dots & 0b100000 != 0 { // 1,2
        mask |= 0b100000
    }
    if dots & 0b1000000 != 0 { // 0,3
        mask |= 0b1000000
    }
    if dots & 0b10000000 != 0 { // 1,3
        mask |= 0b10000000
    }

    return utf16.Decode([]uint16{ 0x2800 | mask })[0]
}

func (img *BrailleImage) ModContrast(amount int) {
    f := (259.0 * (float64(amount) + 255.0)) / (255.0 * (259.0 - float64(amount)))
    bounds := (*img.image).Bounds()

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            c := (*img.image).RGBAAt(x, y)
            (*img.image).SetRGBA(x, y, color.RGBA{
                uint8(Clamp(int(f * (float64(c.R) - 128.0) + 128.0), 0, 255)),
                uint8(Clamp(int(f * (float64(c.G) - 128.0) + 128.0), 0, 255)),
                uint8(Clamp(int(f * (float64(c.B) - 128.0) + 128.0), 0, 255)),
                c.A })
        }
    }
}

func (img *BrailleImage) String() string {
    s := strings.Builder{}

    for y := uint(0); y < img.height; y++ {
        for x := uint(0); x < img.width; x++ {
            s.WriteRune(IntToBrailleRune((*img.dots)[y * img.width + x]));
        }
        s.WriteString("\n")
    }

    return s.String()
}
