package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"os"
	"jasonhightower.com/piet/interpreter"
)

func readPiet(filename string) image.Image {
    file, err := os.Open(filename)
    if err != nil {
        return nil
    }
    defer file.Close()

    image, err := gif.Decode(file)
    return image
}

func print_colorAt(image image.Image, x int, y int) {
    r, b, g, a := image.At(x, y).RGBA()
    fmt.Printf("(%d, %d) - (%d, %d, %d, %d)\n", x, y, r, g, b, a)
}

func main() {
    filename := flag.String("f", "", "name of the piet file to interpret")
    codelsize := flag.Int("codel-size", 1, "Size of codels to support enlarged images for better viewing")
    flag.Parse()

    image := readPiet(*filename)
//    bounds := image.Bounds().Size()

    pi := interpreter.PietInterpreter{Csize: *codelsize}
    pi.Execute(image)

    /*
    fmt.Printf("Codel size (%d)\n", codelsize)
    fmt.Printf("Image (%d, %d)\n", bounds.X, bounds.Y)
    r, b, g, a := image.At(0, 0).RGBA()
    fmt.Printf("Color (%d, %d, %d, %d)\n", r, g, b, a)

    for x := 0; x < bounds.X; x += *codelsize {
        print_colorAt(image, x, 0)
    }
    */
}


