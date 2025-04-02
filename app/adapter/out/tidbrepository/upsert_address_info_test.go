package tidbrepository

import (
	"context"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/paulmach/orb"
)

var _ = Describe("UpsertAddressInfo", func() {

	It("should insert addressInfo if not exists", func() {
		ctx := context.Background()

		addressInfo := domain.AddressInfo{
			AddressLine1: "Av Providencia 1234",
			District:     "Providencia",
			Province:     "Santiago",
			State:        "Metropolitana",
			ZipCode:      "7500000",
			Location:     orb.Point{-70.6506, -33.4372}, // [lon, lat]
			Organization: organization1,
		}

		upsert := NewUpsertAddressInfo(connection)
		err := upsert(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		var dbAddressInfo table.AddressInfo
		err = connection.DB.WithContext(ctx).
			Table("address_infos").
			Where("reference_id = ?", addressInfo.DocID()).
			First(&dbAddressInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbAddressInfo.AddressLine1).To(Equal("Av Providencia 1234"))
		Expect(dbAddressInfo.District).To(Equal("Providencia"))
		Expect(dbAddressInfo.State).To(Equal("Metropolitana"))
	})

	It("should update addressInfo if fields are different", func() {
		ctx := context.Background()

		original := domain.AddressInfo{
			AddressLine1: "Dirección Original",
			District:     "Las Condes",
			Province:     "Santiago",
			State:        "Metropolitana",
			ZipCode:      "7550000",
			Location:     orb.Point{-70.5768, -33.4002}, // [lon, lat]
			Organization: organization1,
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
			Organization: organization1,
		}

		err = upsert(ctx, modified)
		Expect(err).ToNot(HaveOccurred())

		var dbAddressInfo table.AddressInfo
		err = connection.DB.WithContext(ctx).
			Table("address_infos").
			Where("reference_id = ?", modified.DocID()).
			First(&dbAddressInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbAddressInfo.AddressLine1).To(Equal("Dirección Modificada"))
		Expect(dbAddressInfo.ZipCode).To(Equal("7560000"))
	})

	It("should not update if no fields changed", func() {
		ctx := context.Background()

		addressInfo := domain.AddressInfo{
			AddressLine1: "Sin Cambios",
			District:     "Vitacura",
			Province:     "Santiago",
			State:        "Metropolitana",
			ZipCode:      "7630000",
			Location:     orb.Point{-70.5906, -33.3853}, // [lon, lat]
			Organization: organization1,
		}

		upsert := NewUpsertAddressInfo(connection)

		err := upsert(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		// Capturar CreatedAt original para comparar después
		var originalRecord table.AddressInfo
		err = connection.DB.WithContext(ctx).
			Table("address_infos").
			Where("reference_id = ?", addressInfo.DocID()).
			First(&originalRecord).Error
		Expect(err).ToNot(HaveOccurred())

		// Ejecutar de nuevo sin cambios
		err = upsert(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		var dbAddressInfo table.AddressInfo
		err = connection.DB.WithContext(ctx).
			Table("address_infos").
			Where("reference_id = ?", addressInfo.DocID()).
			First(&dbAddressInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbAddressInfo.AddressLine1).To(Equal("Sin Cambios"))
		Expect(dbAddressInfo.CreatedAt).To(Equal(originalRecord.CreatedAt)) // Verificar que no se modificó
	})

	It("should allow same address info for different organizations", func() {
		ctx := context.Background()

		addressInfo1 := domain.AddressInfo{
			AddressLine1: "Multi Org",
			District:     "Santiago",
			Province:     "Santiago",
			State:        "Metropolitana",
			Organization: organization1,
		}

		addressInfo2 := domain.AddressInfo{
			AddressLine1: "Multi Org",
			District:     "Santiago",
			Province:     "Santiago",
			State:        "Metropolitana",
			Organization: organization2,
		}

		upsert := NewUpsertAddressInfo(connection)

		err := upsert(ctx, addressInfo1)
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx, addressInfo2)
		Expect(err).ToNot(HaveOccurred())

		var count int64
		err = connection.DB.WithContext(ctx).
			Table("address_infos").
			Where("address_line1 = ?", addressInfo1.AddressLine1).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))
	})

	It("should update location coordinates correctly", func() {
		ctx := context.Background()

		original := domain.AddressInfo{
			AddressLine1: "Coordenadas Test",
			District:     "Ñuñoa",
			Location:     orb.Point{-70.5975, -33.4566}, // [lon, lat]
			Organization: organization1,
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
			Where("reference_id = ?", modified.DocID()).
			First(&dbAddressInfo).Error
		Expect(err).ToNot(HaveOccurred())

		// Comparar coordenadas con una pequeña tolerancia
		Expect(dbAddressInfo.Longitude).To(BeNumerically("~", -70.6000, 0.0001))
		Expect(dbAddressInfo.Latitude).To(BeNumerically("~", -33.4600, 0.0001))
	})

	It("should correctly handle address without coordinates", func() {
		ctx := context.Background()

		addressInfo := domain.AddressInfo{
			AddressLine1: "Sin Coordenadas",
			District:     "La Florida",
			Province:     "Santiago",
			State:        "Metropolitana",
			Organization: organization1,
		}

		upsert := NewUpsertAddressInfo(connection)
		err := upsert(ctx, addressInfo)
		Expect(err).ToNot(HaveOccurred())

		var dbAddressInfo table.AddressInfo
		err = connection.DB.WithContext(ctx).
			Table("address_infos").
			Where("reference_id = ?", addressInfo.DocID()).
			First(&dbAddressInfo).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbAddressInfo.AddressLine1).To(Equal("Sin Coordenadas"))
		// Las coordenadas deberían ser cero o un valor por defecto
	})

	It("should fail if database has no address_infos table", func() {
		ctx := context.Background()

		addressInfo := domain.AddressInfo{
			AddressLine1: "Error Esperado",
			District:     "Providencia",
			Organization: organization1,
		}

		upsert := NewUpsertAddressInfo(noTablesContainerConnection)
		err := upsert(ctx, addressInfo)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("address_infos"))
	})
})
