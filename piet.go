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
    mode := flag.String("m", "interpret", "Name of the mode to run")
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

    if *mode == "interpret" {
        pi.Execute(image)
    } else {
        shapes := interpreter.Tokenize(image)
        interpreter.ParseGraph(shapes[0][0], shapes)

        printed := map[*interpreter.Shape]bool{}
        printGraph(shapes[0][0], 1, "", printed)
    }
    fmt.Println()

}

func printGraph(shape *interpreter.Shape, depth int, name string, printed map[*interpreter.Shape]bool) {
    for i:=0; i < depth; i++ {
        fmt.Print(" ")
    }
    if contains, _ := printed[shape]; contains {
        fmt.Println("<printed>")
        return
    }

    fmt.Print(name)
    if shape == nil {
        fmt.Println(" nil")
        return
    }
    printed[shape] = true
    fmt.Printf(" %s\n", shape.Color())
    printGraph(shape.Connection(interpreter.DpRight, interpreter.CcLeft), depth + 1, "R-L", printed)
    printGraph(shape.Connection(interpreter.DpRight, interpreter.CcRight), depth + 1, "R-R", printed)
    printGraph(shape.Connection(interpreter.DpDown, interpreter.CcLeft), depth + 1, "D-L", printed)
    printGraph(shape.Connection(interpreter.DpDown, interpreter.CcRight), depth + 1, "D-R", printed)
    printGraph(shape.Connection(interpreter.DpLeft, interpreter.CcLeft), depth + 1, "L-L", printed)
    printGraph(shape.Connection(interpreter.DpLeft, interpreter.CcRight), depth + 1, "L-R", printed)
    printGraph(shape.Connection(interpreter.DpUp, interpreter.CcLeft), depth + 1, "U-L", printed)
    printGraph(shape.Connection(interpreter.DpUp, interpreter.CcRight), depth + 1, "U-R", printed)
}



