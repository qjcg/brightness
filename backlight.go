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

func NewBacklight(cur, max io.Reader) (*Backlight, error) {
	var bl Backlight
	var err error

	bl.CurrentBrightness, err = getVal(cur)
	if err != nil {
		return &bl, err
	}

	bl.MaxBrightness, err = getVal(max)
	if err != nil {
		return &bl, err
	}

	return &bl, nil
}

// getVal retrieves a float value from the provided io.Reader.
func getVal(r io.Reader) (float64, error) {
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

// Set writes a backlight brightness value to the provided io.Writer.
func (bl *Backlight) Set(w io.Writer, pct float64) error {
	b := strconv.AppendFloat(
		[]byte(nil),
		fromPct(bl.MaxBrightness, pct),
		'f', 0, 64)
	_, err := w.Write(b)
	return err
}

// SetIncr writes an increment of backlight brightness to the provided io.ReadWriter.
// For example: -5 -> -5% of max brightness
func (bl *Backlight) SetIncr(w io.Writer, pct float64) error {
	b := strconv.AppendFloat(
		[]byte(nil),
		bl.CurrentBrightness+fromPct(bl.MaxBrightness, pct),
		'f', 0, 64)
	_, err := w.Write(b)
	return err
}

// Set writes backlight brightness to the provided io.ReadWriter.
// pct: overall brightness level, or increment, expressed as a percentage
// 	- (ex: 30 -> 30% of max brightness)
// 	- (ex: -5 -> -5% of max brightness)
// incr: indicates whether the pct value provided is an increment.

// toPct returns a percent value for brightness, given a brightness and
// max brightness value in "brightness units".
// E.g. toPct(4437.0, 1092.0) -> 25 (%)
func toPct(max, b float64) float64 {
	return b / max * 100.0
}

// fromPct takes a percent brightness value and max value in "brightness
// units", and returns the current brightness value in "brightness units".
// E.g. fromPct(25.0, 4437.0) -> 1109.0
func fromPct(max, pct float64) float64 {
	return pct / 100.0 * max
}
