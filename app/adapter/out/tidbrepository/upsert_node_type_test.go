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

var _ = Describe("UpsertNodeType", func() {
	// Helper function to create context with organization
	createOrgContext := func(org domain.Organization) context.Context {
		ctx := context.Background()
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, strconv.FormatInt(org.ID, 10))
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, org.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		return baggage.ContextWithBaggage(ctx, bag)
	}

	It("should insert node type if not exists", func() {
		ctx := createOrgContext(organization1)

		nodeType := domain.NodeType{
			Value: "pickup",
		}

		upsert := NewUpsertNodeType(connection)
		err := upsert(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeType table.NodeType
		err = connection.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nodeType.DocID(ctx)).
			First(&dbNodeType).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeType.Value).To(Equal("pickup"))
	})

	It("should update node type if fields are different", func() {
		ctx := createOrgContext(organization1)

		original := domain.NodeType{
			Value: "original",
		}

		upsert := NewUpsertNodeType(connection)
		err := upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.NodeType{
			Value: "updated",
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeType table.NodeType
		err = connection.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbNodeType).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeType.Value).To(Equal("updated"))
	})

	It("should not update if no fields changed", func() {
		ctx := createOrgContext(organization1)

		nodeType := domain.NodeType{
			Value: "no-change",
		}

		upsert := NewUpsertNodeType(connection)
		err := upsert(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		// Capture original record to verify timestamp doesn't change
		var originalRecord table.NodeType
		err = connection.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nodeType.DocID(ctx)).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())

		// Ejecutar nuevamente sin cambios
		err = upsert(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeType table.NodeType
		err = connection.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nodeType.DocID(ctx)).
			First(&dbNodeType).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeType.Value).To(Equal("no-change"))
		Expect(dbNodeType.UpdatedAt).To(Equal(originalRecord.UpdatedAt)) // Verify timestamp didn't change
	})

	It("should allow same node type for different organizations", func() {
		ctx1 := createOrgContext(organization1)
		ctx2 := createOrgContext(organization2)

		node1 := domain.NodeType{
			Value: "shared-type",
		}

		node2 := domain.NodeType{
			Value: "shared-type",
		}

		upsert := NewUpsertNodeType(connection)

		err := upsert(ctx1, node1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, node2)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(context.Background()).
			Table("node_types").
			Where("value = ?", node1.Value).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verify they have different document IDs
		Expect(node1.DocID(ctx1)).ToNot(Equal(node2.DocID(ctx2)))
	})

	It("should generate predictable DocID with empty value", func() {
		ctx := createOrgContext(organization1)

		nodeType := domain.NodeType{
			Value: "",
		}

		upsert := NewUpsertNodeType(connection)
		err := upsert(ctx, nodeType)
		Expect(err).ToNot(HaveOccurred())

		var dbNodeType table.NodeType
		err = connection.DB.WithContext(ctx).
			Table("node_types").
			Where("document_id = ?", nodeType.DocID(ctx)).
			First(&dbNodeType).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbNodeType.Value).To(BeEmpty())
	})

	It("should fail if database has no node_types table", func() {
		ctx := createOrgContext(organization1)

		nodeType := domain.NodeType{
			Value: "error-case",
		}

		upsert := NewUpsertNodeType(noTablesContainerConnection)
		err := upsert(ctx, nodeType)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("node_types"))
	})
})
