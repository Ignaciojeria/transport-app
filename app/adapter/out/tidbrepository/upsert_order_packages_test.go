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

var _ = Describe("UpsertOrderPackages", func() {
	var ctx context.Context

	BeforeEach(func() {
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, strconv.FormatInt(organization1.ID, 10))
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, organization1.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		ctx = baggage.ContextWithBaggage(context.Background(), bag)

		err := connection.DB.Exec("DELETE FROM order_packages").Error
		Expect(err).ToNot(HaveOccurred())
	})

	It("should insert one package with LPN", func() {
		order := domain.Order{
			ReferenceID: "ORD-LPN-001",
			Headers:     domain.Headers{Commerce: "c1", Consumer: "c2"},
			Packages: []domain.Package{
				{
					Lpn: "PKG-001",
					Items: []domain.Item{
						{Sku: "ITEM-1"},
					},
				},
			},
		}

		uop := NewUpsertOrderPackages(connection)
		err := uop(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(ctx).
			Table("order_packages").
			Where("order_doc = ?", order.DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1)))
	})

	It("should replace old packages with new ones", func() {
		order := domain.Order{
			ReferenceID: "ORD-REPLACE-001",
			Headers:     domain.Headers{Commerce: "c1", Consumer: "c2"},
			Packages: []domain.Package{
				{Lpn: "PKG-OLD", Items: []domain.Item{{Sku: "OLD"}}},
			},
		}

		uop := NewUpsertOrderPackages(connection)
		err := uop(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		updated := order
		updated.Packages = []domain.Package{
			{Lpn: "PKG-NEW", Items: []domain.Item{{Sku: "NEW"}}},
		}
		err = uop(ctx, updated)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(ctx).
			Table("order_packages").
			Where("order_doc = ?", order.DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1)))
	})

	It("should not fail if no packages are provided", func() {
		order := domain.Order{
			ReferenceID: "ORD-EMPTY",
			Headers:     domain.Headers{Commerce: "x", Consumer: "y"},
			Packages:    []domain.Package{},
		}

		uop := NewUpsertOrderPackages(connection)
		err := uop(ctx, order)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should fail if the order_packages table is missing", func() {
		order := domain.Order{
			ReferenceID: "ORD-FAIL",
			Packages: []domain.Package{
				{Lpn: "PKG-ERR", Items: []domain.Item{{Sku: "E"}}},
			},
		}

		uop := NewUpsertOrderPackages(noTablesContainerConnection)
		err := uop(ctx, order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("order_packages"))
	})

	It("should insert placeholder package with a valid docID when no packages are provided", func() {
		order := domain.Order{
			ReferenceID: "ORD-PLACEHOLDER-001",
			Headers:     domain.Headers{Commerce: "a", Consumer: "b"},
			Packages:    []domain.Package{}, // no paquetes
		}

		uop := NewUpsertOrderPackages(connection)
		err := uop(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		var results []table.OrderPackage
		err = connection.DB.WithContext(ctx).
			Table("order_packages").
			Where("order_doc = ?", order.DocID(ctx)).
			Find(&results).Error
		Expect(err).ToNot(HaveOccurred())

		Expect(results).To(HaveLen(1))
		Expect(results[0].PackageDoc).ToNot(BeEmpty()) // docID generado desde pkg vacío + refID
	})

})
