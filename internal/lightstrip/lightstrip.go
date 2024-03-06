package lightstrip

import (
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hsluv/hsluv-go"
	led "github.com/rpi-ws281x/rpi-ws281x-go"
)

type mode string

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

type Lightstrip struct {
	ws          *led.WS2811
	quitWave    chan bool
	quitFade    chan bool
	currentMode string
	lastW       bool
	lastF       bool
	mu          sync.Mutex
}

func (ls *Lightstrip) Init(test bool) error {
	ledOpts := led.DefaultOptions
	ledOpts.Channels[0].Brightness = brightness
	ledOpts.Channels[0].LedCount = ledCounts
	ledOpts.Channels[0].GpioPin = gpioPin
	ledOpts.Frequency = freq
	ws2811, err := led.MakeWS2811(&ledOpts)
	if err != nil {
		return err
	}
	ls.mu = sync.Mutex{}

	ls.ws = ws2811
	ls.quitFade = make(chan bool)
	ls.quitWave = make(chan bool)

	slog.Info("initializing ws2811")
	err = ls.ws.Init()
	if err != nil {
		return err
	}

	if test {
		slog.Info("testing lights")
		ls.Static("#FFFFFF")
		time.Sleep(1 * time.Second)
		slog.Info("lights off")
		ls.Static("#000000")
	}
	return nil
}

func (ls *Lightstrip) Static(hex string) error {
	ls.Cancel()

	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.currentMode = "static"
	h, err := strconv.ParseUint(strings.TrimPrefix(hex, "#"), 16, 32)
	if err != nil {
		return err
	}
	for i := 0; i < len(ls.ws.Leds(0)); i++ {
		ls.ws.Leds(0)[i] = uint32(h)
	}

	leds := ls.ws.Leds(0)
	ls.ws.SetLedsSync(0, leds)

	if err := ls.ws.Render(); err != nil {
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

func (ls *Lightstrip) RgbWave() {
	ls.mu.Lock()
	ls.currentMode = "wave"
	ls.quitWave = make(chan bool)
	ls.mu.Unlock()
	slog.Info("starting wave")

	s, err := createHSVRainbowSlice()
	if err != nil {
		slog.Error("unable to create HSV Rainbow Slice", slog.Any("err", err))
	}

	ticker := time.NewTicker(1 * time.Microsecond)

loop:
	for {
		select {
		case t := <-ticker.C:
			slog.Debug("tick", "tock", t)
			for i := 0; i < ledCounts; i++ {
				ls.ws.Leds(0)[i] = s[i]
				leds := ls.ws.Leds(0)
				err = ls.ws.SetLedsSync(0, leds)
				err = ls.ws.Render()
				if err != nil {
					slog.Error("unable to render rgb wave", slog.Any("err", err))
				}
			}
			s = shiftUint32Slice(s)
		case q := <-ls.quitWave:
			slog.Info("stopping wave", slog.Bool("q", q))
			ticker.Stop()
			break loop
		}
	}
}

func (ls *Lightstrip) RainbowHSVToRGBFade() {
	ls.mu.Lock()
	ls.currentMode = "fade"
	ls.mu.Unlock()
	color := 0

	ticker := time.NewTicker(1 * time.Millisecond)
	select {
	case <-ls.quitFade:
		ticker.Stop()
		break
	case t := <-ticker.C:
		if color > 360 {
			color = 0
		}
		slog.Info("tick", slog.Any("ticker", t))
		hs := hsluv.HsluvToHex(float64(color), 100, 80)
		n, err := strconv.ParseUint(strings.TrimPrefix(hs, "#"), 16, 32)
		if err != nil {
			slog.Error("unable to parse uint", slog.Any("err", err))
		}

		err = ls.renderAllHex(uint32(n))
		if err != nil {
			slog.Error("unable to render all hex in fade")
		}
		color++
		time.Sleep(1 * time.Millisecond)

	}
}

func (ls *Lightstrip) renderAllHex(color uint32) error {
	for i := 0; i < len(ls.ws.Leds(0)); i++ {
		ls.ws.Leds(0)[i] = color
	}

	if err := ls.ws.Render(); err != nil {
		return err
	}

	return nil
}

// shift the leds by 10
func shiftUint32Slice(s []uint32) []uint32 {
	x := s[:5] // get the first 10 index element from slice
	s = s[5:]
	s = append(s, x...)
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

func (ls *Lightstrip) Off() error {
	return ls.Static("#000000")
}

func (ls *Lightstrip) Cancel() error {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	mode := ls.currentMode
	go func() {
		switch mode {
		case "fade":
			slog.Info("sending cancel signal - fade")
			ls.quitFade <- !ls.lastF
			ls.lastF = !ls.lastF
		case "wave":
			slog.Info("sending cancel signal - wave")
			ls.quitWave <- !ls.lastW
			ls.lastW = !ls.lastW
		}
	}()
	return nil
}
