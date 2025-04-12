package tidbrepository

import (
	"context"
	"strconv"

	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

var _ = Describe("UpsertOrderReferences", func() {
	var ctx context.Context

	BeforeEach(func() {
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, strconv.FormatInt(organization1.ID, 10))
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, organization1.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		ctx = baggage.ContextWithBaggage(context.Background(), bag)

		err := connection.DB.Exec("DELETE FROM order_references").Error
		Expect(err).ToNot(HaveOccurred())
	})

	It("should insert order references correctly", func() {
		order := domain.Order{
			ReferenceID: "ORD-REF-001",
			Headers: domain.Headers{
				Commerce: "commerce-1",
				Consumer: "consumer-1",
			},
			References: []domain.Reference{
				{Type: "TRACKING", Value: "123456"},
				{Type: "EXTERNAL", Value: "ABCDEF"},
			},
		}

		uor := NewUpsertOrderReferences(connection)
		err := uor(ctx, order)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(ctx).
			Table("order_references").
			Where("order_doc = ?", order.DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))
	})

	It("should replace old references with new ones", func() {
		order := domain.Order{
			ReferenceID: "ORD-REF-002",
			Headers:     domain.Headers{Commerce: "c", Consumer: "c"},
			References: []domain.Reference{
				{Type: "OLD", Value: "1"},
			},
		}

		uor := NewUpsertOrderReferences(connection)
		_ = uor(ctx, order)

		updated := order
		updated.References = []domain.Reference{
			{Type: "NEW", Value: "2"},
		}
		err := uor(ctx, updated)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(ctx).
			Table("order_references").
			Where("order_doc = ?", order.DocID(ctx)).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1)))
	})

	It("should not fail if no references are provided", func() {
		order := domain.Order{
			ReferenceID: "ORD-REF-003",
			Headers:     domain.Headers{Commerce: "x", Consumer: "y"},
			References:  []domain.Reference{},
		}

		uor := NewUpsertOrderReferences(connection)
		err := uor(ctx, order)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should fail if the order_references table is missing", func() {
		uor := NewUpsertOrderReferences(noTablesContainerConnection)
		order := domain.Order{
			ReferenceID: "ORD-REF-ERR",
			References: []domain.Reference{
				{Type: "ERR", Value: "ERR"},
			},
		}
		err := uor(ctx, order)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("order_references"))
	})
})
