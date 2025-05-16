package tidbrepository

import (
	"context"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/orb"
	"go.opentelemetry.io/otel/baggage"
)

var _ = Describe("UpsertAddressInfo", func() {

	// Helper function to create context with organization
	createOrgContext := func(org domain.Tenant) context.Context {
		ctx := context.Background()
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, org.ID.String())
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, org.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		return baggage.ContextWithBaggage(ctx, bag)
	}

	It("should insert addressInfo if not exists", func() {
		// Create context with organization1
		ctx := createOrgContext(organization1)

		addressInfo := domain.AddressInfo{
			AddressLine1: "Av Providencia 1234",
			District:     "Providencia",
			Province:     "Santiago",
			State:        "Metropolitana",
			ZipCode:      "7500000",
			Location:     orb.Point{-70.6506, -33.4372}, // [lon, lat]
		}

		upsert := NewUpsertAddressInfo(connection)
		err := upsert(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		var dbAddressInfo table.AddressInfo
		err = connection.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", addressInfo.DocID(ctx)).
			First(&dbAddressInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbAddressInfo.AddressLine1).To(Equal("Av Providencia 1234"))
		Expect(dbAddressInfo.District).To(Equal("Providencia"))
		Expect(dbAddressInfo.State).To(Equal("Metropolitana"))
	})

	It("should update addressInfo if fields are different", func() {
		ctx := createOrgContext(organization1)

		original := domain.AddressInfo{
			AddressLine1: "Dirección Original",
			District:     "Las Condes",
			Province:     "Santiago",
			State:        "Metropolitana",
			ZipCode:      "7550000",
			Location:     orb.Point{-70.5768, -33.4002}, // [lon, lat]
		}

		upsert := NewUpsertAddressInfo(connection)

		err := upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		modified := domain.AddressInfo{
			AddressLine1: "Dirección Modificada",
			District:     "Las Condes",
			Province:     "Santiago",
			State:        "Metropolitana",
			ZipCode:      "7560000", // Cambiado
			Location:     orb.Point{-70.5768, -33.4002},
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbAddressInfo table.AddressInfo
		err = connection.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbAddressInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbAddressInfo.AddressLine1).To(Equal("Dirección Modificada"))
		Expect(dbAddressInfo.ZipCode).To(Equal("7560000"))
	})

	It("should not update if no fields changed", func() {
		ctx := createOrgContext(organization1)

		addressInfo := domain.AddressInfo{
			AddressLine1: "Sin Cambios",
			District:     "Vitacura",
			Province:     "Santiago",
			State:        "Metropolitana",
			ZipCode:      "7630000",
			Location:     orb.Point{-70.5906, -33.3853}, // [lon, lat]
		}

		upsert := NewUpsertAddressInfo(connection)

		err := upsert(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		// Capturar CreatedAt original para comparar después
		var originalRecord table.AddressInfo
		err = connection.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", addressInfo.DocID(ctx)).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())

		// Ejecutar de nuevo sin cambios
		err = upsert(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		var dbAddressInfo table.AddressInfo
		err = connection.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", addressInfo.DocID(ctx)).
			First(&dbAddressInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbAddressInfo.AddressLine1).To(Equal("Sin Cambios"))
		Expect(dbAddressInfo.CreatedAt).To(Equal(originalRecord.CreatedAt)) // Verificar que no se modificó
	})

	It("should allow same address info for different organizations", func() {
		ctx1 := createOrgContext(organization1)
		ctx2 := createOrgContext(organization2)

		addressInfo1 := domain.AddressInfo{
			AddressLine1: "Multi Org",
			District:     "Santiago",
			Province:     "Santiago",
			State:        "Metropolitana",
		}

		addressInfo2 := domain.AddressInfo{
			AddressLine1: "Multi Org",
			District:     "Santiago",
			Province:     "Santiago",
			State:        "Metropolitana",
		}

		upsert := NewUpsertAddressInfo(connection)

		err := upsert(ctx1, addressInfo1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, addressInfo2)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(context.Background()).
			Table("address_infos").
			Where("address_line1 = ?", addressInfo1.AddressLine1).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))
	})

	It("should update location coordinates correctly", func() {
		ctx := createOrgContext(organization1)

		original := domain.AddressInfo{
			AddressLine1: "Coordenadas Test",
			District:     "Ñuñoa",
			Location:     orb.Point{-70.5975, -33.4566}, // [lon, lat]
		}

		upsert := NewUpsertAddressInfo(connection)
		err := upsert(ctx, original)
		Expect(err).ToNot(HaveOccurred())

		// Actualizar con nuevas coordenadas
		modified := original
		modified.Location = orb.Point{-70.6000, -33.4600} // Nuevas coordenadas

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbAddressInfo table.AddressInfo
		err = connection.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", modified.DocID(ctx)).
			First(&dbAddressInfo).Error
		Expect(err).ToNot(HaveOccurred())

		// Comparar coordenadas con una pequeña tolerancia
		Expect(dbAddressInfo.Longitude).To(BeNumerically("~", -70.6000, 0.0001))
		Expect(dbAddressInfo.Latitude).To(BeNumerically("~", -33.4600, 0.0001))
	})

	It("should correctly handle address without coordinates", func() {
		ctx := createOrgContext(organization1)

		addressInfo := domain.AddressInfo{
			AddressLine1: "Sin Coordenadas",
			District:     "La Florida",
			Province:     "Santiago",
			State:        "Metropolitana",
		}

		upsert := NewUpsertAddressInfo(connection)
		err := upsert(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		var dbAddressInfo table.AddressInfo
		err = connection.DB.WithContext(ctx).
			Table("address_infos").
			Where("document_id = ?", addressInfo.DocID(ctx)).
			First(&dbAddressInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbAddressInfo.AddressLine1).To(Equal("Sin Coordenadas"))
		// Las coordenadas deberían ser cero o un valor por defecto
	})

	It("should fail if database has no address_infos table", func() {
		ctx := createOrgContext(organization1)

		addressInfo := domain.AddressInfo{
			AddressLine1: "Error Esperado",
			District:     "Providencia",
		}

		upsert := NewUpsertAddressInfo(noTablesContainerConnection)
		err := upsert(ctx, addressInfo)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("address_infos"))
	})
})
