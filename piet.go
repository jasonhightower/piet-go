package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"os"
	"jasonhightower.com/piet/interpreter"
    log "github.com/sirupsen/logrus"
)

func init() {
  // Log as JSON instead of the default ASCII formatter.
  //  log.SetFormatter(&log.JSONFormatter{})

  log.SetOutput(os.Stdout)

}

func readPiet(filename string) image.Image {
    file, err := os.Open(filename)
    if err != nil {
        return nil
    }
    defer file.Close()

    image, err := gif.Decode(file)
    return image
}

func main() {
    filename := flag.String("f", "", "name of the piet file to interpret")
    codelsize := flag.Int("codel-size", 1, "Size of codels to support enlarged images for better viewing")
    capacity := flag.Int("capacity", 512, "Capacity of the stack")
    loglevel := flag.String("log", "info", "Log Level")
    flag.Parse()

    image := readPiet(*filename)

    pi := interpreter.NewStackMachine(*capacity)

    if *codelsize > 1 {
        image = interpreter.NewCodelImage(*codelsize, image)
    }

    switch *loglevel {
    case "debug":
        log.SetLevel(log.DebugLevel)
    case "warn":
        log.SetLevel(log.WarnLevel)
    case "error":
        log.SetLevel(log.ErrorLevel)
    case "info":
        log.SetLevel(log.InfoLevel)
    }

    pi.Execute(image)
    fmt.Println()
}


