# img2braille

Convert images to braille ascii-art. Best results are achieved with images that
have clear distinct shapes and outlines.

### Usage

```
Usage:
   img2braille [OPTIONS] IMAGE_FILE
   img2braille [OPTIONS] -

If either one of width or height is set, the other is calculated based on the
aspect ratio of the image. If no dimensions are specified, a default width of
80 is applied. If both dimensions are specified, aspect ratio is not preserved.

Specify IMAGE_FILE as '-' to read the image file from stdin.

  -c int
        Contrast. A positive or negative integer.
  -h int
        Height in number of characters. (default -1)
  -i    Invert the output.
  -r string
        Image scaling algorithm. (bicubic, bilinear, nearest-neighbor) (default "bicubic")
  -s    Treat and parse IMAGE_FILE as a braille image text file.
  -t int
        Luminance threshold. (default 85)
  -w int
        Width in number of characters. (default -1)
```

If you define the env variable `DEBUG=1` before you run the program, it will
output up to two files depending on your options: `contrast.png` and `lum.png`.
These files make it easier to understand exactly how the applied contrast
affects the luminance of the image. The luminance is used to determine which
dots are present in the braille characters.

### Examples

![Example-1](https://i.imgur.com/2vXmmmz.png)

![Example-2](https://i.imgur.com/lgMQRxn.png)

![Example-3](https://i.imgur.com/NUxjf4s.png)


## License

This work is licensed under the [GNU General Public License v3.0 only](LICENSE)
