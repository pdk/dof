package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

type millimeter float64

type size struct {
	x, y millimeter
}

var format = map[string]size{
	"35mm":          {36.0, 24.0},
	"aps-c-canon":   {22.5, 15.0},
	"aps-c-fuji":    {23.6, 15.6},
	"apx-c-generic": {24.0, 16.0},
	"m43":           {18.0, 13.5},
}

func (s size) diag() millimeter {
	return millimeter(math.Sqrt(float64(s.x*s.x + s.y*s.y)))
}

func (s size) circleOfConfusion() millimeter {

	// circleOfConfusion = 0.025mm for 35mm film
	basis := 0.025 / format["35mm"].diag()

	return s.diag() * basis
}

func (s size) depthOfField(focalLength millimeter, aperture float64, focalDistance millimeter) float64 {

	// depthOfField = 2 * focalDistance^2 * fStop * circleOfConfusion / focalLength^2
	// https://en.wikipedia.org/wiki/Depth_of_field

	return 2.0 *
		float64(focalDistance*focalDistance) *
		aperture *
		float64(s.circleOfConfusion()) /
		float64(focalLength*focalLength)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: dof fmt\n")
	for k, v := range format {
		fmt.Fprintf(os.Stderr, "%-20s %3.1fmm x %3.1fmm\n", k, v.x, v.y)
	}
}

func usage2() {
	fmt.Fprintf(os.Stderr, "usage: dof fmt focalLength aperture focusDistance\n")
}

func parse(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatalf("cannot parse %s: %v", s, err)
	}

	return v
}

func run(args []string, stdout io.Writer) error {

	if len(args) < 2 {
		usage()
		return nil
	}

	s, ok := format[args[1]]
	if !ok {
		usage()
		return nil
	}

	// fmt.Printf("diagonal of %s (%3.1fmm x %3.1fmm) is %.3fmm\n", args[1], s.x, s.y, s.diag())

	fmt.Printf("format               %s\n", args[1])
	fmt.Printf("image size           %3.1fmm x %3.1fmm\n", s.x, s.y)
	fmt.Printf("diagonal             %.3fmm\n", s.diag())
	fmt.Printf("circle of confusion  %.6fmm\n", s.circleOfConfusion())

	if len(args) == 2 {
		return nil
	}

	if len(args) < 5 {
		usage2()
		return nil
	}

	focalLength := millimeter(parse(args[2]))
	aperture := parse(args[3])
	focusDistance := millimeter(parse(args[4]) * 1000)

	fmt.Printf("focal length         %.1fmm\n", focalLength)
	fmt.Printf("aperture             f/%.1f\n", aperture)
	fmt.Printf("focus distance       %.1fmm\n", focusDistance)

	fmt.Printf("depth of field       %.1fmm\n", s.depthOfField(focalLength, aperture, focusDistance))

	return nil
}

// "standard" apertures (https://en.wikipedia.org/wiki/F-number)
// 1.0, 1.1, 1.2,
// 1.4, 1.6, 1.8,
// 2.0, 2.2, 2.5,
// 2.8, 3.2, 3.6,
// 4.0, 4.5, 5.0,
// 5.6, 6.4, 7.1,
// 8.0, 9.0, 10,
// 11, 13, 14,
// 16, 18, 20,
// 22, 25, 29,
// 32, 36, 40,
// 45, 51,

// go run dof.go 35mm 35 4.0 5
// go run dof.go aps-c-fuji 35 4.0 5
// go run dof.go aps-c-fuji 35 5.6 5

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
