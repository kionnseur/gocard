package ui

import "math"

// Fait tourner la couleur sur le cercle chromatique (arc-en-ciel)
func ColorCycle(step int) (uint8, uint8, uint8) {
	// 0 <= step < 1536 pour un cycle complet (256*6)
	phase := step % 1536
	section := phase / 256
	offset := phase % 256

	var r, g, b uint8
	switch section {
	case 0: // Rouge -> Jaune
		r, g, b = 255, uint8(offset), 0
	case 1: // Jaune -> Vert
		r, g, b = uint8(255-offset), 255, 0
	case 2: // Vert -> Cyan
		r, g, b = 0, 255, uint8(offset)
	case 3: // Cyan -> Bleu
		r, g, b = 0, uint8(255-offset), 255
	case 4: // Bleu -> Magenta
		r, g, b = uint8(offset), 0, 255
	case 5: // Magenta -> Rouge
		r, g, b = 255, 0, uint8(255-offset)
	}
	return r, g, b
}

// Variante sinus pour un effet de "respiration" colorée douce
func ColorBreathSin(step int) (uint8, uint8, uint8) {
	f := float64(step) * 0.05
	r := uint8((math.Sin(f) + 1) * 127)
	g := uint8((math.Sin(f+2) + 1) * 127)
	b := uint8((math.Sin(f+4) + 1) * 127)
	return r, g, b
}

func ColorBreath(r, g, b uint8, step int) (uint8, uint8, uint8) {

	// On fait varier les couleurs en fonction du step
	r = uint8((int(r) + step) % 256)
	g = uint8((int(g) + step) % 256)
	b = uint8((int(b) + step) % 256)

	// On retourne les nouvelles valeurs de couleur
	return r, g, b
}
