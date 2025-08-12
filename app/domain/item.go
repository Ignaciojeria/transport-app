package domain

type Item struct {
	Sku         string
	Quantity    int
	Price       int64
	Description string
	Dimensions  Dimensions
	Weight      int64
}

type ItemReference struct {
	Sku      string
	Quantity int
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
