package main

import (
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

// Backlight represents a screen backlight.
type Backlight struct {
	CurrentBrightness float64
	MaxBrightness     float64
}

// NewBacklight returns a new Backlight value.
func NewBacklight(cur, max io.Reader) (*Backlight, error) {
	var bl Backlight
	var err error

	bl.CurrentBrightness, err = get(cur)
	if err != nil {
		return &bl, err
	}

	bl.MaxBrightness, err = get(max)
	if err != nil {
		return &bl, err
	}

	return &bl, nil
}

// get retrieves a float value from the provided io.Reader.
func get(r io.Reader) (float64, error) {
	var val float64

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return val, err
	}

	s := strings.TrimSpace(string(b))
	val, err = strconv.ParseFloat(s, 64)
	if err != nil {
		return val, err
	}

	return val, nil
}

// Set writes a backlight brightness value (in percent) to the provided io.Writer.
func (bl *Backlight) Set(w io.Writer, pct float64) error {
	b := strconv.AppendFloat([]byte(nil), bl.PercentToBrightness(pct), 'f', 0, 64)
	_, err := w.Write(b)
	return err
}

// SetIncr writes an increment of backlight brightness to the provided io.ReadWriter.
// For example: -5 -> -5% of max brightness
func (bl *Backlight) SetIncr(w io.Writer, pct float64) error {
	totalBrightness := bl.CurrentBrightness + bl.PercentToBrightness(pct)
	b := strconv.AppendFloat([]byte(nil), totalBrightness, 'f', 0, 64)
	_, err := w.Write(b)
	return err
}

// Percent returns a percent value for brightness.
func (bl *Backlight) Percent() float64 {
	return bl.CurrentBrightness / bl.MaxBrightness * 100
}

// PercentToBrightness takes a percent value and returns the percent of max in "brightness units".
func (bl *Backlight) PercentToBrightness(pct float64) float64 {
	return pct / 100 * bl.MaxBrightness
}
