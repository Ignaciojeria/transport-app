package tidbrepository

import (
	"context"
	"strconv"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

var _ = Describe("UpsertContact", func() {
	// Helper function to create context with organization
	createOrgContext := func(org domain.Organization) context.Context {
		ctx := context.Background()
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, strconv.FormatInt(org.ID, 10))
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, org.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		return baggage.ContextWithBaggage(ctx, bag)
	}

	It("should insert contact if not exists", func() {
		ctx := createOrgContext(organization1)

		contact := domain.Contact{
			FullName:     "Juan Pérez",
			PrimaryEmail: "juan@example.com",
			PrimaryPhone: "123456789",
		}

		upsert := NewUpsertContact(connection)
		err := upsert(ctx, contact)
		Expect(err).ToNot(HaveOccurred())

		var dbContact table.Contact
		err = connection.DB.WithContext(ctx).
			Table("contacts").
			Where("document_id = ?", contact.DocID(ctx)).
			First(&dbContact).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact.FullName).To(Equal("Juan Pérez"))
		Expect(dbContact.Phone).To(Equal("123456789"))
	})

	It("should update contact if fields are different", func() {
		ctx := createOrgContext(organization1)

		original := domain.Contact{
			FullName:     "Nombre Original",
			PrimaryEmail: "update@example.com",
			PrimaryPhone: "111111111",
		}

		upsert := NewUpsertContact(connection)

		err := upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.Contact{
			FullName:     "Nombre Modificado",
			PrimaryEmail: "update@example.com",
			PrimaryPhone: "222222222",
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbContact table.Contact
		err = connection.DB.WithContext(ctx).
			Table("contacts").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbContact).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact.FullName).To(Equal("Nombre Modificado"))
		Expect(dbContact.Phone).To(Equal("222222222"))
	})

	It("should not update if no fields changed", func() {
		ctx := createOrgContext(organization1)

		contact := domain.Contact{
			FullName:     "Sin Cambios",
			PrimaryEmail: "same@example.com",
			PrimaryPhone: "333333333",
		}

		upsert := NewUpsertContact(connection)

		err := upsert(ctx, contact)
		Expect(err).ToNot(HaveOccurred())

		// Capturar CreatedAt original para comparar después
		var originalRecord table.Contact
		err = connection.DB.WithContext(ctx).
			Table("contacts").
			Where("document_id = ?", contact.DocID(ctx)).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())

		// Ejecutar de nuevo sin cambios
		err = upsert(ctx, contact)
		Expect(err).ToNot(HaveOccurred())

		var dbContact table.Contact
		err = connection.DB.WithContext(ctx).
			Table("contacts").
			Where("document_id = ?", contact.DocID(ctx)).
			First(&dbContact).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact.FullName).To(Equal("Sin Cambios"))
		Expect(dbContact.Phone).To(Equal("333333333"))
		Expect(dbContact.CreatedAt).To(Equal(originalRecord.CreatedAt)) // Verificar que no se modificó
	})

	It("should allow same contact info for different organizations", func() {
		ctx1 := createOrgContext(organization1)
		ctx2 := createOrgContext(organization2)

		contact1 := domain.Contact{
			FullName:     "Multi Org",
			PrimaryEmail: "multi@example.com",
			PrimaryPhone: "444444444",
		}

		contact2 := domain.Contact{
			FullName:     "Multi Org",
			PrimaryEmail: "multi@example.com",
			PrimaryPhone: "444444444",
		}

		upsert := NewUpsertContact(connection)

		err := upsert(ctx1, contact1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, contact2)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(context.Background()).
			Table("contacts").
			Where("email = ?", contact1.PrimaryEmail).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))
	})

	It("should insert using NationalID if email and phone are missing", func() {
		ctx := createOrgContext(organization1)

		contact := domain.Contact{
			FullName:   "Identificado por RUN",
			NationalID: "12345678-9",
		}

		upsert := NewUpsertContact(connection)
		err := upsert(ctx, contact)
		Expect(err).ToNot(HaveOccurred())

		var dbContact table.Contact
		err = connection.DB.WithContext(ctx).
			Table("contacts").
			Where("document_id = ?", contact.DocID(ctx)).
			First(&dbContact).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbContact.NationalID).To(Equal("12345678-9"))
	})

	It("should generate new ReferenceID if all keys are missing", func() {
		ctx := createOrgContext(organization1)

		contact := domain.Contact{
			FullName: "Sin Identificación",
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
		Expect(dbContact.DocumentID).ToNot(BeEmpty())
	})

	It("should fail if database has no contacts table", func() {
		ctx := createOrgContext(organization1)

		contact := domain.Contact{
			FullName:     "Error Esperado",
			PrimaryEmail: "fail@example.com",
		}

		upsert := NewUpsertContact(noTablesContainerConnection)
		err := upsert(ctx, contact)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("contacts"))
	})
})
