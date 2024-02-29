package app

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	led "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	brightness = 128
	ledCounts  = 300
	gpioPin    = 18
	freq       = 800000
	sleepTime  = 1
)

type color struct {
	r uint8
	g uint8
	b uint8
}

type LED struct {
	ws  *led.WS2811
	clr *color
}

func (l *LED) init() error {
	err := l.ws.Init()
	if err != nil {
		return err
	}
	return nil
}

type Handler struct {
	led  *LED
	quit chan bool
}

func NewHandler(led *LED) *Handler {
	return &Handler{led: led, quit: make(chan bool)}
}

type lightsRequest struct {
	Mode     string `json:"mode"`
	RGBColor string `json:"rgbColor"`
}

func (h *Handler) on(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "turning leds on")
	go func() {
		h.quit <- false
	}()
	for i := 0; i < len(h.led.ws.Leds(0)); i++ {
		h.led.ws.Leds(0)[i] = rgbToColor(254, 196, 127)
	}

	if err := h.led.ws.Render(); err != nil {
		w.WriteHeader(500)
		slog.ErrorContext(ctx, "unable to set LEDs", slog.Any("err", err))
		return
	}

	slog.InfoContext(ctx, "leds turned on")
	w.WriteHeader(200)
}

func (h *Handler) off(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "turning leds off")
	for i := 0; i < len(h.led.ws.Leds(0)); i++ {
		h.led.ws.Leds(0)[i] = rgbToColor(0, 0, 0)
	}

	if err := h.led.ws.Render(); err != nil {
		w.WriteHeader(500)
		slog.ErrorContext(ctx, "unable to set LEDs", slog.Any("err", err))
		return
	}
	slog.InfoContext(ctx, "leds turned off")
	w.WriteHeader(200)
}

func (h *Handler) setMode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	go func() {
		h.quit <- false
	}()
	bts, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("internal server error"))
		w.WriteHeader(500)
		return
	}
	req := &lightsRequest{}
	err = json.Unmarshal(bts, req)
	if err != nil {
		slog.ErrorContext(ctx, "unable to unmarshal request", slog.Any("err", err))
		w.WriteHeader(500)
		return
	}

	switch m := req.Mode; strings.ToLower(m) {
	case "static":
		rgb, err := getRGBColor(req.RGBColor)
		if err != nil {
			slog.ErrorContext(ctx, "unable to set light mode", slog.Any("err", err))
			w.WriteHeader(500)
			return
		}
		err = h.static(rgb[0], rgb[1], rgb[2])
	case "warm":
		err = h.static(254, 169, 72)
	case "rgbwave":
		err = h.rgbWave()
	default:
		err = h.static(0, 0, 0)
	}
	if err != nil {
		slog.ErrorContext(ctx, "unable to set light mode", slog.Any("err", err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
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
