package response

import "transport-app/app/domain"

type SearchNodesResponse struct {
	Name        string `json:"name"`
	ReferenceID string `json:"referenceId"`
}

func MapSearchNodesResponse(nodes []domain.NodeInfo) []SearchNodesResponse {
	var responses []SearchNodesResponse
	for _, node := range nodes {
		responses = append(responses, SearchNodesResponse{
			Name:        getStringOrDefault(node.Name, "Unknown"),
			ReferenceID: string(node.ReferenceID),
		})
	}
	return responses
}

func getStringOrDefault(value *string, defaultValue string) string {
	if value != nil {
		return *value
	}
	return defaultValue
}
