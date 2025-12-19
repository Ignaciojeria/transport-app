package tools

import "google.golang.org/genai"

func ReplaceFullMenuTool() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "replaceFullMenu",
		Description: "Borra el menú digital actual por completo y lo reemplaza con la nueva estructura o ítems proporcionados por el usuario. Esta acción es irreversible.",
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				// CLAVE: Recibimos el contenido del nuevo menú, que puede ser una lista de ítems o la estructura completa.
				"newMenuContent": {
					Type:        genai.TypeString,
					Description: "El contenido completo que debe reemplazar al menú actual, usualmente una lista de ítems o un JSON.",
				},
				// Se puede añadir un booleano para reforzar la intención en el prompt.
				"confirmOverwrite": {
					Type:        genai.TypeBoolean,
					Description: "Debe ser true si el usuario usa frases como 'borra todo', 'empecemos de cero' o 'reemplázalo'.",
				},
			},
			// Es requerido tener el contenido del nuevo menú para reemplazarlo.
			Required: []string{"newMenuContent"},
		},
	}
}
