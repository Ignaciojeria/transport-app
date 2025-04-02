package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Contact ReferenceID", func() {
	org := Organization{ID: 1, Country: countries.CL}

	It("should use primaryEmail when available", func() {
		contact := Contact{
			NationalID:   "12345678-9",
			PrimaryEmail: "test@example.com",
			PrimaryPhone: "+56912345678",
			Organization: org,
		}

		refID := contact.DocID()
		Expect(refID).ToNot(BeEmpty())
		Expect(refID).To(Equal(Hash(org, "test@example.com")))
	})

	It("should fallback to primaryEmail if nationalID is missing", func() {
		contact := Contact{
			PrimaryEmail: "test@example.com",
			PrimaryPhone: "+56912345678",
			Organization: org,
		}

		refID := contact.DocID()
		Expect(refID).To(Equal(Hash(org, "test@example.com")))
	})

	It("should use nationalID if primaryEmail and primaryPhone are missing", func() {
		contact := Contact{
			NationalID:   "12345678-9",
			Organization: org,
		}

		refID := contact.DocID()
		Expect(refID).To(Equal(Hash(org, "12345678-9")))
	})

	It("should fallback to primaryPhone if primaryEmail and nationalID are missing", func() {
		contact := Contact{
			PrimaryPhone: "+56912345678",
			Organization: org,
		}

		refID := contact.DocID()
		Expect(refID).To(Equal(Hash(org, "+56912345678")))
	})

	It("should generate a UUID if no identifiers are present", func() {
		contact1 := Contact{Organization: org}
		contact2 := Contact{Organization: org}

		ref1 := contact1.DocID()
		ref2 := contact2.DocID()

		Expect(ref1).ToNot(BeEmpty())
		Expect(ref2).ToNot(BeEmpty())
		Expect(ref1).ToNot(Equal(ref2)) // UUIDs aleatorios → distintos
	})
})

var _ = Describe("Contact UpdateIfChanged", func() {
	var original Contact
	org := Organization{ID: 1, Country: countries.CL}

	BeforeEach(func() {
		original = Contact{
			FullName:     "Juan Pérez",
			PrimaryEmail: "juan@correo.com",
			PrimaryPhone: "+56900000000",
			NationalID:   "12345678-9",
			Documents:    []Document{},
			Organization: org,
		}
	})

	It("should not update if nothing changes", func() {
		updated, changed := original.UpdateIfChanged(original)
		Expect(changed).To(BeFalse())
		Expect(updated).To(Equal(original))
	})

	It("should update FullName", func() {
		updated, changed := original.UpdateIfChanged(Contact{FullName: "Juan Pablo"})
		Expect(changed).To(BeTrue())
		Expect(updated.FullName).To(Equal("Juan Pablo"))
	})

	It("should update PrimaryEmail", func() {
		updated, changed := original.UpdateIfChanged(Contact{PrimaryEmail: "nuevo@correo.com"})
		Expect(changed).To(BeTrue())
		Expect(updated.PrimaryEmail).To(Equal("nuevo@correo.com"))
	})

	It("should update PrimaryPhone", func() {
		updated, changed := original.UpdateIfChanged(Contact{PrimaryPhone: "+56912345678"})
		Expect(changed).To(BeTrue())
		Expect(updated.PrimaryPhone).To(Equal("+56912345678"))
	})

	It("should update NationalID", func() {
		updated, changed := original.UpdateIfChanged(Contact{NationalID: "98765432-1"})
		Expect(changed).To(BeTrue())
		Expect(updated.NationalID).To(Equal("98765432-1"))
	})

	It("should update Documents", func() {
		newDocs := []Document{{Type: "RUT", Value: "1234"}}
		updated, changed := original.UpdateIfChanged(Contact{Documents: newDocs})
		Expect(changed).To(BeTrue())
		Expect(updated.Documents).To(Equal(newDocs))
	})

	It("should update multiple fields", func() {
		updated, changed := original.UpdateIfChanged(Contact{
			FullName:     "Nombre nuevo",
			PrimaryPhone: "+56911111111",
			NationalID:   "11111111-1",
		})
		Expect(changed).To(BeTrue())
		Expect(updated.FullName).To(Equal("Nombre nuevo"))
		Expect(updated.PrimaryPhone).To(Equal("+56911111111"))
		Expect(updated.NationalID).To(Equal("11111111-1"))
	})
})
