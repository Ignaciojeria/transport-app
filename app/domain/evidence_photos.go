package domain

import (
	"context"
	"sort"
	"time"
)

type EvidencePhoto struct {
	URL     string
	Type    string
	TakenAt time.Time
}

type EvidencePhotos []EvidencePhoto

func (e EvidencePhotos) DocID(ctx context.Context) DocumentID {
	// Extract and sort URLs
	urls := make([]string, len(e))
	for i, photo := range e {
		urls[i] = photo.URL
	}
	sort.Strings(urls)

	// Generate hash using sorted URLs
	return HashByTenant(ctx, urls...)
}
