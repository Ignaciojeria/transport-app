// En transport-app/app/adapter/out/agents/tools/tools.go
package tools

import "google.golang.org/genai"

// GetAllMenuTools retorna todas las herramientas disponibles para Gemini.
func GetAllMenuTools() []*genai.Tool {
	menuTools := []*genai.FunctionDeclaration{
		CreateMenuTool(),
		// TODO: Implementar eventos de dominio para las demás tools
		//	UpdateBusinessInfoTool(),
		//	RemoveMenuItemTool(),
		//	UpdateMenuItemTool(),
		//	ReplaceFullMenuTool(),
	}
	return []*genai.Tool{{FunctionDeclarations: menuTools}}
}
