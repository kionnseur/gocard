package ui

import "math"

// ColorCycle cycles through rainbow colors.
func ColorCycle(step int) (uint8, uint8, uint8) {
	phase := step % 1536
	section := phase / 256
	offset := phase % 256

	var r, g, b uint8
	switch section {
	case 0: // Red to Yellow
		r, g, b = 255, uint8(offset), 0
	case 1: // Yellow to Green
		r, g, b = uint8(255-offset), 255, 0
	case 2: // Green to Cyan
		r, g, b = 0, 255, uint8(offset)
	case 3: // Cyan to Blue
		r, g, b = 0, uint8(255-offset), 255
	case 4: // Blue to Magenta
		r, g, b = uint8(offset), 0, 255
	case 5: // Magenta to Red
		r, g, b = 255, 0, uint8(255-offset)
	}
	return r, g, b
}

// Creates a soft breathing color effect using sine waves.
func ColorBreathSin(step int) (uint8, uint8, uint8) {
	f := float64(step) * 0.05
	r := uint8((math.Sin(f) + 1) * 127)
	g := uint8((math.Sin(f+2) + 1) * 127)
	b := uint8((math.Sin(f+4) + 1) * 127)
	return r, g, b
}

// Add a step to RGB values for a simple color shift.
func ColorBreath(r, g, b uint8, step int) (uint8, uint8, uint8) {
	r = uint8((int(r) + step) % 256)
	g = uint8((int(g) + step) % 256)
	b = uint8((int(b) + step) % 256)
	return r, g, b
}
