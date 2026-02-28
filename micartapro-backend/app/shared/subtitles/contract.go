package subtitles

// SchemaVersion actual del contrato de subtítulos. Incrementar al cambiar el schema.
const SchemaVersion = 2

// Estilos de subtítulo soportados.
const (
	StyleOneLine          = "ONE_LINE"
	StyleTwoLines         = "TWO_LINES"
	StyleLineByLine       = "LINE_BY_LINE"
	StyleLineByLineBig    = "LINE_BY_LINE_BIG"
	StyleCinematicDynamic = "CINEMATIC_DYNAMIC" // Placement dinámico por segmento
	StyleTrailer          = "TRAILER"           // Fragmentación agresiva, CENTER dominante, micro gaps, impacto
	StyleKaraokeWords     = "KARAOKE_WORDS"
	StyleAlexHormozi      = "ALEX_HORMOZI" // MAYÚSCULAS, palabra activa amarilla, stroke negro, scale en palabra actual
)

// PlacementStrategy: FIXED = todo igual; DYNAMIC = placement por segmento.
const (
	PlacementStrategyFixed   = "FIXED"
	PlacementStrategyDynamic = "DYNAMIC"
)

// Placements válidos.
const (
	PlacementBottom = "BOTTOM"
	PlacementCenter = "CENTER"
	PlacementTop    = "TOP"
)

// Animations para el renderer.
const (
	AnimationFadeIn  = "FADE_IN"
	AnimationScaleIn = "SCALE_IN"
	AnimationNone    = "NONE"
)

// Sizes para emphasis y presets.
const (
	SizeNormal = "NORMAL"
	SizeM      = "M" // Preset: BODY
	SizeL      = "L" // Preset: HOOK, CONTRAST
	SizeBig    = "BIG"
	SizeXL     = "XL"
	SizeXXL    = "XXL" // TRAILER: emphasis máximo
)

// DirectionPreset preset de dirección cinematográfica.
const (
	DirectionPresetCinematicDynamicV1 = "CINEMATIC_DYNAMIC_V1"
	DirectionPresetTrailerV1          = "TRAILER_V1"
	DirectionPresetDocumentary        = "DOCUMENTARY" // BOTTOM fijo, minimalista
)

// DefaultStyle estilo por defecto.
const DefaultStyle = StyleLineByLine

var ValidStyles = map[string]bool{
	StyleOneLine: true, StyleTwoLines: true, StyleLineByLine: true,
	StyleLineByLineBig: true, StyleCinematicDynamic: true, StyleTrailer: true,
	StyleKaraokeWords: true, StyleAlexHormozi: true,
}

var ValidPlacements = map[string]bool{
	PlacementBottom: true, PlacementCenter: true, PlacementTop: true,
}

// OverflowStrategy define qué hacer cuando lines excede maxLines.
const (
	OverflowAllow           = "ALLOW"                 // No hacer nada (legacy)
	OverflowRebalance       = "REBALANCE"             // Re-optimizar cortes
	OverflowShrink          = "SHRINK"                // Fusionar/truncar
	OverflowRebalanceShrink = "REBALANCE_THEN_SHRINK" // Rebalance primero, si no alcanza: shrink
)

// LayoutConstraints restricciones para placement dinámico.
type LayoutConstraints struct {
	AvoidCenterLongSentences bool `json:"avoidCenterLongSentences,omitempty"`
	PreferCenterForEmphasis  bool `json:"preferCenterForEmphasis,omitempty"`
}

// SafeArea fracciones del frame (0-1) para evitar tapar UI/caras en TikTok/IG.
type SafeArea struct {
	Top    float64 `json:"top"`
	Bottom float64 `json:"bottom"`
	Left   float64 `json:"left"`
	Right  float64 `json:"right"`
}

// DynamicRules reglas para dampening de placement (evitar cambios muy seguidos).
type DynamicRules struct {
	MinSecondsBetweenPlacementChanges float64 `json:"minSecondsBetweenPlacementChanges,omitempty"`
	MaxCenterBlocksInRow              int     `json:"maxCenterBlocksInRow,omitempty"`
}

// SubtitleLayout intención de layout (fallback cuando placementStrategy=FIXED).
type SubtitleLayout struct {
	PlacementStrategy string             `json:"placementStrategy"` // FIXED | DYNAMIC
	DefaultPlacement  string             `json:"defaultPlacement"`  // Fallback
	MaxLines          int                `json:"maxLines"`
	MaxCharsPerLine   int                `json:"maxCharsPerLine"`
	OverflowStrategy  string             `json:"overflowStrategy,omitempty"` // REBALANCE | SHRINK | REBALANCE_THEN_SHRINK
	Constraints       *LayoutConstraints `json:"constraints,omitempty"`
	SafeArea          *SafeArea          `json:"safeArea,omitempty"`
	DynamicRules      *DynamicRules      `json:"dynamicRules,omitempty"`
}

// RenderHints hints para el renderer.
type RenderHints struct {
	ExtendEndSeconds    float64 `json:"extendEndSeconds,omitempty"`
	TimingOffsetSeconds float64 `json:"timingOffsetSeconds,omitempty"` // + = subtítulos aparecen antes, - = después (corrige desfase con audio)
}

// EmphasisObject normaliza emphasis: mismo tipo para segmentos y líneas.
type EmphasisObject struct {
	Phrases []string `json:"phrases,omitempty"` // Frases a destacar (solo en segmento)
	Level   string   `json:"level,omitempty"`   // BIG | XL | XXL
	Color   string   `json:"color,omitempty"`   // Color accent para las frases (ej. #7C3AED)
}

// SubtitleTheme preset de fuentes y colores para dirección cinematográfica.
// Regla: máximo 2 fuentes, 2 colores + blanco. Primary=impacto, Secondary=soporte.
type SubtitleTheme struct {
	PrimaryFont    string `json:"primaryFont,omitempty"`    // Bebas Neue para CENTER/énfasis
	SecondaryFont  string `json:"secondaryFont,omitempty"`  // Montserrat para BOTTOM
	PrimaryColor   string `json:"primaryColor,omitempty"`   // #FFFFFF texto principal
	AccentColor    string `json:"accentColor,omitempty"`    // #7C3AED palabras clave
	SecondaryColor string `json:"secondaryColor,omitempty"` // #D1D5DB texto secundario (opcional)
}

// SubtitleWord palabra con timing para estilos karaoke/ALEX_HORMOZI.
type SubtitleWord struct {
	Text  string  `json:"text"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

// SubtitleLine representa una línea con timing propio. Hereda placement/animation del padre; solo overrides.
type SubtitleLine struct {
	Text      string          `json:"text"`
	Start     float64         `json:"start"`
	End       float64         `json:"end"`
	Placement string          `json:"placement,omitempty"`
	Animation string          `json:"animation,omitempty"`
	Emphasis  *EmphasisObject `json:"emphasis,omitempty"`
	Size      string          `json:"size,omitempty"`
	Words     []SubtitleWord  `json:"words,omitempty"` // Para ALEX_HORMOZI: palabra activa amarilla + scale
}

// SubtitleSegment representa un segmento con placement/animation por bloque.
type SubtitleSegment struct {
	Text      string          `json:"text"`
	Start     float64         `json:"start"`
	End       float64         `json:"end"`
	Placement string          `json:"placement,omitempty"`
	Animation string          `json:"animation,omitempty"`
	Emphasis  *EmphasisObject `json:"emphasis,omitempty"` // { phrases: [...], level: "BIG" }
	Size      string          `json:"size,omitempty"`
	Style     *string         `json:"style,omitempty"`
	Lines     []SubtitleLine  `json:"lines,omitempty"`
	// ImageURL imagen asociada al segmento (ej. frame de video en ese timestamp). Requiere includeImagesPerSegment.
	ImageURL string `json:"imageUrl,omitempty"`
	// VideoURL video asociado al segmento. Requiere includeImagesPerSegment con outputMediaType=video.
	VideoURL string `json:"videoUrl,omitempty"`
}

// SubtitleResponse contrato de respuesta (legacy, interno).
type SubtitleResponse struct {
	SubtitleSchemaVersion int               `json:"subtitleSchemaVersion"`
	SubtitleStyle         string            `json:"subtitleStyle"`
	DirectionPreset       string            `json:"directionPreset,omitempty"` // ej. CINEMATIC_DYNAMIC_V1
	SubtitleLayout        SubtitleLayout    `json:"subtitleLayout"`
	Theme                 *SubtitleTheme    `json:"theme,omitempty"` // Fuentes y colores (TRAILER_DEFAULT, etc.)
	RenderHints           *RenderHints      `json:"renderHints,omitempty"`
	AudioURL              string            `json:"audioUrl"`
	Subtitles             []SubtitleSegment `json:"subtitles"`
	DurationSeconds       float64           `json:"durationSeconds"`
}

// ─── TimelineResponse (formato scenes, consumido por ai-editor) ───

// SceneVoice texto del voice en la escena.
type SceneVoice struct {
	Text         string `json:"text"`
	LanguageCode string `json:"languageCode"`
}

// SceneSubtitleLine línea con timing para el renderer.
type SceneSubtitleLine struct {
	Text     string          `json:"text"`
	Start    float64         `json:"start"`
	End      float64         `json:"end"`
	Size     string          `json:"size"`
	Emphasis *EmphasisObject `json:"emphasis,omitempty"`
}

// SceneSubtitle bloque de subtítulo de la escena.
type SceneSubtitle struct {
	Placement string              `json:"placement"`
	Animation string              `json:"animation"`
	Size      string              `json:"size"`
	Lines     []SceneSubtitleLine `json:"lines"`
	Emphasis  *EmphasisObject     `json:"emphasis,omitempty"`
}

// SceneVisual imagen/video de la escena.
type SceneVisual struct {
	ImageURL  string `json:"imageUrl"`
	Animation string `json:"animation"`
	VideoURL  string `json:"videoUrl,omitempty"`
}

// SceneLayout placement y safeArea de la escena.
type SceneLayout struct {
	Placement string   `json:"placement"`
	SafeArea  SafeArea `json:"safeArea"`
}

// Scene representa una escena del timeline (1 imagen + subtítulos).
type Scene struct {
	SceneID  string        `json:"sceneId"`
	Index    int           `json:"index"`
	Start    float64       `json:"start"`
	End      float64       `json:"end"`
	Duration float64       `json:"duration"`
	Voice    SceneVoice    `json:"voice"`
	Subtitle SceneSubtitle `json:"subtitle"`
	Visual   SceneVisual   `json:"visual"`
	Layout   SceneLayout   `json:"layout"`
}

// AudioInfo info del audio para el timeline.
type AudioInfo struct {
	URL             string  `json:"url"`
	DurationSeconds float64 `json:"durationSeconds"`
}

// TimelineResponse contrato de respuesta con scenes (ai-editor).
type TimelineResponse struct {
	SubtitleSchemaVersion int            `json:"subtitleSchemaVersion"`
	TimelineID            string         `json:"timelineId"`
	Version               int            `json:"version"`
	SubtitleStyle         string         `json:"subtitleStyle"`
	DirectionPreset       string         `json:"directionPreset,omitempty"`
	SubtitleLayout        SubtitleLayout `json:"subtitleLayout"`
	Theme                 *SubtitleTheme `json:"theme,omitempty"`
	RenderHints           *RenderHints   `json:"renderHints,omitempty"`
	Audio                 AudioInfo      `json:"audio"`
	Scenes                []Scene        `json:"scenes"`
	DurationSeconds       float64        `json:"durationSeconds"`
}
