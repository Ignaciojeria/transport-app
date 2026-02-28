package speechtotext

import (
	"strings"

	"micartapro/app/shared/configuration"
	"micartapro/app/shared/infrastructure/observability"

	"cloud.google.com/go/storage"
	ioc "github.com/Ignaciojeria/ioc"
	"google.golang.org/genai"
)

func init() {
	ioc.Register(NewTranscribeAudioProvider)
}

// NewTranscribeAudioProvider retorna el TranscribeAudio según SPEECH_TO_TEXT_PROVIDER.
// "gemini" = Gemini Audio Understanding; "chirp" (default) = Chirp 3.
func NewTranscribeAudioProvider(
	obs observability.Observability,
	gcsClient *storage.Client,
	genaiClient *genai.Client,
	conf configuration.Conf,
) (TranscribeAudio, error) {
	provider := strings.TrimSpace(strings.ToLower(conf.SPEECH_TO_TEXT_PROVIDER))
	if provider == "gemini" {
		g, err := NewSpeechToTextGemini(genaiClient, obs)
		if err != nil {
			return nil, err
		}
		return g.TranscribeAudio, nil
	}
	c, err := NewSpeechToText(obs, gcsClient, conf)
	if err != nil {
		return nil, err
	}
	return c.TranscribeAudio, nil
}
