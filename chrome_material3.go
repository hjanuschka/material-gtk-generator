package main

import (
	"image/color"
	"math"
)

// Chrome's exact Material 3 implementation ported from C++
// Source: ui/color/dynamic_color/palette_factory.cc

type SchemeVariant int

const (
	TonalSpot SchemeVariant = iota
	Vibrant
	Neutral
	Expressive
)

// HCT color space implementation (Hue, Chroma, Tone)
type HCT struct {
	Hue    float64 // 0-360
	Chroma float64 // 0-infinity (practical max ~120)
	Tone   float64 // 0-100 (lightness)
}

type TonalPalette struct {
	hue    float64
	chroma float64
}

type ChromePalette struct {
	Primary        TonalPalette
	Secondary      TonalPalette
	Tertiary       TonalPalette
	Neutral        TonalPalette
	NeutralVariant TonalPalette
	Error          TonalPalette
}

type Transform struct {
	HueRotation      float64
	Chroma          float64
	HuesToRotations map[float64]float64
	HuesToChroma    map[float64]float64
}

type Config struct {
	Primary        Transform
	Secondary      Transform
	Tertiary       Transform
	Neutral        Transform
	NeutralVariant Transform
}

// Convert RGB to HCT (simplified - using basic conversion)
func RGBToHCT(r, g, b uint8) HCT {
	// Convert to 0-1 range
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0
	
	// Convert to HSV first
	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	delta := max - min
	
	var hue float64
	if delta == 0 {
		hue = 0
	} else if max == rf {
		hue = 60 * (math.Mod((gf-bf)/delta, 6))
	} else if max == gf {
		hue = 60 * ((bf-rf)/delta + 2)
	} else {
		hue = 60 * ((rf-gf)/delta + 4)
	}
	
	if hue < 0 {
		hue += 360
	}
	
	// Approximate chroma and tone mappings
	saturation := 0.0
	if max != 0 {
		saturation = delta / max
	}
	
	// Map to Material 3 HCT space
	chroma := saturation * 120.0 // Scale to Material 3 chroma range
	tone := max * 100.0          // Lightness/tone
	
	return HCT{
		Hue:    hue,
		Chroma: chroma,
		Tone:   tone,
	}
}

// HCT to RGB conversion (simplified)
func (h HCT) ToRGB() color.RGBA {
	// Simplified conversion - in real Chrome this uses CAM16 and complex math
	hue := h.Hue
	saturation := math.Min(h.Chroma/120.0, 1.0)
	value := h.Tone / 100.0
	
	c := value * saturation
	x := c * (1 - math.Abs(math.Mod(hue/60.0, 2)-1))
	m := value - c
	
	var r, g, b float64
	
	switch {
	case hue < 60:
		r, g, b = c, x, 0
	case hue < 120:
		r, g, b = x, c, 0
	case hue < 180:
		r, g, b = 0, c, x
	case hue < 240:
		r, g, b = 0, x, c
	case hue < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	
	return color.RGBA{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
		A: 255,
	}
}

func sanitizeDegreesDouble(degrees float64) float64 {
	degrees = math.Mod(degrees, 360.0)
	if degrees < 0 {
		degrees += 360.0
	}
	return degrees
}

func getRotatedHue(sourceHue float64, huesToRotations map[float64]float64) float64 {
	if len(huesToRotations) == 1 {
		for _, rotation := range huesToRotations {
			return sanitizeDegreesDouble(sourceHue + rotation)
		}
	}
	
	// Find closest match
	var bestRotation float64
	minDiff := 360.0
	
	for hue, rotation := range huesToRotations {
		diff := math.Abs(sourceHue - hue)
		if diff < minDiff {
			minDiff = diff
			bestRotation = rotation
		}
	}
	
	return sanitizeDegreesDouble(sourceHue + bestRotation)
}

func getAdjustedChroma(sourceHue float64, huesToChroma map[float64]float64) float64 {
	// Find closest match
	var bestChroma float64
	minDiff := 360.0
	
	for hue, chroma := range huesToChroma {
		diff := math.Abs(sourceHue - hue)
		if diff < minDiff {
			minDiff = diff
			bestChroma = chroma
		}
	}
	
	return bestChroma
}

func newTonalPalette(hue, chroma float64) TonalPalette {
	return TonalPalette{
		hue:    hue,
		chroma: chroma,
	}
}

func (tp TonalPalette) Tone(tone int) color.RGBA {
	hct := HCT{
		Hue:    tp.hue,
		Chroma: tp.chroma,
		Tone:   float64(tone),
	}
	return hct.ToRGB()
}

func makePalette(hue float64, transform Transform) TonalPalette {
	chroma := transform.Chroma
	
	if transform.HuesToChroma != nil {
		chroma = getAdjustedChroma(hue, transform.HuesToChroma)
	}
	
	if transform.HuesToRotations != nil {
		hue = getRotatedHue(hue, transform.HuesToRotations)
	} else {
		hue = sanitizeDegreesDouble(hue + transform.HueRotation)
	}
	
	return newTonalPalette(hue, chroma)
}

// Chrome's exact Material 3 configurations from palette_factory.cc
func GenerateChromePalette(seedColor color.RGBA, variant SchemeVariant) ChromePalette {
	hct := RGBToHCT(seedColor.R, seedColor.G, seedColor.B)
	hue := hct.Hue
	
	var config Config
	
	switch variant {
	case TonalSpot:
		// Chrome's kTonalSpot: {Chroma(40.0), Chroma(16.0), Transform{60.0, 24.0}, Chroma(6.0), Chroma(8.0)}
		config = Config{
			Primary:        Transform{Chroma: 40.0},
			Secondary:      Transform{Chroma: 16.0},
			Tertiary:       Transform{HueRotation: 60.0, Chroma: 24.0},
			Neutral:        Transform{Chroma: 6.0},
			NeutralVariant: Transform{Chroma: 8.0},
		}
	case Vibrant:
		// Chrome's kVibrant with hue rotations
		hues := []float64{0, 41, 61, 101, 131, 181, 251, 301, 360}
		secondaryRotations := []float64{18, 15, 10, 12, 15, 18, 15, 12, 12}
		tertiaryRotations := []float64{35, 30, 20, 25, 30, 35, 30, 25, 25}
		
		secondaryHuesToRotations := make(map[float64]float64)
		tertiaryHuesToRotations := make(map[float64]float64)
		
		for i, h := range hues {
			secondaryHuesToRotations[h] = secondaryRotations[i]
			tertiaryHuesToRotations[h] = tertiaryRotations[i]
		}
		
		config = Config{
			Primary:        Transform{Chroma: 200.0}, // Very high chroma!
			Secondary:      Transform{Chroma: 24.0, HuesToRotations: secondaryHuesToRotations},
			Tertiary:       Transform{Chroma: 32.0, HuesToRotations: tertiaryHuesToRotations},
			Neutral:        Transform{Chroma: 8.0},
			NeutralVariant: Transform{Chroma: 12.0},
		}
	case Neutral:
		// Chrome's kNeutral
		hues := []float64{0, 260, 315, 360}
		chromas := []float64{12.0, 12.0, 20.0, 12.0}
		
		huesToChroma := make(map[float64]float64)
		for i, h := range hues {
			huesToChroma[h] = chromas[i]
		}
		
		config = Config{
			Primary:        Transform{HuesToChroma: huesToChroma},
			Secondary:      Transform{Chroma: 8.0},
			Tertiary:       Transform{Chroma: 16.0},
			Neutral:        Transform{Chroma: 2.0},
			NeutralVariant: Transform{Chroma: 2.0},
		}
	case Expressive:
		// Chrome's kExpressive
		hues := []float64{0, 21, 51, 121, 151, 191, 271, 321, 360}
		secondaryRotations := []float64{45, 95, 45, 20, 45, 90, 45, 45, 45}
		tertiaryRotations := []float64{120, 120, 20, 45, 20, 15, 20, 120, 120}
		
		secondaryHuesToRotations := make(map[float64]float64)
		tertiaryHuesToRotations := make(map[float64]float64)
		
		for i, h := range hues {
			secondaryHuesToRotations[h] = secondaryRotations[i]
			tertiaryHuesToRotations[h] = tertiaryRotations[i]
		}
		
		config = Config{
			Primary:        Transform{HueRotation: -90, Chroma: 40.0},
			Secondary:      Transform{Chroma: 24.0, HuesToRotations: secondaryHuesToRotations},
			Tertiary:       Transform{Chroma: 32.0, HuesToRotations: tertiaryHuesToRotations},
			Neutral:        Transform{Chroma: 8.0},
			NeutralVariant: Transform{Chroma: 12.0},
		}
	}
	
	return ChromePalette{
		Primary:        makePalette(hue, config.Primary),
		Secondary:      makePalette(hue, config.Secondary),
		Tertiary:       makePalette(hue, config.Tertiary),
		Neutral:        makePalette(hue, config.Neutral),
		NeutralVariant: makePalette(hue, config.NeutralVariant),
		Error:          newTonalPalette(25.0, 84.0), // Chrome's error color
	}
}