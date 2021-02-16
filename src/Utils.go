package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"reflect"

	"golang.org/x/image/bmp"
)

func WriteBmpImage(img *image.RGBA, path string) {
    outputFile, err := os.Create(path)
    if err != nil {
        println("Failed to write image!")
        defer outputFile.Close()
        return
    }

    subimg := img.SubImage(img.Bounds())

    bmp.Encode(outputFile, subimg)
    if err != nil {
        println("Failed to encode image!")
    }

    outputFile.Close()
}

func WritePngImage(img *image.RGBA, path string) {
    outputFile, err := os.Create(path)
    if err != nil {
        println("Failed to write image!")
        defer outputFile.Close()
        return
    }

    subimg := img.SubImage(img.Bounds())

    err = png.Encode(outputFile, subimg)
    if err != nil {
        println("Failed to encode image!")
    }

    outputFile.Close()
}

func ImageToRGBA(img *image.Image) *image.RGBA {
    bounds := (*img).Bounds()
    result := image.NewRGBA(bounds)

    for y := 0; y < bounds.Max.Y; y++ {
        for x := 0; x < bounds.Max.X; x++ {
            r, g, b, a := (*img).At(x, y).RGBA()
            result.SetRGBA(
                x, y, color.RGBA{
                    uint8(r >> 8),
                    uint8(g >> 8),
                    uint8(b >> 8),
                    uint8(a >> 8) })
        }
    }

    return result
}

func Clamp(value int, min int, max int) int {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}

func Typeof(v interface{}) string {
    return reflect.TypeOf(v).String()
}
