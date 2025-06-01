package domain

import "github.com/cockroachdb/errors"

// Pagination defines Relay-style pagination parameters.
//
// Fields:
// - First:  number of items to return starting after the 'After' cursor (forward pagination).
// - Last:   number of items to return ending before the 'Before' cursor (backward pagination).
// - After:  base64-encoded cursor indicating where to start the page (used with 'First').
// - Before: base64-encoded cursor indicating where to end the page (used with 'Last').
//
// ✅ Visual example:
//
// Item    Generated Cursor
// -----   -----------------
//   1     "cursor:1"
//   2     "cursor:2"
//   ...   ...
//  10     "cursor:10"
//
// To get items 11 to 20, use:
//
// query {
//   deliveryUnits(first: 10, after: "cursor:10")
// }
//
// The backend interprets this as:
// → “Give me the next 10 results that come after the item with cursor 'cursor:10'.”
//
// ⚠️ Pagination rules (Relay compliant):
//
// - Use either `first` or `last`, but not both.
// - `first` can be used alone or with `after`, but never with `before`.
// - `last` can be used alone or with `before`, but never with `after`.
//
// Valid combinations:
//   ✓ first
//   ✓ first + after
//   ✓ last
//   ✓ last + before
//
// Invalid combinations:
//   ✗ first + last
//   ✗ first + before
//   ✗ last + after
type Pagination struct {
	First *int    // ✅ number of items to return after the cursor (forward pagination)
	After *string // ✅ base64 cursor indicating where to start (used with 'First')

	Last   *int    // ✅ number of items to return before the cursor (backward pagination)
	Before *string // ✅ base64 cursor indicating where to end (used with 'Last')
}

func (p Pagination) IsValid() error {
	if p.First != nil && p.Last != nil {
		return errors.New("invalid pagination: cannot use both 'first' and 'last'")
	}
	if p.First != nil && p.Before != nil {
		return errors.New("invalid pagination: cannot use 'first' with 'before'")
	}
	if p.Last != nil && p.After != nil {
		return errors.New("invalid pagination: cannot use 'last' with 'after'")
	}
	return nil
}

func (p Pagination) IsForward() bool {
	return p.First != nil
}

func (p Pagination) IsBackward() bool {
	return p.Last != nil
}

func (p Pagination) HasAfter() bool {
	return p.After != nil && *p.After != ""
}

func (p Pagination) HasBefore() bool {
	return p.Before != nil && *p.Before != ""
}
