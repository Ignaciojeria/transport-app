package domain

type Item struct {
	Sku    string
	Skills []struct {
		Type        string `json:"type"`
		Value       string `json:"value"`
		Description string `json:"description"`
	} `json:"skills"`
	Quantity    Quantity
	Insurance   Insurance
	Description string
	Dimensions  Dimensions
	Weight      Weight
}

type ItemReference struct {
	Sku      string
	Quantity Quantity
}

// Función auxiliar para comparar arreglos de referencias de ítems
func compareItemReferences(oldRefs, newRefs []ItemReference) bool {
	if len(oldRefs) != len(newRefs) {
		return false
	}
	for i := range oldRefs {
		if oldRefs[i] != newRefs[i] {
			return false
		}
	}
	return true
}
