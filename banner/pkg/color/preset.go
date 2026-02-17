package color

// Preset gradients for common use cases.
var (
	// GradientSunset - Warm sunset colors (red to orange to pink).
	GradientSunset = Gradient{
		Colors: []RGB{{255, 107, 107}, {254, 202, 87}, {255, 159, 243}},
	}

	// GradientOcean - Cool ocean blues.
	GradientOcean = Gradient{
		Colors: []RGB{{0, 82, 212}, {67, 100, 247}, {111, 177, 252}},
	}

	// GradientNeon - Bright neon colors (magenta to cyan).
	GradientNeon = Gradient{
		Colors: []RGB{{255, 0, 255}, {0, 255, 255}},
	}

	// GradientCyberpunk - Cyberpunk purples and blues.
	GradientCyberpunk = Gradient{
		Colors: []RGB{{247, 37, 133}, {114, 9, 183}, {58, 12, 163}},
	}

	// GradientMiami - Miami vice colors (pink to cyan).
	GradientMiami = Gradient{
		Colors: []RGB{{247, 37, 133}, {76, 201, 240}},
	}

	// GradientFire - Hot fire colors (red to orange to yellow).
	GradientFire = Gradient{
		Colors: []RGB{{241, 39, 17}, {245, 175, 25}},
	}

	// GradientForest - Natural forest greens.
	GradientForest = Gradient{
		Colors: []RGB{{19, 78, 94}, {113, 178, 128}},
	}

	// GradientGalaxy - Deep purple galaxy colors.
	GradientGalaxy = Gradient{
		Colors: []RGB{{127, 0, 255}, {225, 0, 255}},
	}

	// GradientRetro - Retro pink and blue.
	GradientRetro = Gradient{
		Colors: []RGB{{252, 70, 107}, {63, 94, 251}},
	}

	// GradientAurora - Northern lights colors (cyan to purple to pink).
	GradientAurora = Gradient{
		Colors: []RGB{{0, 198, 255}, {0, 114, 255}, {114, 9, 183}, {247, 37, 133}},
	}

	// GradientMint - Fresh mint green.
	GradientMint = Gradient{
		Colors: []RGB{{0, 176, 155}, {150, 201, 61}},
	}

	// GradientPeach - Soft peach colors.
	GradientPeach = Gradient{
		Colors: []RGB{{255, 154, 158}, {250, 208, 196}},
	}

	// GradientLavender - Soft lavender colors.
	GradientLavender = Gradient{
		Colors: []RGB{{150, 131, 236}, {246, 191, 255}},
	}

	// GradientGold - Rich gold colors.
	GradientGold = Gradient{
		Colors: []RGB{{255, 215, 0}, {255, 165, 0}, {184, 134, 11}},
	}

	// GradientIce - Cool ice blues.
	GradientIce = Gradient{
		Colors: []RGB{{230, 240, 255}, {135, 206, 250}, {70, 130, 180}},
	}

	// GradientBlood - Dark red blood colors.
	GradientBlood = Gradient{
		Colors: []RGB{{139, 0, 0}, {220, 20, 60}, {178, 34, 34}},
	}

	// GradientMatrix - Matrix green colors.
	GradientMatrix = Gradient{
		Colors: []RGB{{0, 50, 0}, {0, 255, 65}, {0, 100, 0}},
	}

	// GradientVaporwave - Vaporwave pink, cyan, and purple.
	GradientVaporwave = Gradient{
		Colors: []RGB{{255, 113, 206}, {1, 205, 254}, {185, 103, 255}},
	}

	// GradientRainbow - Full rainbow spectrum.
	GradientRainbow = Gradient{
		Colors: []RGB{
			{255, 0, 0}, {255, 127, 0}, {255, 255, 0},
			{0, 255, 0}, {0, 0, 255}, {75, 0, 130}, {148, 0, 211},
		},
	}

	// GradientTerminal - Classic terminal green.
	GradientTerminal = Gradient{
		Colors: []RGB{{0, 255, 0}, {0, 200, 0}},
	}

	// GradientRose - Romantic rose colors.
	GradientRose = Gradient{
		Colors: []RGB{{255, 0, 128}, {255, 102, 178}, {255, 179, 217}},
	}

	// GradientSky - Sky blue gradient.
	GradientSky = Gradient{
		Colors: []RGB{{135, 206, 235}, {70, 130, 180}, {25, 25, 112}},
	}
)

// gradientMap provides name-based lookup for gradients.
var gradientMap = map[string]Gradient{
	"sunset":    GradientSunset,
	"ocean":     GradientOcean,
	"neon":      GradientNeon,
	"cyberpunk": GradientCyberpunk,
	"miami":     GradientMiami,
	"fire":      GradientFire,
	"forest":    GradientForest,
	"galaxy":    GradientGalaxy,
	"retro":     GradientRetro,
	"aurora":    GradientAurora,
	"mint":      GradientMint,
	"peach":     GradientPeach,
	"lavender":  GradientLavender,
	"gold":      GradientGold,
	"ice":       GradientIce,
	"blood":     GradientBlood,
	"matrix":    GradientMatrix,
	"vaporwave": GradientVaporwave,
	"rainbow":   GradientRainbow,
	"terminal":  GradientTerminal,
	"rose":      GradientRose,
	"sky":       GradientSky,
}

// GetGradient returns a gradient by name. Returns GradientAurora if not found.
func GetGradient(name string) Gradient {
	if g, ok := gradientMap[name]; ok {
		return g
	}
	return GradientAurora
}

// ListGradients returns a list of all available gradient names.
func ListGradients() []string {
	names := make([]string, 0, len(gradientMap))
	for name := range gradientMap {
		names = append(names, name)
	}
	return names
}

// HasGradient checks if a gradient with the given name exists.
func HasGradient(name string) bool {
	_, ok := gradientMap[name]
	return ok
}
