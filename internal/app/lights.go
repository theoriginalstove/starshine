package app

import (
	"strconv"
	"strings"
	"time"

	"github.com/hsluv/hsluv-go"
)

type mode string

const (
	Static  = "STATIC"
	RGBWave = "RGBWave"
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

func shiftUint32Slice(s []uint32) []uint32 {
	x := s[0] // get the 0 index element from slice
	s = s[1:]
	s = append(s, x)
	return s
}
