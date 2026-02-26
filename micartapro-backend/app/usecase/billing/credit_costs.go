package billing

// Costos de créditos por tipo de operación.
// Permite ajustar el consumo según el procesamiento real (texto, imagen, TTS, etc.).
const (
	// CreditsPerAgentUsage créditos por uso del agente de menú (LLM).
	CreditsPerAgentUsage = 1

	// CreditsPerImageGeneration créditos por generación de imagen (Imagen 4).
	CreditsPerImageGeneration = 1

	// CreditsPerImageEdition créditos por edición de imagen (image-to-image).
	CreditsPerImageEdition = 1

	// CreditsPerSpeechTTS créditos por generación de audio TTS.
	CreditsPerSpeechTTS = 1

	// CreditsPerSubtitleTiming créditos por estimación de timing de subtítulos (LLM).
	CreditsPerSubtitleTiming = 1

	// CreditsPerSceneImage créditos por imagen de escena para reels (Imagen 4).
	CreditsPerSceneImage = 1

	// CreditsPerSceneVideo créditos por video de escena para reels (Veo 3.1).
	CreditsPerSceneVideo = 5

	// CreditsPerSpeechSubtitle créditos por transcripción de audio existente (Speech-to-Text).
	CreditsPerSpeechSubtitle = 1

	// CreditsPerSpeechVoiceClone créditos por TTS con clonación de voz (ElevenLabs IVC).
	CreditsPerSpeechVoiceClone = 3
)
