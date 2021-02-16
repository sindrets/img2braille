package main

import (
    "flag"
    "fmt"
    "image"
    "image/jpeg"
    "image/png"
    "log"
    "math"
    "os"

    "github.com/nfnt/resize"
)

var debugMode bool = false
var oWidth, oHeight, oThreshold, oContrast *int
var oInvert *bool

func init() {
    flag.Usage = func() {
        fmt.Printf(
            "Usage: %s [OPTIONS] IMAGE_FILE\n\n" +
            "If either one of width or height is set, the other is calculated based on the\n" +
            "aspect ratio of the image. If no dimensions are specified, a default width of\n" +
            "80 is applied. If both dimensions are specified, aspect ratio is not preserved.\n\n",
            os.Args[0])
        flag.PrintDefaults()
    }

    oWidth = flag.Int("w", -1, "Width in number of characters.")
    oHeight = flag.Int("h", -1, "Height in number of characters.")
    oThreshold = flag.Int("t", 85, "Luminance threshold.")
    oContrast = flag.Int("c", 0, "Contrast. A positive or negative integer.")
    oInvert = flag.Bool("i", false, "Invert the output.")
}

func main() {
    flag.Parse()
    var args []string = flag.Args()

    if os.Getenv("DEBUG") == "1" {
        debugMode = true
        println("DEBUG MODE")
    }

    if len(args) < 1 {
        println("No image file specified!")
        return
    }

    imgFile, err := os.Open(args[0])
    if err != nil {
        println("Failed to open image file!")
        imgFile.Close()
        return
    }
    defer imgFile.Close()

    _, imgType, err := image.Decode(imgFile)
    if err != nil {
        println("Failed to decode image!")
        return
    }
    imgFile.Seek(0,0)

    var imgData image.Image
    switch (imgType) {
    case "png":
        imgData, err = png.Decode(imgFile)
        break

    case "jpeg":
        imgData, err = jpeg.Decode(imgFile)
        break

    default:
        log.Fatalln("Unsupported image format!")
    }

    if err != nil {
        println("Failed to decode image!")
        return
    }

    var width uint
    var height uint
    bounds := imgData.Bounds()
    var imgWidth uint = uint(bounds.Max.X - bounds.Min.X)
    var imgHeight uint = uint(bounds.Max.Y - bounds.Min.Y)

    // devide by 2 and 4 respectively because each braille character represents 2x4 pixels: ⣿
    var aspect float64 =
        (float64(imgWidth - (imgWidth % 2)) / 2.0) / (float64(imgHeight - (imgHeight % 4)) / 4.0)

    // if no dimensions were specified: use default
    if *oWidth == -1 && *oHeight == -1 {
        *oWidth = 80
    }

    if *oWidth != -1 && *oHeight != -1 {
        // both width and height set by user: use custom aspect ratio
        width = uint(*oWidth)
        height = uint(*oHeight)
        aspect = (float64(width) * 2.0) / (float64(height) * 4.0)
    } else if *oWidth != -1 || (*oWidth == -1 && *oHeight == -1) {
        width = uint(*oWidth)
        height = uint(math.Round(float64(width) / aspect))
    } else if *oHeight != -1 {
        height = uint(*oHeight)
        width = uint(math.Round(float64(height) * aspect))
    }

    var resWidth uint = width * 2
    var resHeight uint = height * 4

    if debugMode {
        fmt.Printf("w: %d, h: %d, resW: %d, resH: %d, contrast: %d\n", 
            width, height, resWidth, resHeight, *oContrast)
    }

    imgResized := resize.Resize(resWidth, resHeight, imgData, resize.Bicubic)
    imgRGBA := ImageToRGBA(&imgResized)

    imgBraille := NewBrailleImage(width, height, imgRGBA)

    if *oContrast != 0 {
        imgBraille.ModContrast(*oContrast)
        if debugMode {
            WritePngImage(imgRGBA, "contrast.png")
        }
    }

    imgBraille.FillDotsData(uint(*oThreshold), *oInvert)
    fmt.Println(imgBraille.String())

    return
}
