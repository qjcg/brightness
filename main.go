// Set backlight brightness on Linux via sysfs.
// User running this command must have write access to FCtl (default: 0644/root:root).
package main // import "github.com/qjcg/brightness"

import (
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
// pct: overall brightness level, or increment, expressed as a percentage
// 	- (ex: 30 -> 30% of max brightness)
// 	- (ex: -5 -> -5% of max brightness)
// incr: indicates whether the pct value provided is an increment.
func Set(fctl string, max, pct float64, incr bool) error {
	var b []byte

	if incr {
		cur, err := Get(fctl)
		if err != nil {
			return err
		}
		b = strconv.AppendFloat([]byte(nil), cur+fromPct(max, pct), 'f', 0, 64)
	} else {
		b = strconv.AppendFloat([]byte(nil), fromPct(max, pct), 'f', 0, 64)
	}

	if err := ioutil.WriteFile(fctl, b, 0644); err != nil {
		return err
	}
	return nil
}

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

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [LEVEL 0-100 | INCREMENT +/-N]\n", os.Args[0])
}

func main() {

	// Get maximum brightness.
	max, err := Get(FMax)
	if err != nil {
		log.Fatal(err)
	}

	// Act based on the number of arguments, excluding this command.
	switch len(os.Args) - 1 {

	// Print current brightness to stdout.
	case 0:
		b, err := Get(FCtl)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%0.0f\n", toPct(max, b))

	// Set brightness to provided pct value.
	case 1:

		// If argument is help, print usage message.
		switch os.Args[1] {
		case "-h", "-help", "--help", "help":
			usage()
			os.Exit(0)
		}

		level, err := strconv.ParseFloat(os.Args[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		// Provided level is an increment (ex: +5, or -20)
		if strings.ContainsAny(string(os.Args[1][0]), "+-") {
			if err := Set(FCtl, max, level, true); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}

		// Provided level is an absolute percentage (ex: 25%).
		if err := Set(FCtl, max, level, false); err != nil {
			log.Fatal(err)
		}

	default:
		usage()
	}
}
