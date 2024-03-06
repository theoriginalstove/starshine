package fakestrip

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// Fakestrip is used to mock the LED Lightstrip functionality for local development
// of the UI but not for testing actual LED functions
type Fakestrip struct{}

func (f *Fakestrip) Init(test bool) error {
	return nil
}

func (f *Fakestrip) Cancel() error {
	return nil
}

func (f *Fakestrip) RgbWave() {
}

func (f *Fakestrip) RainbowHSVToRGBFade() {
}

func (f *Fakestrip) Static(hex string) error {
	return nil
}

func (f *Fakestrip) Off() error {
	return f.Static("")
}

func (f *Fakestrip) Mode() string {
	mode := fmt.Sprintf("fake mode - %s", randStr(8))
	return mode
}

// RandStr generates a random string of length len unless len is less than 8
func randStr(len int) string {
	if len < 8 {
		len = 8
	}
	b := make([]byte, len)
	if _, err := rand.Read(b); err != nil {
		// rand.Read should never error unless we run out of entropy
		panic(err)
	}

	b = b[:4]
	return hex.EncodeToString(b)
}
