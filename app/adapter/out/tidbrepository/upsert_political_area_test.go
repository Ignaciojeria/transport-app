package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertPoliticalArea", func() {
	var (
		conn   database.ConnectionFactory
		upsert UpsertPoliticalArea
	)

	BeforeEach(func() {
		conn = connection
		upsert = NewUpsertPoliticalArea(conn)
	})

	It("should insert political area if not exists", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		politicalArea := domain.PoliticalArea{
			Code:            "cl-rm-la-florida",
			AdminAreaLevel1: "region metropolitana de santiago",
			AdminAreaLevel2: "santiago",
			AdminAreaLevel3: "la florida",
			TimeZone:        "America/Santiago",
		}

		err = upsert(ctx, politicalArea)
		Expect(err).ToNot(HaveOccurred())

		var dbPoliticalArea table.PoliticalArea
		err = conn.DB.WithContext(ctx).
			Table("political_areas").
			Where("document_id = ?", politicalArea.DocID(ctx)).
			First(&dbPoliticalArea).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPoliticalArea.Code).To(Equal("cl-rm-la-florida"))
		Expect(dbPoliticalArea.AdminAreaLevel1).To(Equal("region metropolitana de santiago"))
		Expect(dbPoliticalArea.AdminAreaLevel2).To(Equal("santiago"))
		Expect(dbPoliticalArea.AdminAreaLevel3).To(Equal("la florida"))
		Expect(dbPoliticalArea.TimeZone).To(Equal("America/Santiago"))
		Expect(dbPoliticalArea.TenantID.String()).To(Equal(tenant.ID.String()))
	})

	It("should update political area only when timezone changes", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		original := domain.PoliticalArea{
			Code:            "cl-rm-la-florida",
			AdminAreaLevel1: "region metropolitana de santiago",
			AdminAreaLevel2: "santiago",
			AdminAreaLevel3: "la florida",
			TimeZone:        "America/Santiago",
		}

		err = upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Get the original record
		var originalRecord table.PoliticalArea
		err = conn.DB.WithContext(ctx).
			Table("political_areas").
			Where("document_id = ?", original.DocID(ctx)).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(originalRecord.TenantID.String()).To(Equal(tenant.ID.String()))
		initialUpdatedAt := originalRecord.UpdatedAt

		// Try to update with different code but same timezone
		modifiedCode := domain.PoliticalArea{
			Code:            "cl-rm-la-florida-new",
			AdminAreaLevel1: "region metropolitana de santiago",
			AdminAreaLevel2: "santiago",
			AdminAreaLevel3: "la florida",
			TimeZone:        "America/Santiago",
		}

		err = upsert(ctx, modifiedCode)
		Expect(err).ToNot(HaveOccurred())

		var dbPoliticalArea table.PoliticalArea
		err = conn.DB.WithContext(ctx).
			Table("political_areas").
			Where("document_id = ?", original.DocID(ctx)).
			First(&dbPoliticalArea).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPoliticalArea.UpdatedAt).To(Equal(initialUpdatedAt)) // No debería actualizarse
		Expect(dbPoliticalArea.Code).To(Equal("cl-rm-la-florida"))    // Debería mantener el código original

		// Ahora actualizamos el timezone
		modifiedTimeZone := domain.PoliticalArea{
			Code:            "cl-rm-la-florida",
			AdminAreaLevel1: "region metropolitana de santiago",
			AdminAreaLevel2: "santiago",
			AdminAreaLevel3: "la florida",
			TimeZone:        "America/New_York",
		}

		err = upsert(ctx, modifiedTimeZone)
		Expect(err).ToNot(HaveOccurred())

		err = conn.DB.WithContext(ctx).
			Table("political_areas").
			Where("document_id = ?", original.DocID(ctx)).
			First(&dbPoliticalArea).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPoliticalArea.TimeZone).To(Equal("America/New_York"))
		Expect(dbPoliticalArea.UpdatedAt).To(BeNumerically(">", initialUpdatedAt)) // Debería actualizarse
		Expect(dbPoliticalArea.TenantID.String()).To(Equal(tenant.ID.String()))
	})
})
