package response

import "transport-app/app/domain"

type SearchNodesResponse struct {
	Name        string `json:"name"`
	ReferenceID string `json:"referenceID"`
}

func MapSearchNodesResponse(nodes []domain.NodeInfo) []SearchNodesResponse {
	var responses []SearchNodesResponse
	for _, node := range nodes {
		responses = append(responses, SearchNodesResponse{
			Name:        node.Name,
			ReferenceID: string(node.ReferenceID),
		})
	}
	return responses
}
