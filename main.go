// Set backlight brightness on a Linux system via sysfs.
// User running this command must have write access to ControlFile (default: 0644/root:root).
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

	// Maximum brightness (float64).
	Max = Get(FMax)
)

// Get retrieves a brightness value from the provided control file. Result
// is a float64 representing the brightness in arbitrary "brightness units"
// relative to Max.
func Get(fctl string) float64 {
	var brightness float64

	b, err := ioutil.ReadFile(fctl)
	check(err)

	s := strings.TrimSpace(string(b))
	brightness, err = strconv.ParseFloat(s, 64)
	check(err)

	return brightness
}

// Set writes backlight brightness to the provided control file.
// levelpct: overall brightness level expressed as a percentage (ex: 30 -> 30% of max brightness)
func Set(fctl string, pct float64) {
	b := strconv.AppendFloat([]byte(nil), FromPct(pct, Max), 'f', 0, 64)
	err := ioutil.WriteFile(fctl, b, 0644)
	check(err)
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

	if len(os.Args) != 2 {
		fmt.Printf("%0.0f\n", ToPct(Get(FCtl), Max))
		os.Exit(0)
	}

	level, err := strconv.ParseFloat(os.Args[1], 64)
	check(err)
	Set(FCtl, level)
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
