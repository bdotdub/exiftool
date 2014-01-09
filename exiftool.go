package exiftool

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	gpsPrecisionFmt = "%.15f"
)

type Exif struct {
	DateTimeOriginal string
	GPS              struct {
		Latitude    float64
		Longitude   float64
	}
}

func DecodeFileAtPath(p string) (*Exif, error) {
	var e *Exif
	args := []string{"-c", gpsPrecisionFmt, p}

	out, err := exec.Command("exiftool", args...).Output()
	if err != nil {
		return e, err
	}

	e = new(Exif)

	for _, l := range strings.Split(string(out), "\n") {
		parts := strings.SplitN(l, ": ", 2)
		if len(parts) != 2 {
			continue
		}

		field, value := parts[0], parts[1]

		if ok, _ := regexp.MatchString("Date/Time Original", field); ok {
			e.DateTimeOriginal = value
		} else if ok, _ := regexp.MatchString("GPS Latitude +$", field); ok {
			v, err := valueForCoordinateString(value)
			if err == nil {
				e.GPS.Latitude = v
			}
		} else if ok, _ := regexp.MatchString("GPS Longitude +$", field); ok {
			v, err := valueForCoordinateString(value)
			if err == nil {
				e.GPS.Longitude = v
			}
		}
	}

	return e, nil
}

func valueForCoordinateString(coord string) (float64, error) {
	parts := strings.Split(coord, " ")
	numStr, dir := parts[0], parts[1]

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0.0, err
	}

	sign := 1.0
	if dir == "W" || dir == "S" {
		sign = -1.0
	}

	return (num * sign), nil
}
