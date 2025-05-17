package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

var _ = Describe("UpsertPlanHeaders", func() {
	// Helper function to create context with organization
	createOrgContext := func(org domain.Tenant) context.Context {
		ctx := context.Background()
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, org.ID.String())
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, org.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		return baggage.ContextWithBaggage(ctx, bag)
	}

	BeforeEach(func() {
		// Clean the table before each test
		err := connection.DB.Exec("DELETE FROM plan_headers").Error
		Expect(err).ToNot(HaveOccurred())
	})

	It("should insert headers if not exists", func() {
		ctx := createOrgContext(organization1)

		headers := domain.Headers{
			Commerce: "CORP",
			Consumer: "CORP",
			Channel:  "WEB",
		}

		upsert := NewUpsertPlanHeaders(connection)
		err := upsert(ctx, headers)
		Expect(err).ToNot(HaveOccurred())

		var dbHeaders table.PlanHeaders
		err = connection.DB.WithContext(ctx).
			Table("plan_headers").
			Where("document_id = ?", headers.DocID(ctx)).
			First(&dbHeaders).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbHeaders.Commerce).To(Equal("CORP"))
		Expect(dbHeaders.Consumer).To(Equal("CORP"))
		Expect(dbHeaders.Channel).To(Equal("WEB"))
	})

	It("should update headers if fields are different", func() {
		ctx := createOrgContext(organization1)

		original := domain.Headers{
			Commerce: "CORP",
			Consumer: "CORP",
			Channel:  "WEB",
		}

		upsert := NewUpsertPlanHeaders(connection)
		err := upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.Headers{
			Commerce: "CORP",
			Consumer: "RETAIL",
			Channel:  "WEB",
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbHeaders table.PlanHeaders
		err = connection.DB.WithContext(ctx).
			Table("plan_headers").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbHeaders).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbHeaders.Consumer).To(Equal("RETAIL"))
	})

	It("should not update headers if fields are the same", func() {
		ctx := createOrgContext(organization1)

		original := domain.Headers{
			Commerce: "CORP",
			Consumer: "CORP",
			Channel:  "WEB",
		}

		upsert := NewUpsertPlanHeaders(connection)
		err := upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Try to update with same values
		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		var dbHeaders table.PlanHeaders
		err = connection.DB.WithContext(ctx).
			Table("plan_headers").
			Where("document_id = ?", original.DocID(ctx)).
			First(&dbHeaders).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbHeaders.Commerce).To(Equal("CORP"))
		Expect(dbHeaders.Consumer).To(Equal("CORP"))
		Expect(dbHeaders.Channel).To(Equal("WEB"))
	})

	It("should allow same headers for different organizations", func() {
		ctx1 := createOrgContext(organization1)
		ctx2 := createOrgContext(organization2)

		headers1 := domain.Headers{
			Commerce: "CORP",
			Consumer: "CORP",
			Channel:  "WEB",
		}

		headers2 := domain.Headers{
			Commerce: "CORP",
			Consumer: "CORP",
			Channel:  "WEB",
		}

		upsert := NewUpsertPlanHeaders(connection)

		err := upsert(ctx1, headers1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, headers2)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(context.Background()).
			Table("plan_headers").
			Where("commerce = ? AND consumer = ? AND channel = ?", headers1.Commerce, headers1.Consumer, headers1.Channel).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verify they have different document IDs
		Expect(headers1.DocID(ctx1)).ToNot(Equal(headers2.DocID(ctx2)))
	})

	It("should fail if database has no plan_headers table", func() {
		ctx := createOrgContext(organization1)

		headers := domain.Headers{
			Commerce: "CORP",
			Consumer: "CORP",
			Channel:  "WEB",
		}

		upsert := NewUpsertPlanHeaders(noTablesContainerConnection)
		err := upsert(ctx, headers)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("plan_headers"))
	})
})
