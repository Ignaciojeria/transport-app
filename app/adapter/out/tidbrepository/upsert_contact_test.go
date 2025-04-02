package tidbrepository

import (
	"context"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertContact", func() {

	It("should insert contact if not exists", func() {
		ctx := context.Background()

		contact := domain.Contact{
			FullName:     "Juan Pérez",
			PrimaryEmail: "juan@example.com",
			PrimaryPhone: "123456789",
			Organization: organization1,
		}

		upsert := NewUpsertContact(connection)
		err := upsert(ctx, contact)
		Expect(err).ToNot(HaveOccurred())

		var dbContact table.Contact
		err = connection.DB.WithContext(ctx).
			Table("contacts").
			Where("reference_id = ?", contact.DocID()).
			First(&dbContact).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact.FullName).To(Equal("Juan Pérez"))
		Expect(dbContact.Phone).To(Equal("123456789"))
	})

	It("should update contact if fields are different", func() {
		ctx := context.Background()

		original := domain.Contact{
			FullName:     "Nombre Original",
			PrimaryEmail: "update@example.com",
			PrimaryPhone: "111111111",
			Organization: organization1,
		}

		upsert := NewUpsertContact(connection)

		err := upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.Contact{
			FullName:     "Nombre Modificado",
			PrimaryEmail: "update@example.com",
			PrimaryPhone: "222222222",
			Organization: organization1,
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbContact table.Contact
		err = connection.DB.WithContext(ctx).
			Table("contacts").
			Where("reference_id = ?", modified.DocID()).
			First(&dbContact).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact.FullName).To(Equal("Nombre Modificado"))
		Expect(dbContact.Phone).To(Equal("222222222"))
	})

	It("should not update if no fields changed", func() {
		ctx := context.Background()

		contact := domain.Contact{
			FullName:     "Sin Cambios",
			PrimaryEmail: "same@example.com",
			PrimaryPhone: "333333333",
			Organization: organization1,
		}

		upsert := NewUpsertContact(connection)

		err := upsert(ctx, contact)
		Expect(err).ToNot(HaveOccurred())

		// Ejecutar de nuevo sin cambios
		err = upsert(ctx, contact)
		Expect(err).ToNot(HaveOccurred())

		var dbContact table.Contact
		err = connection.DB.WithContext(ctx).
			Table("contacts").
			Where("reference_id = ?", contact.DocID()).
			First(&dbContact).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact.FullName).To(Equal("Sin Cambios"))
		Expect(dbContact.Phone).To(Equal("333333333"))
	})

	It("should allow same contact info for different organizations", func() {
		ctx := context.Background()

		contact1 := domain.Contact{
			FullName:     "Multi Org",
			PrimaryEmail: "multi@example.com",
			PrimaryPhone: "444444444",
			Organization: organization1,
		}

		contact2 := domain.Contact{
			FullName:     "Multi Org",
			PrimaryEmail: "multi@example.com",
			PrimaryPhone: "444444444",
			Organization: organization2,
		}

		upsert := NewUpsertContact(connection)

		err := upsert(ctx, contact1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx, contact2)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(ctx).
			Table("contacts").
			Where("email = ?", contact1.PrimaryEmail).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))
	})

	It("should insert using NationalID if email and phone are missing", func() {
		ctx := context.Background()

		contact := domain.Contact{
			FullName:     "Identificado por RUN",
			NationalID:   "12345678-9",
			Organization: organization1,
		}

		upsert := NewUpsertContact(connection)
		err := upsert(ctx, contact)
		Expect(err).ToNot(HaveOccurred())

		var dbContact table.Contact
		err = connection.DB.WithContext(ctx).
			Table("contacts").
			Where("reference_id = ?", contact.DocID()).
			First(&dbContact).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact.NationalID).To(Equal("12345678-9"))
	})

	It("should generate new ReferenceID if all keys are missing", func() {
		ctx := context.Background()

		contact := domain.Contact{
			FullName:     "Sin Identificación",
			Organization: organization1,
		}

		upsert := NewUpsertContact(connection)
		err := upsert(ctx, contact)
		Expect(err).ToNot(HaveOccurred())

		var dbContact table.Contact
		err = connection.DB.WithContext(ctx).
			Table("contacts").
			Where("full_name = ?", contact.FullName).
			First(&dbContact).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact.ReferenceID).ToNot(BeEmpty())
	})

	It("should fail if database has no contacts table", func() {
		ctx := context.Background()

		contact := domain.Contact{
			FullName:     "Error Esperado",
			PrimaryEmail: "fail@example.com",
			Organization: organization1,
		}

		upsert := NewUpsertContact(noTablesContainerConnection)
		err := upsert(ctx, contact)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("contacts"))
	})
})
