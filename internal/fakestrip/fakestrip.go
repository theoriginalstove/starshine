package fakestrip

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
