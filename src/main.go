package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"

	"github.com/nfnt/resize"
)

func eprintln(msg string) {
    fmt.Fprintln(os.Stderr, msg)
}

func eprintf(format string, args ...any) {
    fmt.Fprintf(os.Stderr, format, args...)
}

var debugMode bool = false
var oWidth, oHeight, oThreshold, oContrast *int
var oInvert *bool
var oImageFromText *bool
var oScalingMode *string

func init() {
    flag.Usage = func() {
        eprintf(
            "Usage:\n" +
            "   %s [OPTIONS] IMAGE_FILE\n" +
            "   %s [OPTIONS] -\n\n" +
            "If either one of width or height is set, the other is calculated based on the\n" +
            "aspect ratio of the image. If no dimensions are specified, a default width of\n" +
            "80 is applied. If both dimensions are specified, aspect ratio is not preserved.\n" +
            "\n" +
            "Specify IMAGE_FILE as '-' to read the image file from stdin.\n\n",
            os.Args[0],
            os.Args[0],
        )
        flag.PrintDefaults()
    }

    oWidth = flag.Int("w", -1, "Width in number of characters.")
    oHeight = flag.Int("h", -1, "Height in number of characters.")
    oThreshold = flag.Int("t", 85, "Luminance threshold.")
    oContrast = flag.Int("c", 0, "Contrast. A positive or negative integer.")
    oInvert = flag.Bool("i", false, "Invert the output.")
    oImageFromText = flag.Bool(
        "s",
        false,
        "Treat and parse IMAGE_FILE as a braille image text file.",
    )
    oScalingMode = flag.String(
        "r",
        "bicubic",
        "Image scaling algorithm. (bicubic, bilinear, nearest-neighbor)",
    )
}

func readImageFile(arg string) (image.Image, error) {
    var fileData []byte

    if arg == "-" {
        // read from stdin
        data, err := io.ReadAll(os.Stdin)
        if err != nil {
            return nil, fmt.Errorf("Failed to read from stdin: %w", err)
        }
        fileData = data
    } else {
        imgFile, err := os.Open(arg)
        if err != nil {
            return nil, errors.New("Failed to open image file!")
        }
        defer imgFile.Close()

        data, err := io.ReadAll(imgFile)
        if err != nil {
            return nil, fmt.Errorf("Failed to read image file: %w", err)
        }
        fileData = data
    }

    var imgData image.Image

    if *oImageFromText {
        imgData = ImageFromBrailleString(string(fileData))

        if debugMode {
            WritePngImage(ImageToRGBA(&imgData), "parsed.png")
        }
    } else {
        fileReader := bytes.NewReader(fileData)
        _, imgType, err := image.Decode(fileReader)
        if err != nil {
            return nil, errors.New("Failed to decode image!")
        }
        fileReader.Seek(0,0)

        switch (imgType) {
        case "png":
            imgData, err = png.Decode(fileReader)

        case "jpeg":
            imgData, err = jpeg.Decode(fileReader)

        default:
            return nil, errors.New("Unsupported image format!")
        }

        if err != nil {
            return nil, errors.New("Failed to decode image!")
        }
    }

    return imgData, nil
}

func main() {
    flag.Parse()
    var args []string = flag.Args()

    if os.Getenv("DEBUG") == "1" {
        debugMode = true
        eprintln("DEBUG MODE")
    }

    if len(args) < 1 {
        eprintln("No image file specified!")
        os.Exit(1)
    }

    imgData, err := readImageFile(args[0])

    if err != nil {
        eprintln(err.Error())
        os.Exit(1)
    }

    var cols uint
    var rows uint
    bounds := imgData.Bounds()
    var imgWidth uint = uint(bounds.Max.X - bounds.Min.X)
    var imgHeight uint = uint(bounds.Max.Y - bounds.Min.Y)

    var aspect float64 =
        (float64(imgWidth - (imgWidth % uint(BRAILLE_WIDTH))) / float64(BRAILLE_WIDTH)) /
        (float64(imgHeight - (imgHeight % uint(BRAILLE_HEIGHT))) / float64(BRAILLE_HEIGHT))

    // if no dimensions were specified: use default
    if *oWidth == -1 && *oHeight == -1 {
        *oWidth = 80
    }

    if *oWidth != -1 && *oHeight != -1 {
        // both width and height set by user: use custom aspect ratio
        cols = uint(*oWidth)
        rows = uint(*oHeight)
        aspect = (float64(cols) * 2.0) / (float64(rows) * 4.0)
    } else if *oWidth != -1 || (*oWidth == -1 && *oHeight == -1) {
        cols = uint(*oWidth)
        rows = uint(math.Round(float64(cols) / aspect))
    } else if *oHeight != -1 {
        rows = uint(*oHeight)
        cols = uint(math.Round(float64(rows) * aspect))
    }

    var inWidth uint = uint(imgData.Bounds().Dx())
    var inHeight uint = uint(imgData.Bounds().Dy())
    var outWidth uint = cols * 2
    var outHeight uint = rows * 4
    var imgRGBA *image.RGBA
    var scalingMode resize.InterpolationFunction
    var useScaling = inWidth != outWidth || inHeight != outHeight

    switch *oScalingMode {
    case "bicubic":
        scalingMode = resize.Bicubic
    case "bilinear":
        scalingMode = resize.Bilinear
    case "nearest-neighbor":
        scalingMode = resize.NearestNeighbor
    default:
        eprintf("Unknown scaling mode: '%s'\n", *oScalingMode)
        os.Exit(1)
    }

    if debugMode {
        scalingInfo := "none"
        if useScaling { scalingInfo = *oScalingMode }
        inCols := float64(inWidth) / float64(BRAILLE_WIDTH)
        inRows := float64(inHeight) / float64(BRAILLE_HEIGHT)

        eprintf(
            "inCols: %f, inRows: %f, inW: %d, inH: %d\n",
            inCols,
            inRows,
            inWidth,
            inHeight,
        )
        eprintf(
            "outCols: %d, outRows: %d, outW: %d, outH: %d\n",
            cols,
            rows,
            outWidth,
            outHeight,
        )
        eprintf(
            "threshold: %d, contrast: %d, scaling: %s\n\n",
            *oThreshold,
            *oContrast,
            scalingInfo,
        )
    }

    if useScaling {
        imgResized := resize.Resize(outWidth, outHeight, imgData, scalingMode)
        imgRGBA = ImageToRGBA(&imgResized)
    } else {
        imgRGBA = ImageToRGBA(&imgData)
    }

    imgBraille := NewBrailleImage(cols, rows, imgRGBA)

    if *oContrast != 0 {
        imgBraille.ModContrast(*oContrast)
        if debugMode {
            WritePngImage(imgRGBA, "contrast.png")
        }
    }

    imgBraille.FillDotsData(uint(*oThreshold), *oInvert)
    fmt.Println(imgBraille.String())
}
