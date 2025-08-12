package mapper

func buildIndex(syn map[string][]string) map[string]string {
	idx := make(map[string]string, 256)
	for canonical, alts := range syn {
		idx[norm(canonical)] = canonical
		for _, alt := range alts {
			idx[norm(alt)] = canonical
		}
	}
	return idx
}
