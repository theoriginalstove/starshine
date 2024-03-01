package app

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/hsluv/hsluv-go"
)

type mode string

const (
	Static  = "STATIC"
	RGBWave = "RGBWave"
	RGBFade = "RGBFade"
)

func (h *Handler) static(r, g, b uint8) error {
	for i := 0; i < len(h.led.ws.Leds(0)); i++ {
		h.led.ws.Leds(0)[i] = rgbToColor(r, g, b)
	}

	if err := h.led.ws.Render(); err != nil {
		return err
	}
	return nil
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

func (h *Handler) rgbWave() error {
	s, err := createHSVRainbowSlice()
	if err != nil {
		return err
	}

	for {
		select {
		case <-h.quit:
			break
		default:
			for i := 0; i < ledCounts; i++ {
				h.led.ws.Leds(0)[i] = s[i]
				err = h.led.ws.Render()
				if err != nil {
					return err
				}
			}

			s = shiftUint32Slice(s)

			time.Sleep(sleepTime * time.Millisecond)
		}
		break
	}
	return nil
}

func (h *Handler) rainbowHSVToRGBFade() error {
	for {
		select {
		case <-h.quit:
			break
		default:
			for i := 0; i <= 360; i++ {
				hs := hsluv.HsluvToHex(float64(i), 100, 80)
				n, err := strconv.ParseUint(strings.TrimPrefix(hs, "#"), 16, 32)
				if err != nil {
					return err
				}

				err = h.renderAllHex(uint32(n))
				if err != nil {
					return err
				}
				time.Sleep(sleepTime * time.Millisecond)
			}
		}
		break
	}
	return nil
}

func (h *Handler) renderAll() error {
	for i := 0; i < len(h.led.ws.Leds(0)); i++ {
		h.led.ws.Leds(0)[i] = rgbToColor(h.led.clr.r, h.led.clr.g, h.led.clr.b)
	}

	if err := h.led.ws.Render(); err != nil {
		return err
	}

	return nil
}

func (h *Handler) renderAllHex(color uint32) error {
	for i := 0; i < len(h.led.ws.Leds(0)); i++ {
		h.led.ws.Leds(0)[i] = color
	}

	if err := h.led.ws.Render(); err != nil {
		return err
	}

	return nil
}

func shiftUint32Slice(s []uint32) []uint32 {
	x := s[0] // get the 0 index element from slice
	s = s[1:]
	s = append(s, x)
	return s
}

func rgbToColor(r uint8, g uint8, b uint8) uint32 {
	return uint32(uint32(r)<<16 | uint32(g)<<8 | uint32(b))
}

func getRGBColor(c string) ([]uint8, error) {
	vals := strings.Split(c, ",")
	if len(vals) != 3 {
		return []uint8{0, 0, 0}, errors.New("not a valid RGB color")
	}
	r, err := strconv.Atoi(vals[0])
	if err != nil {
		return []uint8{0, 0, 0}, errors.New("not a valid RGB color")
	}
	g, err := strconv.Atoi(vals[1])
	if err != nil {
		return []uint8{0, 0, 0}, errors.New("not a valid RGB color")
	}
	b, err := strconv.Atoi(vals[2])
	if err != nil {
		return []uint8{0, 0, 0}, errors.New("not a valid RGB color")
	}

	if r < 0 || r > 255 {
		r = 0
	}

	if g < 0 || g > 255 {
		g = 0
	}

	if b < 0 || b > 255 {
		b = 0
	}
	return []uint8{uint8(r), uint8(g), uint8(b)}, nil
}
