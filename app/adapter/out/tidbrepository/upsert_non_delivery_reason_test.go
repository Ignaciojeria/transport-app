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

var _ = Describe("UpsertNonDeliveryReason", func() {
	createOrgContext := func(org domain.Organization) context.Context {
		ctx := context.Background()
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, strconv.FormatInt(org.ID, 10))
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, org.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		return baggage.ContextWithBaggage(ctx, bag)
	}

	It("should insert non delivery reason if not exists", func() {
		ctx := createOrgContext(organization1)
		reason := domain.NonDeliveryReason{
			ReferenceID: "REF-001",
			Reason:      "NO_ONE_HOME",
			Details:     "Cliente no estaba presente",
		}

		upsert := NewUpsertNonDeliveryReason(connection)
		err := upsert(ctx, reason)
		Expect(err).ToNot(HaveOccurred())

		var dbReason table.NonDeliveryReason
		err = connection.DB.WithContext(ctx).
			Table("non_delivery_reasons").
			Where("document_id = ?", reason.DocID(ctx)).
			First(&dbReason).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbReason.Reason).To(Equal("NO_ONE_HOME"))
	})

	It("should update non delivery reason if fields are different", func() {
		ctx := createOrgContext(organization1)

		original := domain.NonDeliveryReason{
			ReferenceID: "REF-002",
			Reason:      "BAD_ADDRESS",
			Details:     "Dirección no existe",
		}
		upsert := NewUpsertNonDeliveryReason(connection)
		_ = upsert(ctx, original)

		modified := domain.NonDeliveryReason{
			ReferenceID: "REF-002",
			Reason:      "WRONG_PHONE",
			Details:     "Número incorrecto",
		}
		err := upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbReason table.NonDeliveryReason
		err = connection.DB.WithContext(ctx).
			Table("non_delivery_reasons").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbReason).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbReason.Reason).To(Equal("WRONG_PHONE"))
	})

	It("should not update non delivery reason if no changes", func() {
		ctx := createOrgContext(organization1)

		reason := domain.NonDeliveryReason{
			ReferenceID: "REF-003",
			Reason:      "INCOMPLETE_ADDRESS",
		}

		upsert := NewUpsertNonDeliveryReason(connection)
		err := upsert(ctx, reason)
		Expect(err).ToNot(HaveOccurred())

		var original table.NonDeliveryReason
		err = connection.DB.WithContext(ctx).
			Table("non_delivery_reasons").
			Where("document_id = ?", reason.DocID(ctx)).
			First(&original).Error
		Expect(err).ToNot(HaveOccurred())

		// Ejecutar nuevamente sin cambios
		err = upsert(ctx, reason)
		Expect(err).ToNot(HaveOccurred())

		var dbReason table.NonDeliveryReason
		err = connection.DB.WithContext(ctx).
			Table("non_delivery_reasons").
			Where("document_id = ?", reason.DocID(ctx)).
			First(&dbReason).Error
		Expect(err).ToNot(HaveOccurred())

		Expect(dbReason.UpdatedAt).To(Equal(original.UpdatedAt))
	})

	It("should fail if database has no non_delivery_reasons table", func() {
		ctx := createOrgContext(organization1)

		reason := domain.NonDeliveryReason{
			ReferenceID: "REF-999",
			Reason:      "DATABASE_MISSING",
		}
		upsert := NewUpsertNonDeliveryReason(noTablesContainerConnection)
		err := upsert(ctx, reason)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("non_delivery_reasons"))
	})
})
