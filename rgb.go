package led

import (
	"fmt"
	"image/color"
)

type RGBLED interface {
	LinuxLED
	SetColor(color.Color) error
	GetColor() color.Color
}

type rgbled struct {
	red, green, blue, global *linuxled
	color color.Color

	trigger Trigger
}

func NewRGBLED(red, green, blue, global string) (RGBLED, error) {
	var rgb rgbled
	var led LinuxLED

	led, err := NewLED(red)
	if err != nil {
		return nil, err
	}
	rgb.red = led.(*linuxled)

	led, err = NewLED(green)
	if err != nil {
		return nil, err
	}
	rgb.green = led.(*linuxled)

	led, err = NewLED(blue)
	if err != nil {
		return nil, err
	}
	rgb.blue = led.(*linuxled)

	if global != "" {
		led, err = NewLED(global)
		if err != nil {
			return nil, err
		}
		rgb.global = led.(*linuxled)
	}

	return &rgb, nil
}

func (rgb *rgbled) SetBrightness(brightness float32) error {
	if rgb.global != nil {
		return rgb.global.SetBrightness(brightness)
	}
	return fmt.Errorf("ENOSYS")
}

func (rgb *rgbled) Off() error {
	if rgb.global != nil {
		return rgb.global.Off()
	}

	err := rgb.red.Off()
	if err != nil {
		return err
	}
	err = rgb.green.Off()
	if err != nil {
		return err
	}
	err = rgb.blue.Off()
	if err != nil {
		return err
	}

	return nil
}

func (rgb *rgbled) SetTrigger(trigger Trigger) error {
	if rgb.global != nil {
		err := rgb.global.SetTrigger(trigger)
		if err == nil {
			rgb.trigger = trigger
		}
		return err
	}

	rtrig := rgb.red.trigger
	err := rgb.red.Off()
	if err != nil {
		return err
	}

	gtrig := rgb.green.trigger
	err = rgb.green.Off()
	if err != nil {
		rgb.red.SetTrigger(rtrig)
		return err
	}

	err = rgb.blue.Off()
	if err != nil {
		rgb.green.SetTrigger(gtrig)
		rgb.red.SetTrigger(rtrig)
		return err
	}

	rgb.trigger = trigger

	return nil
}

func (rgb *rgbled) GetTrigger() Trigger {
	return rgb.trigger
}

func (rgb *rgbled) SetColor(c color.Color) error {
	if c == rgb.color {
		return nil
	}

	r, g, b, a := c.RGBA()
	if a != 0 && a != 0xffff {
		r = (r << 16) / a
		g = (g << 16) / a
		b = (b << 16) / a
	}

	rf := float32(r) / 65535.0
	gf := float32(g) / 65535.0
	bf := float32(b) / 65535.0

	err1 := rgb.red.SetBrightness(rf)
	err2 := rgb.green.SetBrightness(gf)
	err3 := rgb.blue.SetBrightness(bf)

	if err1 != nil || err2 != nil || err3 != nil {
		return fmt.Errorf("Failed")
	}

	rgb.color = c

	return nil
}

func (rgb *rgbled) GetColor() color.Color {
	return rgb.color
}
