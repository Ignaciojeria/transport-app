package tidbrepository

import (
	"context"
	"strconv"
	"time"

	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"
	"transport-app/app/shared/sharedcontext"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/baggage"
)

var _ = Describe("UpsertPackages", func() {
	var (
		ctx1, ctx2 context.Context
	)

	// Helper function to create context with organization
	createOrgContext := func(org domain.Organization) context.Context {
		ctx := context.Background()
		orgIDMember, _ := baggage.NewMember(sharedcontext.BaggageTenantID, strconv.FormatInt(org.ID, 10))
		countryMember, _ := baggage.NewMember(sharedcontext.BaggageTenantCountry, org.Country.String())
		bag, _ := baggage.New(orgIDMember, countryMember)
		return baggage.ContextWithBaggage(ctx, bag)
	}

	BeforeEach(func() {
		// Create contexts with different organizations
		ctx1 = createOrgContext(organization1)
		ctx2 = createOrgContext(organization2)

		// Limpia la tabla antes de cada test
		err := connection.DB.Exec("DELETE FROM packages").Error
		Expect(err).ToNot(HaveOccurred())
	})

	It("should handle empty package slice", func() {
		upsert := NewUpsertPackages(connection)
		err := upsert(ctx1, []domain.Package{})
		Expect(err).ToNot(HaveOccurred())
	})

	It("should insert new packages when they don't exist", func() {
		// Crear paquetes para insertar
		package1 := domain.Package{
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

		package2 := domain.Package{
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
		packages := []domain.Package{package1, package2}
		upsert := NewUpsertPackages(connection)
		err := upsert(ctx1, packages)
		Expect(err).ToNot(HaveOccurred())

		// Verificar que se insertaron correctamente
		var dbPackages []table.Package
		err = connection.DB.WithContext(ctx1).
			Table("packages").
			Find(&dbPackages).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPackages).To(HaveLen(2))

		// Verificar el primer paquete
		var dbPackage1 table.Package
		err = connection.DB.WithContext(ctx1).
			Table("packages").
			Where("lpn = ?", "PKG001").
			First(&dbPackage1).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbPackage1.Lpn).To(Equal("PKG001"))
		Expect(dbPackage1.OrganizationID).To(Equal(organization1.ID))

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
		// Crear un paquete inicial
		initialPackage := domain.Package{
			Lpn: "PKG-UPDATE",
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
					Description: "Item inicial",
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

		// Importante: Guardamos el DocID del paquete inicial para usarlo en la actualización
		initialDocID := initialPackage.DocID(ctx1)

		// Insertar el paquete inicial
		upsert := NewUpsertPackages(connection)
		err := upsert(ctx1, []domain.Package{initialPackage})
		Expect(err).ToNot(HaveOccurred())

		// Obtener el registro creado y su timestamp
		var initialDBPackage table.Package
		err = connection.DB.WithContext(ctx1).
			Table("packages").
			Where("document_id = ?", string(initialDocID)).
			First(&initialDBPackage).Error
		Expect(err).ToNot(HaveOccurred())
		initialCreatedAt := initialDBPackage.CreatedAt
		initialID := initialDBPackage.ID

		// Esperar un momento para asegurar que el timestamp de actualización sea diferente
		time.Sleep(1 * time.Millisecond)

		// Crear versión actualizada del paquete
		updatedPackage := domain.Package{
			Lpn: "PKG-UPDATE", // Mismo LPN para que se actualice
			Dimensions: domain.Dimensions{
				Length: 15.0, // Cambiar dimensiones
				Width:  25.0,
				Height: 35.0,
				Unit:   "mm", // Cambiar unidad
			},
			Weight: domain.Weight{
				Value: 7.5,  // Cambiar peso
				Unit:  "lb", // Cambiar unidad
			},
			Insurance: domain.Insurance{
				UnitValue: 2000.0, // Cambiar valor de seguro
				Currency:  "EUR",  // Cambiar moneda
			},
			Items: []domain.Item{
				{
					Sku:         "ITEM002", // Cambiar referencia
					Description: "Item actualizado",
					Quantity: domain.Quantity{
						QuantityNumber: 3,
						QuantityUnit:   "box",
					},
					Weight: domain.Weight{
						Value: 1.5,
						Unit:  "lb",
					},
				},
			},
		}

		// Verificar que generan el mismo DocID
		updatedDocID := updatedPackage.DocID(ctx1)
		Expect(updatedDocID).To(Equal(initialDocID), "Los DocIDs deben ser iguales para la actualización")

		// Actualizar el paquete
		err = upsert(ctx1, []domain.Package{updatedPackage})
		Expect(err).ToNot(HaveOccurred())

		// Verificar que se actualizó correctamente
		var updatedDBPackage table.Package
		err = connection.DB.WithContext(ctx1).
			Table("packages").
			Where("document_id = ?", string(initialDocID)).
			First(&updatedDBPackage).Error
		Expect(err).ToNot(HaveOccurred())

		// Verificar que mantiene el mismo ID y CreatedAt
		Expect(updatedDBPackage.ID).To(Equal(initialID))
		Expect(updatedDBPackage.CreatedAt).To(Equal(initialCreatedAt))

		// Verificar que los campos JSON se actualizaron correctamente
		updatedDimensions := updatedDBPackage.JSONDimensions.Map()
		Expect(updatedDimensions.Length).To(Equal(15.0))
		Expect(updatedDimensions.Width).To(Equal(25.0))
		Expect(updatedDimensions.Height).To(Equal(35.0))
		Expect(updatedDimensions.Unit).To(Equal("mm"))

		updatedWeight := updatedDBPackage.JSONWeight.Map()
		Expect(updatedWeight.Value).To(Equal(7.5))
		Expect(updatedWeight.Unit).To(Equal("lb"))

		updatedInsurance := updatedDBPackage.JSONInsurance.Map()
		Expect(updatedInsurance.UnitValue).To(Equal(2000.0))
		Expect(updatedInsurance.Currency).To(Equal("EUR"))

		// Verificar que los items se actualizaron
		updatedItems := updatedDBPackage.JSONItems.Map()
		Expect(updatedItems).To(HaveLen(1))
		Expect(updatedItems[0].Sku).To(Equal("ITEM002"))
		Expect(updatedItems[0].Description).To(Equal("Item actualizado"))
		Expect(updatedItems[0].Quantity.QuantityNumber).To(Equal(3))
		Expect(updatedItems[0].Quantity.QuantityUnit).To(Equal("box"))
		Expect(updatedItems[0].Weight.Value).To(Equal(1.5))
		Expect(updatedItems[0].Weight.Unit).To(Equal("lb"))
	})

	It("should handle mix of new and existing packages", func() {
		// Crear un paquete inicial
		existingPackage := domain.Package{
			Lpn: "PKG-EXISTING",
			Dimensions: domain.Dimensions{
				Length: 10.0,
				Width:  20.0,
				Height: 30.0,
				Unit:   "cm",
			},
		}

		// Guardar el DocID del paquete existente
		existingDocID := existingPackage.DocID(ctx1)

		// Insertar el paquete inicial
		upsert := NewUpsertPackages(connection)
		err := upsert(ctx1, []domain.Package{existingPackage})
		Expect(err).ToNot(HaveOccurred())

		// Crear un nuevo paquete para inserción
		newPackage := domain.Package{
			Lpn: "PKG-NEW",
			Dimensions: domain.Dimensions{
				Length: 15.0,
				Width:  25.0,
				Height: 35.0,
				Unit:   "cm",
			},
		}

		// Crear versión actualizada del paquete existente
		updatedExistingPackage := domain.Package{
			Lpn: "PKG-EXISTING",
			Dimensions: domain.Dimensions{
				Length: 15.0,
				Width:  25.0,
				Height: 35.0,
				Unit:   "mm",
			},
		}

		// Verificar que el DocID del paquete actualizado coincide con el original
		updatedDocID := updatedExistingPackage.DocID(ctx1)
		Expect(updatedDocID).To(Equal(existingDocID), "Los DocIDs deben ser iguales para actualizar")

		// Insertar ambos paquetes
		mixedPackages := []domain.Package{newPackage, updatedExistingPackage}
		err = upsert(ctx1, mixedPackages)
		Expect(err).ToNot(HaveOccurred())

		// Verificar que hay dos paquetes en la base de datos
		var count int64
		err = connection.DB.WithContext(ctx1).
			Table("packages").
			Count(&count).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(int64(2)))

		// Verificar el paquete existente se actualizó
		var updatedDBPackage table.Package
		err = connection.DB.WithContext(ctx1).
			Table("packages").
			Where("document_id = ?", string(existingDocID)).
			First(&updatedDBPackage).Error
		Expect(err).ToNot(HaveOccurred())

		updatedDimensions := updatedDBPackage.JSONDimensions.Map()
		Expect(updatedDimensions.Unit).To(Equal("mm"))

		// Verificar el nuevo paquete se insertó
		var newDBPackage table.Package
		err = connection.DB.WithContext(ctx1).
			Table("packages").
			Where("lpn = ?", "PKG-NEW").
			First(&newDBPackage).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(newDBPackage.Lpn).To(Equal("PKG-NEW"))
	})

	It("should allow same LPN in different organizations", func() {
		// Crear paquete para organización 1
		package1 := domain.Package{
			Lpn: "PKG-MULTIORG",
			Dimensions: domain.Dimensions{
				Length: 10.0,
				Width:  20.0,
				Height: 30.0,
				Unit:   "cm",
			},
		}

		// Crear paquete para organización 2 con el mismo LPN
		package2 := domain.Package{
			Lpn: "PKG-MULTIORG",
			Dimensions: domain.Dimensions{
				Length: 15.0,
				Width:  25.0,
				Height: 35.0,
				Unit:   "cm",
			},
		}

		// Insertar paquete para org1
		upsert := NewUpsertPackages(connection)
		err := upsert(ctx1, []domain.Package{package1})
		Expect(err).ToNot(HaveOccurred())

		// Insertar paquete para org2
		err = upsert(ctx2, []domain.Package{package2})
		Expect(err).ToNot(HaveOccurred())

		// Verificar que hay dos paquetes en la base de datos
		var packages []table.Package
		err = connection.DB.WithContext(context.Background()).
			Table("packages").
			Where("lpn = ?", "PKG-MULTIORG").
			Find(&packages).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(packages).To(HaveLen(2))

		// Verificar que tienen diferentes organizaciones
		orgs := map[int64]bool{}
		for _, pkg := range packages {
			orgs[pkg.OrganizationID] = true
		}
		Expect(orgs).To(HaveLen(2))
		Expect(orgs[organization1.ID]).To(BeTrue())
		Expect(orgs[organization2.ID]).To(BeTrue())

		// Verify they have different document IDs
		Expect(package1.DocID(ctx1)).ToNot(Equal(package2.DocID(ctx2)))
	})

	It("should fail when database has no packages table", func() {
		package1 := domain.Package{
			Lpn: "PKG-ERROR",
		}

		upsert := NewUpsertPackages(noTablesContainerConnection)
		err := upsert(ctx1, []domain.Package{package1})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("packages"))
	})

	It("should fail when saving packages if the table does not exist", func() {
		// Crear un paquete válido
		pkg := domain.Package{
			Lpn: "PKG-NOTABLE",
			Dimensions: domain.Dimensions{
				Length: 10.0,
				Width:  10.0,
				Height: 10.0,
				Unit:   "cm",
			},
		}

		// Usar conexión sin tablas
		upsert := NewUpsertPackages(noTablesContainerConnection)
		err := upsert(ctx1, []domain.Package{pkg})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("packages"))
	})
})
