package domain

type OrderSearchResult struct {
	Orders      []Order
	HasNextPage bool
	EndCursor   *string // opcional si quieres soportar cursor-based
}
