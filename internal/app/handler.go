package app

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

const (
	brightness = 128
	ledCounts  = 300
	gpioPin    = 18
	freq       = 800000
	sleepTime  = 100
)

type color struct {
	r uint8
	g uint8
	b uint8
}

type Handler struct {
	led lighter
}

func NewHandler(l lighter) *Handler {
	return &Handler{led: l}
}

type lightsRequest struct {
	Mode     string `json:"mode"`
	RGBColor string `json:"rgbColor"`
}

func (h *Handler) reset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "resetting lighter")
	err := h.led.Init(true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) on(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "turning leds on")
	err := h.led.Cancel()
	if err != nil {
		slog.ErrorContext(ctx, "error encountered cancelling lights", slog.Any("err", err))
		w.WriteHeader(500)
		return
	}

	err = h.led.Static("#FFFFFF")
	if err != nil {
		slog.ErrorContext(ctx, "error encountered cancelling lights", slog.Any("err", err))
		w.WriteHeader(500)
		return
	}
	slog.InfoContext(ctx, "leds turned on")
	w.WriteHeader(200)
}

func (h *Handler) off(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.InfoContext(ctx, "turning leds off")
	err := h.led.Cancel()
	if err != nil {
		slog.ErrorContext(ctx, "error encountered cancelling lights", slog.Any("err", err))
		w.WriteHeader(500)
		return
	}

	err = h.led.Off()
	if err != nil {
		slog.ErrorContext(ctx, "error encountered turning off lights", slog.Any("err", err))
		w.WriteHeader(500)
		return
	}
	slog.InfoContext(ctx, "leds turned off")
	w.WriteHeader(200)
}

func (h *Handler) setMode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	err = h.led.Cancel()
	if err != nil {
		slog.ErrorContext(ctx, "error encountered cancelling lights", slog.Any("err", err))
		w.WriteHeader(500)
		return
	}

	switch m := req.Mode; strings.ToLower(m) {
	case "static":
		err = h.led.Static(req.RGBColor)
		if err != nil {
			slog.ErrorContext(ctx, "unable to set light mode", slog.Any("err", err))
			w.WriteHeader(500)
			return
		}
	case "warm":
		slog.InfoContext(ctx, "setting to warm")
		err = h.led.Static("#ffb348")
		if err != nil {
			slog.ErrorContext(ctx, "unable to set light mode", slog.Any("err", err), slog.String("mode", m))
			w.WriteHeader(500)
			return
		}
	case "rgbwave":
		slog.InfoContext(ctx, "setting to rgbwave")
		go h.led.RgbWave()
	case "rgbfade":
		slog.InfoContext(ctx, "setting to rgbfade")
		go h.led.RainbowHSVToRGBFade()
	default:
		slog.InfoContext(ctx, "defaulting to off")
		err = h.led.Static("#000000")
		if err != nil {
			slog.ErrorContext(ctx, "unable to set light mode", slog.Any("err", err), slog.String("case", "default"), slog.String("mode", m))
			w.WriteHeader(500)
			return
		}
	}
	if err != nil {
		slog.ErrorContext(ctx, "unable to set light mode", slog.Any("err", err))
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

type lighter interface {
	Init(test bool) error
	Cancel() error
	RgbWave()
	RainbowHSVToRGBFade()
	Static(hex string) error
	Off() error
}
