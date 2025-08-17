package main

import (
	"image"
	"image/color"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// each braille character represents 2x4 pixels: ⣿
const BRAILLE_WIDTH int = 2
const BRAILLE_HEIGHT int = 4

type BrailleImage struct {
    cols uint
    rows uint
    image *image.RGBA
    dots *[]uint8
}

func NewBrailleImage(cols uint, rows uint, img *image.RGBA) *BrailleImage {
    data := make([]uint8, cols * rows)

    return &BrailleImage{ cols: cols, rows: rows, image: img, dots: &data }
}

func ImageFromBrailleString(s string) *image.RGBA {
    lines := strings.Split(s, "\n")
    height := (len(lines) - 1) * BRAILLE_HEIGHT
    var width int

    {
        maxWidth := 0
        for _, line := range lines {
            maxWidth = max(maxWidth, utf8.RuneCountInString(line))
        }

        width = maxWidth * BRAILLE_WIDTH
    }

    img := image.NewRGBA(image.Rect(0, 0, width, height))

    for row, line := range lines {
        for col, ch := range utf16.Encode([]rune(line)) {
            mask := ch & 0b11111111
            var v uint8 = 0

            // 0,0
            if mask & 0b1 != 0 { v = 255 } else { v = 0 }
            img.SetRGBA(col * BRAILLE_WIDTH + 0, row * BRAILLE_HEIGHT + 0, color.RGBA{v, v, v, 255})

            // 1,0
            if mask & 0b1000 != 0 { v = 255 } else { v = 0 }
            img.SetRGBA(col * BRAILLE_WIDTH + 1, row * BRAILLE_HEIGHT + 0, color.RGBA{v, v, v, 255})

            // 0,1
            if mask & 0b10 != 0 { v = 255 } else { v = 0 }
            img.SetRGBA(col * BRAILLE_WIDTH + 0, row * BRAILLE_HEIGHT + 1, color.RGBA{v, v, v, 255})

            // 1,1
            if mask & 0b10000 != 0 { v = 255 } else { v = 0 }
            img.SetRGBA(col * BRAILLE_WIDTH + 1, row * BRAILLE_HEIGHT + 1, color.RGBA{v, v, v, 255})

            // 0,2
            if mask & 0b100 != 0 { v = 255 } else { v = 0 }
            img.SetRGBA(col * BRAILLE_WIDTH + 0, row * BRAILLE_HEIGHT + 2, color.RGBA{v, v, v, 255})

            // 1,2
            if mask & 0b100000 != 0 { v = 255 } else { v = 0 }
            img.SetRGBA(col * BRAILLE_WIDTH + 1, row * BRAILLE_HEIGHT + 2, color.RGBA{v, v, v, 255})

            // 0,3
            if mask & 0b1000000 != 0 { v = 255 } else { v = 0 }
            img.SetRGBA(col * BRAILLE_WIDTH + 0, row * BRAILLE_HEIGHT + 3, color.RGBA{v, v, v, 255})

            // 1,3
            if mask & 0b10000000 != 0 { v = 255 } else { v = 0 }
            img.SetRGBA(col * BRAILLE_WIDTH + 1, row * BRAILLE_HEIGHT + 3, color.RGBA{v, v, v, 255})
        }
    }

    return img
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

    for y := uint(0); y < h && yd < img.rows; y++ {
        for x := uint(0); x < w && xd < img.cols; x++ {
            lum = img.GetLuminance(x - uint(bounds.Min.X), y - uint(bounds.Min.Y))

            if debugMode {
                imgLum.SetRGBA(int(x), int(y), color.RGBA{lum, lum, lum, 255})
            }

            if invert {
                lum = 255 - lum
            }

            if uint(lum) >= threshold {
                (*img.dots)[yd * img.cols + xd] |=
                    ((1 << (x % uint(BRAILLE_WIDTH))) << ((y % uint(BRAILLE_HEIGHT)) << 1))
            }

            if (x + 1) % uint(BRAILLE_WIDTH) == 0 {
                xd++;
            }
        }

        xd = 0
        if (y + 1) % uint(BRAILLE_HEIGHT) == 0 {
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
            (*img.image).SetRGBA(
                x,
                y,
                color.RGBA{
                    uint8(Clamp(int(f * (float64(c.R) - 128.0) + 128.0), 0, 255)),
                    uint8(Clamp(int(f * (float64(c.G) - 128.0) + 128.0), 0, 255)),
                    uint8(Clamp(int(f * (float64(c.B) - 128.0) + 128.0), 0, 255)),
                    c.A,
                },
            )
        }
    }
}

func (img *BrailleImage) String() string {
    s := strings.Builder{}

    for y := uint(0); y < img.rows; y++ {
        for x := uint(0); x < img.cols; x++ {
            s.WriteRune(IntToBrailleRune((*img.dots)[y * img.cols + x]));
        }
        s.WriteString("\n")
    }

    return s.String()
}
