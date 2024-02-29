package app

import (
	"log/slog"
	"net/http"

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
	led *LED
}

func NewHandler(led *LED) *Handler {
	return &Handler{led: led}
}

func (h *Handler) on(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "turning leds on")
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

func rgbToColor(r uint8, g uint8, b uint8) uint32 {
	return uint32(uint32(r)<<16 | uint32(g)<<8 | uint32(b))
}
