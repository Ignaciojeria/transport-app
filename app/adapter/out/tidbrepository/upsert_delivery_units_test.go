package tidbrepository

import (
	"context"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpsertDeliveryUnits", func() {
	var (
		conn database.ConnectionFactory
	)

	BeforeEach(func() {
		conn = connection
	})

	It("should handle empty package slice", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		upsert := NewUpsertDeliveryUnits(conn)
		err = upsert(ctx, []domain.DeliveryUnit{}, "")
		Expect(err).ToNot(HaveOccurred())
	})

	It("should insert new packages when they don't exist", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear paquetes para insertar
		package1 := domain.DeliveryUnit{
			Lpn: "PKG001",
			Dimensions: domain.Dimensions{
				Length: 10.0,
				Width:  20.0,
				Height: 30.0,
				Unit:   "cm",
			},
			Weight: domain.Weight{
				Value: 5.0,
				Unit:  "kg",
			},
			Insurance: domain.Insurance{
				UnitValue: 1000.0,
				Currency:  "USD",
			},
			Items: []domain.Item{
				{
					Sku:         "ITEM001",
					Description: "Item de prueba",
					Quantity: domain.Quantity{
						QuantityNumber: 2,
						QuantityUnit:   "unit",
					},
					Weight: domain.Weight{
						Value: 1.0,
						Unit:  "kg",
					},
				},
			},
		}

		package2 := domain.DeliveryUnit{
			Lpn: "PKG002",
			Dimensions: domain.Dimensions{
				Length: 15.0,
				Width:  25.0,
				Height: 35.0,
				Unit:   "cm",
			},
			Weight: domain.Weight{
				Value: 7.5,
				Unit:  "kg",
			},
		}

		// Insertar los paquetes
		packages := []domain.DeliveryUnit{package1, package2}
		upsert := NewUpsertDeliveryUnits(conn)
		err = upsert(ctx, packages, "")
		Expect(err).ToNot(HaveOccurred())

		// Verificar que se insertaron correctamente
		var dbPackages []table.DeliveryUnit
		err = conn.DB.WithContext(ctx).
			Table("delivery_units").
			Where("tenant_id = ?", tenant.ID).
			Find(&dbPackages).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPackages).To(HaveLen(2))

		// Verificar el primer paquete
		var dbPackage1 table.DeliveryUnit
		err = conn.DB.WithContext(ctx).
			Table("delivery_units").
			Where("document_id = ?", package1.DocID(ctx, "")).
			First(&dbPackage1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPackage1.Lpn).To(Equal("PKG001"))
		Expect(dbPackage1.TenantID.String()).To(Equal(tenant.ID.String()))

		// Verificar las dimensiones (que están en JSON)
		dimensions := dbPackage1.JSONDimensions.Map()
		Expect(dimensions.Length).To(Equal(10.0))
		Expect(dimensions.Width).To(Equal(20.0))
		Expect(dimensions.Height).To(Equal(30.0))
		Expect(dimensions.Unit).To(Equal("cm"))

		// Verificar el peso (que está en JSON)
		weight := dbPackage1.JSONWeight.Map()
		Expect(weight.Value).To(Equal(5.0))
		Expect(weight.Unit).To(Equal("kg"))

		// Verificar el seguro (que está en JSON)
		insurance := dbPackage1.JSONInsurance.Map()
		Expect(insurance.UnitValue).To(Equal(1000.0))
		Expect(insurance.Currency).To(Equal("USD"))

		// Verificar items dentro del paquete
		items := dbPackage1.JSONItems.Map()
		Expect(items).To(HaveLen(1))
		Expect(items[0].Sku).To(Equal("ITEM001"))
		Expect(items[0].Description).To(Equal("Item de prueba"))
		Expect(items[0].Quantity.QuantityNumber).To(Equal(2))
		Expect(items[0].Quantity.QuantityUnit).To(Equal("unit"))
		Expect(items[0].Weight.Value).To(Equal(1.0))
		Expect(items[0].Weight.Unit).To(Equal("kg"))
	})

	It("should update existing packages", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear paquete inicial
		originalPackage := domain.DeliveryUnit{
			Lpn: "PKG003",
			Dimensions: domain.Dimensions{
				Length: 10.0,
				Width:  20.0,
				Height: 30.0,
				Unit:   "cm",
			},
			Weight: domain.Weight{
				Value: 5.0,
				Unit:  "kg",
			},
		}

		// Insertar el paquete original
		upsert := NewUpsertDeliveryUnits(conn)
		err = upsert(ctx, []domain.DeliveryUnit{originalPackage}, "")
		Expect(err).ToNot(HaveOccurred())

		// Modificar el paquete
		modifiedPackage := originalPackage
		modifiedPackage.Dimensions.Length = 15.0
		modifiedPackage.Weight.Value = 7.5

		// Actualizar el paquete
		err = upsert(ctx, []domain.DeliveryUnit{modifiedPackage}, "")
		Expect(err).ToNot(HaveOccurred())

		// Verificar que se actualizó correctamente
		var dbPackage table.DeliveryUnit
		err = conn.DB.WithContext(ctx).
			Table("delivery_units").
			Where("document_id = ? AND tenant_id = ?", modifiedPackage.DocID(ctx, ""), tenant.ID).
			First(&dbPackage).Error
		Expect(err).ToNot(HaveOccurred())

		dimensions := dbPackage.JSONDimensions.Map()
		Expect(dimensions.Length).To(Equal(15.0))
		Expect(dimensions.Width).To(Equal(20.0))
		Expect(dimensions.Height).To(Equal(30.0))

		weight := dbPackage.JSONWeight.Map()
		Expect(weight.Value).To(Equal(7.5))
		Expect(weight.Unit).To(Equal("kg"))
	})

	It("should allow same packages for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		package1 := domain.DeliveryUnit{
			Lpn: "PKG004",
			Dimensions: domain.Dimensions{
				Length: 10.0,
				Width:  20.0,
				Height: 30.0,
				Unit:   "cm",
			},
		}

		package2 := domain.DeliveryUnit{
			Lpn: "PKG004",
			Dimensions: domain.Dimensions{
				Length: 10.0,
				Width:  20.0,
				Height: 30.0,
				Unit:   "cm",
			},
		}

		upsert := NewUpsertDeliveryUnits(conn)
		err = upsert(ctx1, []domain.DeliveryUnit{package1}, "")
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, []domain.DeliveryUnit{package2}, "")
		Expect(err).ToNot(HaveOccurred())

		// Verificar que existen dos paquetes con el mismo LPN pero en diferentes tenants
		var count int64
		err = conn.DB.WithContext(context.Background()).
			Table("delivery_units").
			Where("document_id IN ? AND tenant_id IN (?, ?)", []string{string(package1.DocID(ctx1, "")), string(package2.DocID(ctx2, ""))}, tenant1.ID, tenant2.ID).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verificar cada paquete pertenece a su respectivo tenant
		var dbPackage1, dbPackage2 table.DeliveryUnit
		err = conn.DB.WithContext(ctx1).
			Table("delivery_units").
			Where("document_id = ?", package1.DocID(ctx1, "")).
			First(&dbPackage1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPackage1.TenantID.String()).To(Equal(tenant1.ID.String()))

		err = conn.DB.WithContext(ctx2).
			Table("delivery_units").
			Where("document_id = ?", package2.DocID(ctx2, "")).
			First(&dbPackage2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPackage2.TenantID.String()).To(Equal(tenant2.ID.String()))
	})

	It("should fail if database has no delivery_units table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		package1 := domain.DeliveryUnit{
			Lpn: "PKG005",
			Dimensions: domain.Dimensions{
				Length: 10.0,
				Width:  20.0,
				Height: 30.0,
				Unit:   "cm",
			},
		}

		upsert := NewUpsertDeliveryUnits(noTablesContainerConnection)
		err = upsert(ctx, []domain.DeliveryUnit{package1}, "")

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("delivery_units"))
	})

	It("should handle mix of new and existing packages", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear un paquete inicial
		existingPackage := domain.DeliveryUnit{
			Lpn: "PKG-EXISTING",
			Dimensions: domain.Dimensions{
				Length: 10.0,
				Width:  20.0,
				Height: 30.0,
				Unit:   "cm",
			},
		}

		// Guardar el DocID del paquete existente
		existingDocID := existingPackage.DocID(ctx, "")

		// Insertar el paquete inicial
		upsert := NewUpsertDeliveryUnits(conn)
		err = upsert(ctx, []domain.DeliveryUnit{existingPackage}, "")
		Expect(err).ToNot(HaveOccurred())

		// Crear un nuevo paquete para inserción
		newPackage := domain.DeliveryUnit{
			Lpn: "PKG-NEW",
			Dimensions: domain.Dimensions{
				Length: 15.0,
				Width:  25.0,
				Height: 35.0,
				Unit:   "cm",
			},
		}

		// Crear versión actualizada del paquete existente
		updatedExistingPackage := domain.DeliveryUnit{
			Lpn: "PKG-EXISTING",
			Dimensions: domain.Dimensions{
				Length: 15.0,
				Width:  25.0,
				Height: 35.0,
				Unit:   "mm",
			},
		}

		// Verificar que el DocID del paquete actualizado coincide con el original
		updatedDocID := updatedExistingPackage.DocID(ctx, "")
		Expect(updatedDocID).To(Equal(existingDocID), "Los DocIDs deben ser iguales para actualizar")

		// Insertar ambos paquetes
		mixedPackages := []domain.DeliveryUnit{newPackage, updatedExistingPackage}
		err = upsert(ctx, mixedPackages, "")
		Expect(err).ToNot(HaveOccurred())

		// Verificar que hay dos paquetes en la base de datos
		var count int64
		err = conn.DB.WithContext(ctx).
			Table("delivery_units").
			Where("tenant_id = ?", tenant.ID).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verificar el paquete existente se actualizó
		var updatedDBPackage table.DeliveryUnit
		err = conn.DB.WithContext(ctx).
			Table("delivery_units").
			Where("document_id = ? AND tenant_id = ?", string(existingDocID), tenant.ID).
			First(&updatedDBPackage).Error
		Expect(err).ToNot(HaveOccurred())

		updatedDimensions := updatedDBPackage.JSONDimensions.Map()
		Expect(updatedDimensions.Unit).To(Equal("mm"))

		// Verificar el nuevo paquete se insertó
		var newDBPackage table.DeliveryUnit
		err = conn.DB.WithContext(ctx).
			Table("delivery_units").
			Where("lpn = ? AND tenant_id = ?", "PKG-NEW", tenant.ID).
			First(&newDBPackage).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(newDBPackage.Lpn).To(Equal("PKG-NEW"))
	})

	It("should correctly handle packages without LPN", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear un paquete sin LPN pero con items
		package1 := domain.DeliveryUnit{
			Lpn: "",
			Items: []domain.Item{
				{
					Sku:         "NO-LPN-ITEM-001",
					Description: "Item sin LPN",
					Quantity: domain.Quantity{
						QuantityNumber: 2,
						QuantityUnit:   "unit",
					},
				},
			},
		}

		// Insertar el paquete
		orderRef := "ORDER-REF-001"
		upsert := NewUpsertDeliveryUnits(conn)
		err = upsert(ctx, []domain.DeliveryUnit{package1}, orderRef)
		Expect(err).ToNot(HaveOccurred())

		// Verificar DocID generado
		expectedDocID := package1.DocID(ctx, orderRef)

		// Verificar que se insertó correctamente
		var dbPackage table.DeliveryUnit
		err = conn.DB.WithContext(ctx).
			Table("delivery_units").
			Where("document_id = ? AND tenant_id = ?", string(expectedDocID), tenant.ID).
			First(&dbPackage).Error
		Expect(err).ToNot(HaveOccurred())

		// Verificar que el LPN se mantiene vacío en la BD
		Expect(dbPackage.Lpn).To(Equal(""))

		// Verificar items
		items := dbPackage.JSONItems.Map()
		Expect(items).To(HaveLen(1))
		Expect(items[0].Sku).To(Equal("NO-LPN-ITEM-001"))
	})

	It("should correctly update packages without LPN", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear un paquete sin LPN pero con items
		initialPackage := domain.DeliveryUnit{
			Lpn: "",
			Items: []domain.Item{
				{
					Sku:         "UPDATE-NO-LPN-001",
					Description: "Item inicial sin LPN",
					Quantity: domain.Quantity{
						QuantityNumber: 1,
						QuantityUnit:   "unit",
					},
				},
			},
		}

		// Insertar el paquete
		orderRef := "ORDER-REF-002"
		upsert := NewUpsertDeliveryUnits(conn)
		err = upsert(ctx, []domain.DeliveryUnit{initialPackage}, orderRef)
		Expect(err).ToNot(HaveOccurred())

		// Guardar el DocID
		initialDocID := initialPackage.DocID(ctx, orderRef)

		// Crear un paquete actualizado (mismo SKU para generar mismo DocID)
		updatedPackage := domain.DeliveryUnit{
			Lpn: "",
			Items: []domain.Item{
				{
					Sku:         "UPDATE-NO-LPN-001",        // Mismo SKU
					Description: "Item actualizado sin LPN", // Descripción cambiada
					Quantity: domain.Quantity{
						QuantityNumber: 3, // Cantidad cambiada
						QuantityUnit:   "unit",
					},
				},
			},
		}

		// Verificar que generan el mismo DocID
		updatedDocID := updatedPackage.DocID(ctx, orderRef)
		Expect(updatedDocID).To(Equal(initialDocID))

		// Actualizar el paquete
		err = upsert(ctx, []domain.DeliveryUnit{updatedPackage}, orderRef)
		Expect(err).ToNot(HaveOccurred())

		// Verificar que se actualizó correctamente
		var dbPackage table.DeliveryUnit
		err = conn.DB.WithContext(ctx).
			Table("delivery_units").
			Where("document_id = ? AND tenant_id = ?", string(initialDocID), tenant.ID).
			First(&dbPackage).Error
		Expect(err).ToNot(HaveOccurred())

		// Verificar que el paquete se actualizó correctamente
		items := dbPackage.JSONItems.Map()
		Expect(items[0].Description).To(Equal("Item actualizado sin LPN"))
		Expect(items[0].Quantity.QuantityNumber).To(Equal(3))
	})

	It("should handle partial updates correctly", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear un paquete completo
		initialPackage := domain.DeliveryUnit{
			Lpn: "PARTIAL-UPDATE-PKG",
			Dimensions: domain.Dimensions{
				Length: 10.0,
				Width:  20.0,
				Height: 30.0,
				Unit:   "cm",
			},
			Weight: domain.Weight{
				Value: 5.0,
				Unit:  "kg",
			},
			Insurance: domain.Insurance{
				UnitValue: 1000.0,
				Currency:  "USD",
			},
			Items: []domain.Item{
				{
					Sku:         "ITEM-PARTIAL",
					Description: "Item original",
					Quantity: domain.Quantity{
						QuantityNumber: 1,
						QuantityUnit:   "unit",
					},
				},
			},
		}

		// Insertar el paquete
		upsert := NewUpsertDeliveryUnits(conn)
		err = upsert(ctx, []domain.DeliveryUnit{initialPackage}, "")
		Expect(err).ToNot(HaveOccurred())

		// Crear un paquete con solo algunos campos actualizados
		partialUpdate := domain.DeliveryUnit{
			Lpn: "PARTIAL-UPDATE-PKG", // Mismo LPN
			Weight: domain.Weight{ // Solo actualizar peso
				Value: 10.0,
				Unit:  "kg",
			},
			// Otros campos vacíos o con valores por defecto
		}

		// Actualizar el paquete
		err = upsert(ctx, []domain.DeliveryUnit{partialUpdate}, "")
		Expect(err).ToNot(HaveOccurred())

		// Verificar que solo se actualizó el campo de peso
		var dbPackage table.DeliveryUnit
		err = conn.DB.WithContext(ctx).
			Table("delivery_units").
			Where("lpn = ? AND tenant_id = ?", "PARTIAL-UPDATE-PKG", tenant.ID).
			First(&dbPackage).Error
		Expect(err).ToNot(HaveOccurred())

		// Verificar peso actualizado
		weight := dbPackage.JSONWeight.Map()
		Expect(weight.Value).To(Equal(10.0))

		// Verificar que otros campos se mantuvieron igual
		dimensions := dbPackage.JSONDimensions.Map()
		Expect(dimensions.Length).To(Equal(10.0))
		Expect(dimensions.Unit).To(Equal("cm"))

		insurance := dbPackage.JSONInsurance.Map()
		Expect(insurance.UnitValue).To(Equal(1000.0))

		items := dbPackage.JSONItems.Map()
		Expect(items[0].Description).To(Equal("Item original"))
	})

	It("should handle packages with different orderReference", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear un paquete sin LPN
		package1 := domain.DeliveryUnit{
			Lpn: "",
			Items: []domain.Item{
				{
					Sku:         "ORDER-REF-ITEM",
					Description: "Item referencia orden",
				},
			},
		}

		// Insertar con una referencia de orden
		orderRef1 := "ORDER-REF-A"
		upsert := NewUpsertDeliveryUnits(conn)
		err = upsert(ctx, []domain.DeliveryUnit{package1}, orderRef1)
		Expect(err).ToNot(HaveOccurred())

		// DocID con la primera referencia
		docID1 := package1.DocID(ctx, orderRef1)

		// Insertar el mismo paquete con otra referencia de orden
		orderRef2 := "ORDER-REF-B"
		err = upsert(ctx, []domain.DeliveryUnit{package1}, orderRef2)
		Expect(err).ToNot(HaveOccurred())

		// DocID con la segunda referencia
		docID2 := package1.DocID(ctx, orderRef2)

		// Los DocIDs deberían ser diferentes
		Expect(docID1).ToNot(Equal(docID2))

		// Verificar que ambos paquetes existen en la BD
		var count int64
		err = conn.DB.WithContext(ctx).
			Table("delivery_units").
			Where("document_id IN ? AND tenant_id = ?", []string{string(docID1), string(docID2)}, tenant.ID).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))
	})

	It("should group items into a single package when LPN is missing", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear un paquete sin LPN con múltiples items
		original := domain.DeliveryUnit{
			Lpn: "",
			Items: []domain.Item{
				{
					Sku:         "EXPLODE-ITEM-1",
					Description: "Item 1",
					Quantity: domain.Quantity{
						QuantityNumber: 1,
						QuantityUnit:   "unit",
					},
				},
				{
					Sku:         "EXPLODE-ITEM-2",
					Description: "Item 2",
					Quantity: domain.Quantity{
						QuantityNumber: 2,
						QuantityUnit:   "unit",
					},
				},
			},
		}

		orderRef := "ORDER-EXPLODE"
		upsert := NewUpsertDeliveryUnits(conn)
		err = upsert(ctx, []domain.DeliveryUnit{original}, orderRef)
		Expect(err).ToNot(HaveOccurred())

		// Verificar que se creó un solo paquete
		var count int64
		err = conn.DB.WithContext(ctx).
			Table("delivery_units").
			Where("tenant_id = ?", tenant.ID).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(1)))

		// Verificar que el paquete contiene todos los items
		var dbPackage table.DeliveryUnit
		err = conn.DB.WithContext(ctx).
			Table("delivery_units").
			Where("tenant_id = ?", tenant.ID).
			First(&dbPackage).Error
		Expect(err).ToNot(HaveOccurred())

		// Verificar que el LPN está vacío
		Expect(dbPackage.Lpn).To(Equal(""))

		// Verificar que contiene todos los items
		items := dbPackage.JSONItems.Map()
		Expect(items).To(HaveLen(2))

		// Verificar el primer item
		Expect(items[0].Sku).To(Equal("EXPLODE-ITEM-1"))
		Expect(items[0].Description).To(Equal("Item 1"))
		Expect(items[0].Quantity.QuantityNumber).To(Equal(1))
		Expect(items[0].Quantity.QuantityUnit).To(Equal("unit"))

		// Verificar el segundo item
		Expect(items[1].Sku).To(Equal("EXPLODE-ITEM-2"))
		Expect(items[1].Description).To(Equal("Item 2"))
		Expect(items[1].Quantity.QuantityNumber).To(Equal(2))
		Expect(items[1].Quantity.QuantityUnit).To(Equal("unit"))
	})
})
