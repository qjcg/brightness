// Set backlight brightness on a Linux system via sysfs.
// User running this command must have write access to FCtl (default: 0644/root:root).
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	FCtl = "/sys/class/backlight/intel_backlight/brightness"
	FMax = "/sys/class/backlight/intel_backlight/max_brightness"
)

// Get retrieves a brightness value from the provided control file. Result
// is a float64 representing the brightness in arbitrary "brightness units".
func Get(fctl string) (float64, error) {
	var brightness float64

	b, err := ioutil.ReadFile(fctl)
	if err != nil {
		return brightness, err
	}

	s := strings.TrimSpace(string(b))
	brightness, err = strconv.ParseFloat(s, 64)
	if err != nil {
		return brightness, err
	}

	return brightness, nil
}

// Set writes backlight brightness to the provided control file.
// pct: overall brightness level expressed as a percentage (ex: 30 -> 30% of max brightness)
func Set(fctl string, pct, max float64) error {
	b := strconv.AppendFloat([]byte(nil), FromPct(pct, max), 'f', 0, 64)
	if err := ioutil.WriteFile(fctl, b, 0644); err != nil {
		return err
	}
	return nil
}

// ToPct returns a percent value for brightness, given a brightness and
// max brightness value in "brightness units".
// E.g. ToPct(1092.0, 4437.0) -> 25 (%)
func ToPct(b, max float64) float64 {
	return b / max * 100.0
}

// FromPct takes a percent brightness value and max value in "brightness
// units", and returns the current brightness value in "brightness units".
// E.g. FromPct(25.0, 4437.0) -> 1109.0
func FromPct(pct, max float64) float64 {
	return pct / 100.0 * max
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [level 1-100]\n", os.Args[0])
	}
	flag.Parse()

	// Get maximum brightness.
	max, err := Get(FMax)
	if err != nil {
		log.Fatal(err)
	}

	switch flag.NArg() {

	// Print current brightness to stdout.
	case 0:
		b, err := Get(FCtl)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%0.0f\n", ToPct(b, max))

	// Set brightness to provided pct value.
	case 1:
		level, err := strconv.ParseFloat(os.Args[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		if err := Set(FCtl, level, max); err != nil {
			log.Fatal(err)
		}

	default:
		flag.Usage()
	}
}
