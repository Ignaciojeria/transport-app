package subtitles

const (
	ThemeTrailerDefault = "TRAILER_DEFAULT"
)

// ThemeTrailerDefaultPreset fuentes y colores para estilo trailer.
// Primary: Bebas Neue (impacto, CENTER). Secondary: Montserrat (soporte, BOTTOM).
// Accent #7C3AED: tecnológico, funciona sobre fondo oscuro.
var ThemeTrailerDefaultPreset = &SubtitleTheme{
	PrimaryFont:    "Bebas Neue",
	SecondaryFont:  "Montserrat",
	PrimaryColor:   "#FFFFFF",
	AccentColor:    "#7C3AED",
	SecondaryColor: "#D1D5DB",
}

// GetThemePreset retorna el tema por nombre o nil.
func GetThemePreset(name string) *SubtitleTheme {
	switch name {
	case ThemeTrailerDefault:
		return ThemeTrailerDefaultPreset
	default:
		return nil
	}
}
