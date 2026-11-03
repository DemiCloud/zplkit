package main

import (
	"fmt"
	"strconv"
	"time"
)

// adminActions maps an admin action name (the first positional argument,
// in the slot a FILE path otherwise occupies) to a builder that returns
// the exact ZPL/SGD bytes to send. Anything not in this map falls back
// to admin's original "read this file and send it" behavior — see
// adminCmd in admin.go.
var adminActions = map[string]func(args []string) (string, error){
	"clock":     buildClockZPL,
	"darkness":  buildDarknessZPL,
	"calibrate": buildCalibrateZPL,
	"speed":     buildSpeedZPL,
}

// buildClockZPL sets the printer's real-time clock to the current time
// via ^ST (Set Time), in 24-hour/military format so no AM/PM conversion
// is needed. Requires the printer to have the RTC option installed;
// unlike darkness/speed below, the clock is its own battery-backed
// hardware and isn't lost on power-off, so no ^JUS is needed here.
func buildClockZPL(args []string) (string, error) {
	if len(args) != 0 {
		return "", fmt.Errorf("clock takes no arguments")
	}
	now := time.Now()
	return fmt.Sprintf("^XA^ST%d,%d,%d,%d,%d,%d,M^XZ",
		now.Month(), now.Day(), now.Year(), now.Hour(), now.Minute(), now.Second()), nil
}

// buildDarknessZPL sets print darkness via ~SD (0-30). Per Zebra's docs,
// this kind of persistent parameter is lost on power-off unless saved
// with ^JUS, so that's appended here rather than left to the caller.
func buildDarknessZPL(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("darkness takes exactly one argument: a level from 0 to 30")
	}
	n, err := strconv.Atoi(args[0])
	if err != nil || n < 0 || n > 30 {
		return "", fmt.Errorf("invalid darkness %q: must be an integer from 0 to 30", args[0])
	}
	return fmt.Sprintf("~SD%02d^XA^JUS^XZ", n), nil
}

// buildCalibrateZPL triggers media sensor (label length/gap)
// calibration via ~JC. No arguments, and no ^JUS: calibration writes
// directly to the printer's stored settings.
func buildCalibrateZPL(args []string) (string, error) {
	if len(args) != 0 {
		return "", fmt.Errorf("calibrate takes no arguments")
	}
	return "~JC", nil
}

// buildSpeedZPL sets print speed via ^PR (1-14 inches/second), saved
// with ^JUS for the same reason as buildDarknessZPL.
func buildSpeedZPL(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("speed takes exactly one argument: inches per second from 1 to 14")
	}
	n, err := strconv.Atoi(args[0])
	if err != nil || n < 1 || n > 14 {
		return "", fmt.Errorf("invalid speed %q: must be an integer from 1 to 14", args[0])
	}
	return fmt.Sprintf("^XA^PR%d^JUS^XZ", n), nil
}
