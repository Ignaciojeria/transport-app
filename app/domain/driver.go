package domain

import "context"

type Driver struct {
	Headers
	Name       string
	NationalID string
	Email      string
}

func (d Driver) DocID(ctx context.Context) DocumentID {
	return HashByTenant(ctx, d.NationalID)
}

func (d Driver) UpdateIfChanged(in Driver) (Driver, bool) {
	changed := false

	if in.Name != "" && in.Name != d.Name {
		d.Name = in.Name
		changed = true
	}

	if in.NationalID != "" && in.NationalID != d.NationalID {
		d.NationalID = in.NationalID
		changed = true
	}

	if in.Email != "" && in.Email != d.Email {
		d.Email = in.Email
		changed = true
	}

	return d, changed
}

func (d Driver) IsEmpty() bool {
	return d.Name == "" && d.NationalID == "" && d.Email == ""
}
