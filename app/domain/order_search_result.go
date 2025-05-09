package domain

type OrderSearchResult struct {
	Plan        Plan
	HasNextPage bool
	EndCursor   *string // opcional si quieres soportar cursor-based
}
