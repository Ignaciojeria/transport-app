package agents

import "encoding/json"

type AgentResponse struct {
	TextResponse string
	CommandName  string
	CommandArgs  json.RawMessage // JSON crudo de los argumentos
}
