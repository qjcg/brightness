// Set backlight brightness on Linux via sysfs.
// User running this command must have write access to FCtl (default: 0644/root:root).
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	FCtl = "/sys/class/backlight/intel_backlight/brightness"
	FMax = "/sys/class/backlight/intel_backlight/max_brightness"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [LEVEL 0-100 | INCREMENT +/-N]\n", os.Args[0])
}

func main() {
	fCtl, err := os.Open(FCtl)
	if err != nil {
		log.Fatal(err)
	}
	defer fCtl.Close()

	fMax, err := os.Open(FMax)
	if err != nil {
		log.Fatal(err)
	}
	defer fMax.Close()

	bl, err := NewBacklight(fCtl, fMax)

	// Act based on the number of arguments, excluding this command.
	switch len(os.Args) - 1 {

	// Print current brightness to stdout.
	case 0:
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
