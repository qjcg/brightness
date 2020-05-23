// Set backlight brightness on Linux via sysfs.
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// FIXME: Don't hard code these, accept them from config/environment.
var (
	FCtl = "/sys/class/backlight/intel_backlight/brightness"
	FMax = "/sys/class/backlight/intel_backlight/max_brightness"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [LEVEL 0-100 | INCREMENT +/-N]\n", os.Args[0])
}

func main() {

	// Open backlight brightness and max_brightness files.
	fCtl, err := os.OpenFile(FCtl, os.O_RDWR, 0)
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
		fmt.Printf("%0.0f\n", bl.Percent())

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
			if err := bl.SetIncr(fCtl, level); err != nil {
				log.Fatalf("Error setting brightness to level: %v\n%v\n", level, err)
			}
			os.Exit(0)
		}

		// Provided level is an absolute percentage (ex: 25%).
		if err := bl.Set(fCtl, level); err != nil {
			log.Fatalf("Error setting brightness to level: %v\n%v\n", level, err)
		}

	default:
		usage()
	}
}
