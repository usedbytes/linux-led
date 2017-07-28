package led

import (
	"bytes"
	"io/ioutil"
	"strings"
	"strconv"
)

type LinuxLED interface {
	SetBrightness(brightness float32) error
	Off() error
	SetTrigger(trigger Trigger) error
	GetTrigger() Trigger
}

type linuxled struct {
	syspath string
	maxBrightness float32

	triggers []Trigger
	trigger Trigger
}

type Trigger string

const (
	TriggerNone      = "none"
	TriggerPanic     = "panic"
	TriggerDisk      = "disk-activity"
	TriggerHeartbeat = "heartbeat"
)

func NewLED(syspath string) (LinuxLED, error) {
	led := linuxled{
		syspath: syspath,
	}

	b, err := ioutil.ReadFile(syspath + "/max_brightness")
	if err != nil {
		return nil, err
	}
	mb, err := strconv.Atoi(string(bytes.TrimSpace(b)))
	if err != nil {
		return nil, err
	}
	led.maxBrightness = float32(mb)

	b, err = ioutil.ReadFile(syspath + "/trigger")
	if err != nil {
		return nil, err
	}
	triggers := strings.Split(string(b), " ")
	for _, t := range triggers {
		if t[0] == '[' && t[len(t) - 1] == ']' {
			t = strings.Trim(t, "[]")
			led.trigger = Trigger(t)
		}
		led.triggers = append(led.triggers, Trigger(t))
	}

	return &led, nil
}

func (led *linuxled) SetBrightness(brightness float32) error {
	str := strconv.Itoa(int(brightness * led.maxBrightness))
	return ioutil.WriteFile(led.syspath + "/brightness", []byte(str), 0)
}

func (led *linuxled) Off() error {
	return led.SetBrightness(0)
}

func (led *linuxled) SetTrigger(trigger Trigger) error {
	err := ioutil.WriteFile(led.syspath + "/trigger", []byte(trigger), 0)
	if err == nil {
		led.trigger = trigger
	}
	return err
}

func (led *linuxled) GetTrigger() Trigger {
	return led.trigger
}
