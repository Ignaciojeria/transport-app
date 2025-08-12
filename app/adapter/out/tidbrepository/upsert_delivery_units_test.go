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

		upsert := NewUpsertDeliveryUnits(conn, nil)
		err = upsert(ctx, []domain.DeliveryUnit{})
		Expect(err).ToNot(HaveOccurred())
	})

	It("should insert new packages when they don't exist", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear paquetes para insertar
		vol1 := int64(6000000)
		wgt1 := int64(5000)
		ins1 := int64(1000)
		package1 := domain.DeliveryUnit{
			Lpn:          "PKG001",
			SizeCategory: domain.SizeCategory{Code: "SMALL"},
			Volume:       &vol1, // 1000 * 2000 * 3000 = 6000000 cm³
			Weight:       &wgt1, // 5000 g
			Price:        &ins1, // 1000 CLP
			Items: []domain.Item{
				{
					Sku:         "ITEM001",
					Description: "Item de prueba",
					Quantity:    2,
					Weight:      1000,
				},
			},
		}

		vol2 := int64(13125000)
		wgt2 := int64(7500)
		ins2 := int64(0)
		package2 := domain.DeliveryUnit{
			Lpn:          "PKG002",
			SizeCategory: domain.SizeCategory{Code: "MEDIUM"},
			Volume:       &vol2, // 1500 * 2500 * 3500 = 13125000 cm³
			Weight:       &wgt2, // 7500 g
			Price:        &ins2, // Sin seguro
		}

		// Insertar los paquetes
		packages := []domain.DeliveryUnit{package1, package2}
		upsert := NewUpsertDeliveryUnits(conn, nil)
		err = upsert(ctx, packages)
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
			Where("document_id = ?", package1.DocID(ctx)).
			First(&dbPackage1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPackage1.Lpn).To(Equal("PKG001"))
		Expect(dbPackage1.TenantID.String()).To(Equal(tenant.ID.String()))

		// Verificar los campos simplificados
		Expect(dbPackage1.Volume).To(Equal(int64(6000000)))
		Expect(dbPackage1.Weight).To(Equal(int64(5000)))
		Expect(dbPackage1.Price).To(Equal(int64(1000)))

		// Verificar items dentro del paquete
		items := dbPackage1.JSONItems.Map()
		Expect(items).To(HaveLen(1))
		Expect(items[0].Sku).To(Equal("ITEM001"))
		Expect(items[0].Description).To(Equal("Item de prueba"))
		Expect(items[0].Quantity).To(Equal(2))
		Expect(items[0].Weight).To(Equal(int64(1000)))
	})

	It("should update existing packages", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear paquete inicial
		vol3 := int64(6000000)
		wgt3 := int64(5000)
		ins3 := int64(0)
		originalPackage := domain.DeliveryUnit{
			Lpn:       "PKG003",
			Volume:    &vol3, // 1000 * 2000 * 3000 = 6000000 cm³
			Weight: &wgt3, // 5000 g
			Price:  &ins3, // Sin seguro
		}

		// Insertar el paquete original
		upsert := NewUpsertDeliveryUnits(conn, nil)
		err = upsert(ctx, []domain.DeliveryUnit{originalPackage})
		Expect(err).ToNot(HaveOccurred())

		// Modificar el paquete
		vol4 := int64(9000000)
		wgt4 := int64(7500)
		modifiedPackage := originalPackage
		modifiedPackage.Volume = &vol4 // 1500 * 2000 * 3000 = 9000000 cm³
		modifiedPackage.Weight = &wgt4 // 7500 g

		// Actualizar el paquete
		err = upsert(ctx, []domain.DeliveryUnit{modifiedPackage})
		Expect(err).ToNot(HaveOccurred())

		// Verificar que se actualizó correctamente
		var dbPackage table.DeliveryUnit
		err = conn.DB.WithContext(ctx).
			Table("delivery_units").
			Where("document_id = ? AND tenant_id = ?", modifiedPackage.DocID(ctx), tenant.ID).
			First(&dbPackage).Error
		Expect(err).ToNot(HaveOccurred())

		// Verificar los campos simplificados
		Expect(dbPackage.Volume).To(Equal(int64(9000000)))
		Expect(dbPackage.Weight).To(Equal(int64(7500)))
	})

	It("should allow same packages for different tenants", func() {
		// Create two tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())
		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		vol := int64(6000)
		wgt := int64(0)
		ins := int64(0)
		package1 := domain.DeliveryUnit{
			Lpn:       "PKG004",
			Volume:    &vol, // 10 * 20 * 30 = 6000 cm³
			Weight:    &wgt, // Sin peso
			Price:     &ins, // Sin seguro
		}

		vol2 := int64(6000)
		wgt2 := int64(0)
		ins2 := int64(0)
		package2 := domain.DeliveryUnit{
			Lpn:       "PKG004",
			Volume:    &vol2, // 10 * 20 * 30 = 6000 cm³
			Weight:    &wgt2, // Sin peso
			Price:     &ins2, // Sin seguro
		}

		upsert := NewUpsertDeliveryUnits(conn, nil)
		err = upsert(ctx1, []domain.DeliveryUnit{package1})
		Expect(err).ToNot(HaveOccurred())

		err = upsert(ctx2, []domain.DeliveryUnit{package2})
		Expect(err).ToNot(HaveOccurred())

		// Verificar que existen dos paquetes con el mismo LPN pero en diferentes tenants
		var count int64
		err = conn.DB.WithContext(context.Background()).
			Table("delivery_units").
			Where("document_id IN ? AND tenant_id IN (?, ?)", []string{string(package1.DocID(ctx1)), string(package2.DocID(ctx2))}, tenant1.ID, tenant2.ID).
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verificar cada paquete pertenece a su respectivo tenant
		var dbPackage1, dbPackage2 table.DeliveryUnit
		err = conn.DB.WithContext(ctx1).
			Table("delivery_units").
			Where("document_id = ?", package1.DocID(ctx1)).
			First(&dbPackage1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPackage1.TenantID.String()).To(Equal(tenant1.ID.String()))

		err = conn.DB.WithContext(ctx2).
			Table("delivery_units").
			Where("document_id = ?", package2.DocID(ctx2)).
			First(&dbPackage2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPackage2.TenantID.String()).To(Equal(tenant2.ID.String()))
	})

	It("should fail if database has no delivery_units table", func() {
		// Create a new tenant for this test
		_, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		vol3 := int64(6000)
		wgt3 := int64(1000)
		ins3 := int64(0)
		package1 := domain.DeliveryUnit{
			Lpn:       "PKG005",
			Volume:    &vol3, // 10 * 20 * 30 = 6000 cm³
			Weight:    &wgt3, // 1000 g
			Price:     &ins3, // Sin seguro
		}

		upsert := NewUpsertDeliveryUnits(noTablesContainerConnection, nil)
		err = upsert(ctx, []domain.DeliveryUnit{package1})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("delivery_units"))
	})

	It("should handle mix of new and existing packages", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear un paquete inicial
		vol4 := int64(6000)
		wgt4 := int64(1000)
		ins4 := int64(0)
		existingPackage := domain.DeliveryUnit{
			Lpn:       "PKG-EXISTING",
			Volume:    &vol4, // 10 * 20 * 30 = 6000 cm³
			Weight:    &wgt4, // 1000 g
			Price:     &ins4, // Sin seguro
		}

		// Guardar el DocID del paquete existente
		existingDocID := existingPackage.DocID(ctx)

		// Insertar el paquete inicial
		upsert := NewUpsertDeliveryUnits(conn, nil)
		err = upsert(ctx, []domain.DeliveryUnit{existingPackage})
		Expect(err).ToNot(HaveOccurred())

		// Crear un nuevo paquete para inserción
		newVol := int64(13125)
		newWgt := int64(2000)
		newIns := int64(500)
		newPackage := domain.DeliveryUnit{
			Lpn:       "PKG-NEW",
			Volume:    &newVol, // 15 * 25 * 35 = 13125 cm³
			Weight:    &newWgt, // 2000 g
			Price:     &newIns, // 500 CLP
		}

		// Crear versión actualizada del paquete existente
		updatedVol := int64(13125)
		updatedWgt := int64(2000)
		updatedIns := int64(500)
		updatedExistingPackage := domain.DeliveryUnit{
			Lpn:       "PKG-EXISTING",
			Volume:    &updatedVol, // 15 * 25 * 35 = 13125 cm³
			Weight:    &updatedWgt, // 2000 g
			Price:     &updatedIns, // 500 CLP
		}

		// Verificar que el DocID del paquete actualizado coincide con el original
		updatedDocID := updatedExistingPackage.DocID(ctx)
		Expect(updatedDocID).To(Equal(existingDocID), "Los DocIDs deben ser iguales para actualizar")

		// Insertar ambos paquetes
		mixedPackages := []domain.DeliveryUnit{newPackage, updatedExistingPackage}
		err = upsert(ctx, mixedPackages)
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

		// Verificar que los campos se actualizaron correctamente
		Expect(updatedDBPackage.Volume).To(Equal(int64(13125)))
		Expect(updatedDBPackage.Weight).To(Equal(int64(2000)))
		Expect(updatedDBPackage.Price).To(Equal(int64(500)))

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
					Quantity:    2,
				},
			},
		}

		// Insertar el paquete
		//	orderRef := "ORDER-REF-001"
		upsert := NewUpsertDeliveryUnits(conn, nil)
		err = upsert(ctx, []domain.DeliveryUnit{package1})
		Expect(err).ToNot(HaveOccurred())

		// Verificar DocID generado
		expectedDocID := package1.DocID(ctx)

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
					Quantity:    1,
				},
			},
		}

		// Insertar el paquete
		//	orderRef := "ORDER-REF-002"
		upsert := NewUpsertDeliveryUnits(conn, nil)
		err = upsert(ctx, []domain.DeliveryUnit{initialPackage})
		Expect(err).ToNot(HaveOccurred())

		// Guardar el DocID
		initialDocID := initialPackage.DocID(ctx)

		// Crear un paquete actualizado (mismo SKU para generar mismo DocID)
		updatedPackage := domain.DeliveryUnit{
			Lpn: "",
			Items: []domain.Item{
				{
					Sku:         "UPDATE-NO-LPN-001",        // Mismo SKU
					Description: "Item actualizado sin LPN", // Descripción cambiada
					Quantity:    3,                          // Cantidad cambiada
				},
			},
		}

		// Verificar que generan el mismo DocID
		updatedDocID := updatedPackage.DocID(ctx)
		Expect(updatedDocID).To(Equal(initialDocID))

		// Actualizar el paquete
		err = upsert(ctx, []domain.DeliveryUnit{updatedPackage})
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
		Expect(items[0].Quantity).To(Equal(3))
	})

	It("should handle partial updates correctly", func() {
		// Create a new tenant for this test
		tenant, ctx, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		// Crear un paquete completo
		vol := int64(6000)
		wgt := int64(5000)
		ins := int64(1000)
		initialPackage := domain.DeliveryUnit{
			Lpn:       "PARTIAL-UPDATE-PKG",
			Volume:    &vol, // 10 * 20 * 30 = 6000 cm³
			Weight:    &wgt, // 5000 g
			Price:     &ins, // 1000 CLP
			Items: []domain.Item{
				{
					Sku:         "ITEM-PARTIAL",
					Description: "Item original",
					Quantity:    1,
				},
			},
		}

		// Insertar el paquete
		upsert := NewUpsertDeliveryUnits(conn, nil)
		err = upsert(ctx, []domain.DeliveryUnit{initialPackage})
		Expect(err).ToNot(HaveOccurred())

		// Crear un paquete con solo algunos campos actualizados
		wgt2 := int64(10000)
		partialUpdate := domain.DeliveryUnit{
			Lpn:    "PARTIAL-UPDATE-PKG", // Mismo LPN
			Weight: &wgt2,                // Solo actualizar peso
			// Otros campos vacíos o con valores por defecto
		}

		// Actualizar el paquete
		err = upsert(ctx, []domain.DeliveryUnit{partialUpdate})
		Expect(err).ToNot(HaveOccurred())

		// Verificar que solo se actualizó el campo de peso
		var dbPackage table.DeliveryUnit
		err = conn.DB.WithContext(ctx).
			Table("delivery_units").
			Where("lpn = ? AND tenant_id = ?", "PARTIAL-UPDATE-PKG", tenant.ID).
			First(&dbPackage).Error
		Expect(err).ToNot(HaveOccurred())

		// Verificar peso actualizado
		Expect(dbPackage.Weight).To(Equal(int64(10000)))

		// Verificar que otros campos se mantuvieron igual
		Expect(dbPackage.Volume).To(Equal(int64(6000)))
		Expect(dbPackage.Price).To(Equal(int64(1000)))

		items := dbPackage.JSONItems.Map()
		Expect(items[0].Description).To(Equal("Item original"))
	})

	It("should handle packages with different orderReference", func() {
		// Create two different tenants for this test
		tenant1, ctx1, err := CreateTestTenant(context.Background(), conn)
		Expect(err).ToNot(HaveOccurred())

		tenant2, ctx2, err := CreateTestTenant(context.Background(), conn)
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

		// Insertar con el primer tenant
		upsert := NewUpsertDeliveryUnits(conn, nil)
		err = upsert(ctx1, []domain.DeliveryUnit{package1})
		Expect(err).ToNot(HaveOccurred())

		// DocID con el primer tenant
		docID1 := package1.DocID(ctx1)

		// Insertar el mismo paquete con el segundo tenant
		err = upsert(ctx2, []domain.DeliveryUnit{package1})
		Expect(err).ToNot(HaveOccurred())

		// DocID con el segundo tenant
		docID2 := package1.DocID(ctx2)

		// Los DocIDs deberían ser diferentes
		Expect(docID1).ToNot(Equal(docID2))

		// Verificar que ambos paquetes existen en la BD con sus respectivos tenants
		var count1 int64
		err = conn.DB.WithContext(ctx1).
			Table("delivery_units").
			Where("document_id = ? AND tenant_id = ?", string(docID1), tenant1.ID).
			Count(&count1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count1).To(Equal(int64(1)))

		var count2 int64
		err = conn.DB.WithContext(ctx2).
			Table("delivery_units").
			Where("document_id = ? AND tenant_id = ?", string(docID2), tenant2.ID).
			Count(&count2).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count2).To(Equal(int64(1)))
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
					Quantity:    1,
				},
				{
					Sku:         "EXPLODE-ITEM-2",
					Description: "Item 2",
					Quantity:    2,
				},
			},
		}

		//orderRef := "ORDER-EXPLODE"
		upsert := NewUpsertDeliveryUnits(conn, nil)
		err = upsert(ctx, []domain.DeliveryUnit{original})
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
		Expect(items[0].Quantity).To(Equal(1))

		// Verificar el segundo item
		Expect(items[1].Sku).To(Equal("EXPLODE-ITEM-2"))
		Expect(items[1].Description).To(Equal("Item 2"))
		Expect(items[1].Quantity).To(Equal(2))
	})
})
