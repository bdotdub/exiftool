package exiftool

import (
	"bytes"
	"io"
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
	LensModel        string
	LensID           string
	ShutterSpeed     string
	Aperture         string
	ISO              string
	FocalLength      string
	GPS              struct {
		Latitude  float64
		Longitude float64
	}
}

func Decode(r io.Reader) (*Exif, error) {
	args := []string{"-c", gpsPrecisionFmt, "-"}

	cmd := exec.Command("exiftool", args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer stdout.Close()

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	go func(in io.WriteCloser, r io.Reader) {
		defer in.Close()
		io.Copy(in, r)
	}(stdin, r)

	done := make(chan bool)
	out := new(bytes.Buffer)

	go func(stdout io.ReadCloser) {
		defer stdout.Close()
		io.Copy(out, stdout)
		done <- true
	}(stdout)

	<-done

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return parseOutput(out.Bytes())
}

func DecodeFileAtPath(p string) (*Exif, error) {
	args := []string{"-c", gpsPrecisionFmt, p}

	out, err := exec.Command("exiftool", args...).Output()
	if err != nil {
		return nil, err
	}

	return parseOutput(out)
}

func parseOutput(out []byte) (*Exif, error) {
	e := new(Exif)

	for _, l := range strings.Split(string(out), "\n") {
		parts := strings.SplitN(l, ": ", 2)
		if len(parts) != 2 {
			continue
		}

		field, value := parts[0], parts[1]

		if ok, _ := regexp.MatchString("Date/Time Original", field); ok {
			e.DateTimeOriginal = value
		} else if ok, _ := regexp.MatchString("Lens Model", field); ok {
			e.LensModel = value
		} else if ok, _ := regexp.MatchString("Lens ID", field); ok {
			e.LensID = value
		} else if ok, _ := regexp.MatchString("Shutter Speed", field); ok {
			e.ShutterSpeed = value
		} else if ok, _ := regexp.MatchString("Aperture", field); ok {
			e.Aperture = value
		} else if ok, _ := regexp.MatchString("ISO", field); ok {
			e.ISO = value
		} else if ok, _ := regexp.MatchString("Focal Length", field); ok {
			e.FocalLength = value
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
