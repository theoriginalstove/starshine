package main

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/hsluv/hsluv-go"

	led "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	brightness = 64
	ledCounts  = 9
	gpioPin    = 18
	freq       = 700000
	sleepTime  = 1
)

type color struct {
	r uint8
	g uint8
	b uint8
}

type ws struct {
	ws2811 *led.WS2811
	clr    *color
}

func (w *ws) init() error {
	err := w.ws2811.Init()
	if err != nil {
		return err
	}
	return nil
}

func (ws *ws) renderAll() error {
	for i := 0; i < len(ws.ws2811.Leds(0)); i++ {
		ws.ws2811.Leds(0)[i] = rgbToColor(ws.clr.r, ws.clr.g, ws.clr.b)
	}

	if err := ws.ws2811.Render(); err != nil {
		return err
	}

	return nil
}

func (ws *ws) renderAllHex(color uint32) error {
	for i := 0; i < len(ws.ws2811.Leds(0)); i++ {
		ws.ws2811.Leds(0)[i] = color
	}

	if err := ws.ws2811.Render(); err != nil {
		return err
	}

	return nil
}

func (ws *ws) close() {
	ws.ws2811.Fini()
}

func rgbToColor(r uint8, g uint8, b uint8) uint32 {
	return uint32(uint32(r)<<16 | uint32(g)<<8 | uint32(b))
}

func (ws *ws) rgbLinIncRRender() error {
	for ws.clr.r = 0; ws.clr.r < 255; ws.clr.r++ {
		err := ws.renderAll()
		if err != nil {
			return err
		}
		time.Sleep(sleepTime * time.Millisecond)
	}

	return nil
}

func (ws *ws) rgbLinDecRRender() error {
	for ws.clr.r = 255; ws.clr.r > 0; ws.clr.r-- {
		err := ws.renderAll()
		if err != nil {
			return err
		}
		time.Sleep(sleepTime * time.Millisecond)
	}

	return nil
}

func (ws *ws) rgbLinIncGRender() error {
	for ws.clr.g = 0; ws.clr.g < 255; ws.clr.g++ {
		err := ws.renderAll()
		if err != nil {
			return err
		}
		time.Sleep(sleepTime * time.Millisecond)
	}

	return nil
}

func (ws *ws) rgbLinDecGRender() error {
	for ws.clr.g = 255; ws.clr.g > 0; ws.clr.g-- {
		err := ws.renderAll()
		if err != nil {
			return err
		}
		time.Sleep(sleepTime * time.Millisecond)
	}

	return nil
}

func (ws *ws) rgbLinIncBRender() error {
	for ws.clr.b = 0; ws.clr.b < 255; ws.clr.b++ {
		err := ws.renderAll()
		if err != nil {
			return err
		}
		time.Sleep(sleepTime * time.Millisecond)
	}

	return nil
}

func (ws *ws) rgbLinDecBRender() error {
	for ws.clr.b = 255; ws.clr.b > 0; ws.clr.b-- {
		err := ws.renderAll()
		if err != nil {
			return err
		}
		time.Sleep(sleepTime * time.Millisecond)
	}

	return nil
}

func (ws *ws) rainbowRGB() error {
	ws.clr.r = 255
	err := ws.renderAll()
	if err != nil {
		panic(err)
	}
	time.Sleep(sleepTime * time.Millisecond)

	for {
		err = ws.rgbLinIncGRender()
		if err != nil {
			return err
		}

		err = ws.rgbLinDecRRender()
		if err != nil {
			return err
		}

		err = ws.rgbLinIncBRender()
		if err != nil {
			return err
		}

		err = ws.rgbLinDecGRender()
		if err != nil {
			return err
		}

		err = ws.rgbLinIncRRender()
		if err != nil {
			return err
		}

		err = ws.rgbLinDecBRender()
		if err != nil {
			return err
		}
	}
}

func createHSVRainbowSlice() ([]uint32, error) {
	s := []uint32{}

	for i := 0; i <= 360; i++ {
		h := hsluv.HsluvToHex(float64(i), 100, 80)
		n, err := strconv.ParseUint(strings.TrimPrefix(h, "#"), 16, 32)
		if err != nil {
			return nil, err
		}

		s = append(s, uint32(n))
	}

	return s, nil
}

func shiftUint32Slice(s []uint32) []uint32 {
	x := s[0] // get the 0 index element from slice
	s = s[1:]
	s = append(s, x)
	return s
}

func (ws *ws) rainbowHSVToRGBFade() error {
	for {
		for i := 0; i <= 360; i++ {
			h := hsluv.HsluvToHex(float64(i), 100, 80)
			n, err := strconv.ParseUint(strings.TrimPrefix(h, "#"), 16, 32)
			if err != nil {
				return err
			}

			err = ws.renderAllHex(uint32(n))
			if err != nil {
				return err
			}
			time.Sleep(sleepTime * time.Millisecond)
		}
	}
}

func (ws *ws) rainbowHSVToRGBWave() error {
	s, err := createHSVRainbowSlice()
	if err != nil {
		return err
	}

	counter := 0

	for {
		for i := 0; i < ledCounts; i++ {
			ws.ws2811.Leds(0)[i] = s[i]
			ws.ws2811.Render()
		}

		s = shiftUint32Slice(s)

		counter++
		time.Sleep(sleepTime * time.Millisecond)
		if counter == 50 {
			break
		}
	}
	return nil
}

func main() {
	slog.Info("starting up app")
	opts := led.DefaultOptions
	opts.Channels[0].Brightness = brightness
	opts.Channels[0].LedCount = ledCounts
	opts.Channels[0].GpioPin = gpioPin
	opts.Frequency = freq

	ws2811, err := led.MakeWS2811(&opts)
	if err != nil {
		panic(err)
	}

	ws := ws{
		ws2811: ws2811,
		clr:    &color{},
	}

	err = ws.init()
	if err != nil {
		if err != nil {
			panic(err)
		}
	}

	defer ws.close()
	err = ws.rainbowHSVToRGBWave()
	if err != nil {
		panic(err)
	}

	slog.Info("shutting down")
}
