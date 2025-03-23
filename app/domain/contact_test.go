package domain

import (
	"testing"

	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestContactReferenceID(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Contact ReferenceID Suite")
}

var _ = Describe("Contact ReferenceID", func() {
	org := Organization{ID: 1, Country: countries.CL}

	It("should use nationalID when available", func() {
		contact := Contact{
			NationalID:   "12345678-9",
			Email:        "test@example.com",
			Phone:        "+56912345678",
			Organization: org,
		}

		refID := contact.ReferenceID()
		Expect(refID).ToNot(BeEmpty())
		Expect(refID).To(Equal(Hash(org, "12345678-9")))
	})

	It("should fallback to email if nationalID is missing", func() {
		contact := Contact{
			Email:        "test@example.com",
			Phone:        "+56912345678",
			Organization: org,
		}

		refID := contact.ReferenceID()
		Expect(refID).To(Equal(Hash(org, "test@example.com")))
	})

	It("should fallback to phone if nationalID and email are missing", func() {
		contact := Contact{
			Phone:        "+56912345678",
			Organization: org,
		}

		refID := contact.ReferenceID()
		Expect(refID).To(Equal(Hash(org, "+56912345678")))
	})

	It("should generate a UUID if no identifiers are present", func() {
		contact1 := Contact{Organization: org}
		contact2 := Contact{Organization: org}

		ref1 := contact1.ReferenceID()
		ref2 := contact2.ReferenceID()

		Expect(ref1).ToNot(BeEmpty())
		Expect(ref2).ToNot(BeEmpty())
		Expect(ref1).ToNot(Equal(ref2)) // UUIDs aleatorios → distintos
	})
})
