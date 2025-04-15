package gemini

import "github.com/google/generative-ai-go/genai"

var AddressNormalizationSchema *genai.Schema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"addressLine1": {Type: genai.TypeString},
		"addressLine2": {Type: genai.TypeString},
		"district":     {Type: genai.TypeString},
		"state":        {Type: genai.TypeString},
		"province":     {Type: genai.TypeString},
		"latitude":     {Type: genai.TypeNumber},
		"longitude":    {Type: genai.TypeNumber},
	},
	Required: []string{
		"providerAddress",
		"addressLine1",
		"district",
		"state",
		"province",
		"latitude",
		"longitude"},
}
