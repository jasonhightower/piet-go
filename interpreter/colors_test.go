package interpreter

import (
	"image/color"
	"testing"
    "fmt"
)

func assertHue(color *color.Color, expected byte) error {
    hue := Hue(*color)
    if hue != expected {
        return fmt.Errorf("Expected hue of %d but got %d\n", expected, hue)
    }
    return nil
}

func assertLightness(color *color.Color, expected byte) error {
    lightness := Lightness(*color)
    if lightness != expected {
        return fmt.Errorf("Expected lightness of %d but got %d\n", expected, lightness)
    }
    return nil
}

func TestLightness(t *testing.T) {
    err := assertLightness(&l_red, LIGHT)
    if  err != nil {
        t.Error(err)
    }

    err = assertLightness(&n_red, NORMAL)
    if err != nil {
        t.Error(err)
    }

    err = assertLightness(&d_red, DARK)
    if err != nil {
        t.Error(err)
    }

    err = assertLightness(&l_yellow, LIGHT)
    if  err != nil {
        t.Error(err)
    }

    err = assertLightness(&n_yellow, NORMAL)
    if err != nil {
        t.Error(err)
    }

    err = assertLightness(&d_yellow, DARK)
    if err != nil {
        t.Error(err)
    }

    err = assertLightness(&l_green, LIGHT)
    if  err != nil {
        t.Error(err)
    }

    err = assertLightness(&n_green, NORMAL)
    if err != nil {
        t.Error(err)
    }

    err = assertLightness(&d_green, DARK)
    if err != nil {
        t.Error(err)
    }

    err = assertLightness(&l_cyan, LIGHT)
    if  err != nil {
        t.Error(err)
    }

    err = assertLightness(&n_cyan, NORMAL)
    if err != nil {
        t.Error(err)
    }

    err = assertLightness(&d_cyan, DARK)
    if err != nil {
        t.Error(err)
    }

    err = assertLightness(&l_blue, LIGHT)
    if  err != nil {
        t.Error(err)
    }

    err = assertLightness(&n_blue, NORMAL)
    if err != nil {
        t.Error(err)
    }

    err = assertLightness(&d_blue, DARK)
    if err != nil {
        t.Error(err)
    }

    err = assertLightness(&l_magenta, LIGHT)
    if  err != nil {
        t.Error(err)
    }

    err = assertLightness(&n_magenta, NORMAL)
    if err != nil {
        t.Error(err)
    }

    err = assertLightness(&d_magenta, DARK)
    if err != nil {
        t.Error(err)
    }
    // TODO JH other colors besides red
}

func TestHue(t *testing.T) {

    err := assertHue(&l_red, RED)
    if  err != nil {
        t.Error(err)
    }

    err = assertHue(&n_red, RED)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&d_red, RED)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&l_green, GREEN)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&n_green, GREEN)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&d_green, GREEN)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&l_blue, BLUE)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&n_blue, BLUE)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&d_blue, BLUE)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&l_yellow, YELLOW)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&n_yellow, YELLOW)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&d_yellow, YELLOW)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&l_cyan, CYAN)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&n_cyan, CYAN)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&d_cyan, CYAN)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&l_magenta, MAGENTA)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&n_magenta, MAGENTA)
    if err != nil {
        t.Error(err)
    }

    err = assertHue(&d_magenta, MAGENTA)
    if err != nil {
        t.Error(err)
    }
}
