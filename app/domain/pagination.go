package domain

// Pagination defines Relay-style pagination parameters.
//
// Fields:
// - First:  number of items to return starting after the 'After' cursor.
// - Last:   number of items to return ending before the 'Before' cursor.
// - After:  base64-encoded cursor indicating where to start the page.
// - Before: base64-encoded cursor indicating where to end the page.
//
// Example of a base64 cursor: "MjAyNS0wOC0xMlQwMDowMDowMFo=|123"
type Pagination struct {
	First  *int
	Last   *int
	After  *string
	Before *string
}
