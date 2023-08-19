package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	log "github.com/sirupsen/logrus"
	"jasonhightower.com/piet/interpreter"
)

func init() {
  // Log as JSON instead of the default ASCII formatter.
  //  log.SetFormatter(&log.JSONFormatter{})

  log.SetOutput(os.Stdout)
}

func readImage(filename string) (image.Image, error) {
    file, err := os.Open(filename)
    if err != nil {
        // FIXME proper error handling
        return nil, err
    }
    defer file.Close()

    image, _, err := image.Decode(file)
    return image, err
}

func main() {
    filename := flag.String("f", "", "name of the piet file to interpret")
    codelsize := flag.Int("codel-size", 1, "Size of codels to support enlarged images for better viewing")
    capacity := flag.Int("capacity", 512, "Capacity of the stack")
    mode := flag.String("m", "interpret", "Name of the mode to run")
    loglevel := flag.String("log", "info", "Log Level")
    flag.Parse()

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

    if image, err := readImage(*filename); err == nil {

        if *codelsize > 1 {
            image = interpreter.NewCodelImage(*codelsize, image)
        }

        pi := interpreter.NewInterpreter(*capacity)
        if *mode == "interpret" {
            pi.Interpret(image)
        } else {
            tokens := interpreter.Tokenize(image)
            interpreter.ParseTokens(tokens)
            pi.InterpretTokens(tokens[0][0])
        }
        fmt.Println()
    } else {
        log.Errorf("Unable to decode image %s - %v", *filename, err)
    }
}
