package domain

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Package", func() {
	var ctx1, ctx2 context.Context

	BeforeEach(func() {
		ctx1 = buildCtx("org1", "CL")
		ctx2 = buildCtx("org2", "AR")
	})

	Describe("DocID", func() {
		It("should generate unique ID based on context and Lpn when Lpn is present", func() {
			package1 := DeliveryUnit{
				Lpn: "PKG001",
			}
			package2 := DeliveryUnit{
				Lpn: "PKG002",
			}

			Expect(package1.DocID(ctx1)).To(Equal(HashByTenant(ctx1, "PKG001")))
			Expect(package1.DocID(ctx1)).ToNot(Equal(package2.DocID(ctx1)))
			Expect(package1.DocID(ctx1)).ToNot(Equal(package1.DocID(ctx2)))
		})

		It("should generate ID based on reference and sorted items when Lpn is empty", func() {
			pkg := DeliveryUnit{
				Lpn:            "",
				noLPNReference: "REF001",
				Items: []Item{
					{Sku: "SKU002"},
					{Sku: "SKU001"},
				},
			}
			expectedInputs := []string{"REF001", "SKU001", "SKU002"}
			Expect(pkg.DocID(ctx1)).To(Equal(HashByTenant(ctx1, expectedInputs...)))
		})

		It("should generate different IDs for packages with same reference but different items", func() {
			pkg1 := DeliveryUnit{
				Lpn: "",
				Items: []Item{
					{Sku: "SKU001"},
					{Sku: "SKU002"},
				},
			}

			pkg2 := DeliveryUnit{
				Lpn: "",
				Items: []Item{
					{Sku: "SKU001"},
					{Sku: "SKU003"}, // Diferente SKU
				},
			}

			id1 := pkg1.DocID(ctx1)
			id2 := pkg2.DocID(ctx1)

			Expect(id1).ToNot(Equal(id2))
		})

		It("should generate consistent IDs for same reference and items", func() {
			pkg1 := DeliveryUnit{
				Lpn: "",
				Items: []Item{
					{Sku: "SKU001"},
					{Sku: "SKU002"},
				},
			}

			pkg2 := DeliveryUnit{
				Lpn: "",
				Items: []Item{
					{Sku: "SKU001"},
					{Sku: "SKU002"},
				},
			}

			id1 := pkg1.DocID(ctx1)
			id2 := pkg2.DocID(ctx1)

			Expect(id1).To(Equal(id2))
		})

		It("should handle the case with no LPN and no items", func() {
			pkg := DeliveryUnit{
				Lpn:            "",
				noLPNReference: "REF001",
				Items:          []Item{},
			}

			expected := HashByTenant(ctx1, "REF001")
			Expect(pkg.DocID(ctx1)).To(Equal(expected))
		})

	})

	Describe("SearchPackageByLpn", func() {
		var packages []DeliveryUnit

		BeforeEach(func() {
			packages = []DeliveryUnit{
				{Lpn: "PKG001"},
				{Lpn: "PKG002"},
				{Lpn: "PKG003"},
			}
		})

		It("should find package by lpn when it exists", func() {
			result := SearchPackageByLpn(packages, "PKG002")
			Expect(result.Lpn).To(Equal("PKG002"))
		})

		It("should return empty package when lpn doesn't exist", func() {
			result := SearchPackageByLpn(packages, "NONEXISTENT")
			Expect(result).To(Equal(DeliveryUnit{}))
		})

		It("should return first matching package when multiple matches exist", func() {
			// Agregar un duplicado con el mismo Lpn
			duplicatePackages := append(packages, DeliveryUnit{Lpn: "PKG001"})

			result := SearchPackageByLpn(duplicatePackages, "PKG001")
			Expect(result.Lpn).To(Equal("PKG001"))
		})

		It("should handle empty package slice", func() {
			result := SearchPackageByLpn([]DeliveryUnit{}, "PKG001")
			Expect(result).To(Equal(DeliveryUnit{}))
		})

		It("should handle case-sensitive lpn comparison", func() {
			// Agregar paquete con LPN en minúsculas
			packagesWithCase := append(packages, DeliveryUnit{Lpn: "pkg001"})

			// Debería devolver el original (PKG001) y no el nuevo (pkg001)
			result := SearchPackageByLpn(packagesWithCase, "PKG001")
			Expect(result.Lpn).To(Equal("PKG001"))

			// Debería encontrar el lpn en minúsculas
			result = SearchPackageByLpn(packagesWithCase, "pkg001")
			Expect(result.Lpn).To(Equal("pkg001"))
		})
	})

	Describe("UpdateIfChanged", func() {
		var basePackage DeliveryUnit

		BeforeEach(func() {
			vol := int64(6000)
			wgt := int64(5000)
			ins := int64(1000)
			basePackage = DeliveryUnit{
				Lpn:    "PKG-TEST",
				Volume: &vol, // 10 * 20 * 30 = 6000 cm³
				Weight: &wgt, // 5 kg = 5000 g
				Price:  &ins, // 1000 CLP (simplified)
				Items: []Item{
					{
						Sku:         "ITEM001",
						Description: "Item 1 Description",
						Quantity:    2,
						Weight:      1000,
					},
					{
						Sku:         "ITEM002",
						Description: "Item 2 Description",
						Quantity:    1,
						Weight:      500,
					},
				},
			}
		})

		It("should update Lpn", func() {
			newPackage := DeliveryUnit{
				Lpn: "PKG-UPDATED",
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			Expect(changed).To(BeTrue())
			Expect(updated.Lpn).To(Equal("PKG-UPDATED"))
			// Verificar que otros campos se mantienen igual
			Expect(updated.Volume).To(Equal(basePackage.Volume))
			Expect(updated.Weight).To(Equal(basePackage.Weight))
			Expect(updated.Price).To(Equal(basePackage.Price))
			Expect(updated.Items).To(Equal(basePackage.Items))
		})

		It("should update Volume", func() {
			vol := int64(13125)
			newPackage := DeliveryUnit{
				Volume: &vol, // 15 * 25 * 35 = 13125 cm³
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			Expect(changed).To(BeTrue())
			Expect(updated.Volume).To(Equal(newPackage.Volume))
			// Verificar que otros campos se mantienen igual
			Expect(updated.Lpn).To(Equal(basePackage.Lpn))
			Expect(updated.Weight).To(Equal(basePackage.Weight))
			Expect(updated.Price).To(Equal(basePackage.Price))
			Expect(updated.Items).To(Equal(basePackage.Items))
		})

		It("should update Weight", func() {
			wgt := int64(7500)
			newPackage := DeliveryUnit{
				Weight: &wgt, // 7500 g
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			Expect(changed).To(BeTrue())
			Expect(updated.Weight).To(Equal(newPackage.Weight))
			// Verificar que otros campos se mantienen igual
			Expect(updated.Lpn).To(Equal(basePackage.Lpn))
			Expect(updated.Volume).To(Equal(basePackage.Volume))
			Expect(updated.Price).To(Equal(basePackage.Price))
			Expect(updated.Items).To(Equal(basePackage.Items))
		})

		It("should update Insurance", func() {
			ins := int64(2000)
			newPackage := DeliveryUnit{
				Price: &ins, // 2000 CLP
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			Expect(changed).To(BeTrue())
			Expect(updated.Price).To(Equal(newPackage.Price))
			// Verificar que otros campos se mantienen igual
			Expect(updated.Lpn).To(Equal(basePackage.Lpn))
			Expect(updated.Volume).To(Equal(basePackage.Volume))
			Expect(updated.Weight).To(Equal(basePackage.Weight))
			Expect(updated.Items).To(Equal(basePackage.Items))
		})

		It("should update Items", func() {
			newPackage := DeliveryUnit{
				Items: []Item{
					{
						Sku:         "ITEM003",
						Description: "New Item 3",
						Quantity:    3,
						Weight:      1500,
					},
					{
						Sku:         "ITEM004",
						Description: "New Item 4",
						Quantity:    4,
						Weight:      2000,
					},
				},
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			Expect(changed).To(BeTrue())
			Expect(updated.Items).To(Equal(newPackage.Items))
			// Verificar que otros campos se mantienen igual
			Expect(updated.Lpn).To(Equal(basePackage.Lpn))
			Expect(updated.Volume).To(Equal(basePackage.Volume))
			Expect(updated.Weight).To(Equal(basePackage.Weight))
			Expect(updated.Price).To(Equal(basePackage.Price))
		})

		It("should not update fields when new values are empty", func() {
			newPackage := DeliveryUnit{
				Lpn: "",
				// Todos los demás campos están vacíos o con sus valores por defecto
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			// Nada debería cambiar
			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(basePackage))
		})

		It("should update multiple fields at once", func() {
			wgt := int64(8000)
			newPackage := DeliveryUnit{
				Lpn:    "PKG-MULTI-UPDATE",
				Weight: &wgt, // 8000 g
				Items: []Item{
					{
						Sku:         "ITEM005",
						Description: "New Item for Multi-update",
						Quantity:    5,
						Weight:      3000,
					},
				},
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			Expect(changed).To(BeTrue())
			Expect(updated.Lpn).To(Equal("PKG-MULTI-UPDATE"))
			Expect(updated.Weight).To(Equal(newPackage.Weight))
			Expect(updated.Items).To(Equal(newPackage.Items))
			// Estos campos no deberían cambiar
			Expect(updated.Volume).To(Equal(basePackage.Volume))
			Expect(updated.Price).To(Equal(basePackage.Price))
		})

		It("should handle empty Items array", func() {
			// Primero confirmar que tenemos items iniciales
			Expect(basePackage.Items).ToNot(BeEmpty())

			// Intentar actualizar con un array vacío
			newPackage := DeliveryUnit{
				Items: []Item{},
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			// Los items deberían mantenerse sin cambios, ya que el array vacío
			// no debería sobrescribir los valores existentes según la lógica del método
			Expect(changed).To(BeFalse())
			Expect(updated.Items).To(Equal(basePackage.Items))
		})

		It("should replace items partially when some are provided", func() {
			// Verificar el número inicial de items
			Expect(len(basePackage.Items)).To(Equal(2))

			// Actualizar con un solo item
			newPackage := DeliveryUnit{
				Items: []Item{
					{
						Sku:         "NEW-SINGLE-ITEM",
						Description: "Single Item Update",
						Quantity:    10,
					},
				},
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			// Debería reemplazar todos los items con el nuevo array
			Expect(changed).To(BeTrue())
			Expect(len(updated.Items)).To(Equal(1))
			Expect(updated.Items[0].Sku).To(Equal("NEW-SINGLE-ITEM"))
		})

		It("should update Volume to zero", func() {
			baseVal := int64(1000)
			zero := int64(0)
			basePackage := DeliveryUnit{
				Volume: &baseVal,
			}
			newPackage := DeliveryUnit{
				Volume: &zero,
			}
			updated, changed := basePackage.UpdateIfChanged(newPackage)
			Expect(changed).To(BeTrue())
			Expect(updated.Volume).ToNot(BeNil())
			Expect(*updated.Volume).To(Equal(int64(0)))
		})

		It("should not update Volume if new value is nil", func() {
			baseVal := int64(1000)
			basePackage := DeliveryUnit{
				Volume: &baseVal,
			}
			newPackage := DeliveryUnit{
				Volume: nil,
			}
			updated, changed := basePackage.UpdateIfChanged(newPackage)
			Expect(changed).To(BeFalse())
			Expect(updated.Volume).ToNot(BeNil())
			Expect(*updated.Volume).To(Equal(int64(1000)))
		})

		It("should update Weight to zero", func() {
			wgt := int64(500)
			zero := int64(0)
			basePackage := DeliveryUnit{
				Weight: &wgt,
			}
			newPackage := DeliveryUnit{
				Weight: &zero,
			}
			updated, changed := basePackage.UpdateIfChanged(newPackage)
			Expect(changed).To(BeTrue())
			Expect(updated.Weight).ToNot(BeNil())
			Expect(*updated.Weight).To(Equal(int64(0)))
		})

		It("should not update Weight if new value is nil", func() {
			wgt := int64(500)
			basePackage := DeliveryUnit{
				Weight: &wgt,
			}
			newPackage := DeliveryUnit{
				Weight: nil,
			}
			updated, changed := basePackage.UpdateIfChanged(newPackage)
			Expect(changed).To(BeFalse())
			Expect(updated.Weight).ToNot(BeNil())
			Expect(*updated.Weight).To(Equal(int64(500)))
		})

		It("should update Insurance to zero", func() {
			ins := int64(200)
			zero := int64(0)
			basePackage := DeliveryUnit{
				Price: &ins,
			}
			newPackage := DeliveryUnit{
				Price: &zero,
			}
			updated, changed := basePackage.UpdateIfChanged(newPackage)
			Expect(changed).To(BeTrue())
			Expect(updated.Price).ToNot(BeNil())
			Expect(*updated.Price).To(Equal(int64(0)))
		})

		It("should not update Insurance if new value is nil", func() {
			ins := int64(200)
			basePackage := DeliveryUnit{
				Price: &ins,
			}
			newPackage := DeliveryUnit{
				Price: nil,
			}
			updated, changed := basePackage.UpdateIfChanged(newPackage)
			Expect(changed).To(BeFalse())
			Expect(updated.Price).ToNot(BeNil())
			Expect(*updated.Price).To(Equal(int64(200)))
		})
	})

	Describe("Integration scenarios", func() {
		It("should correctly handle packages with LPN changes", func() {
			// Crear un paquete con LPN
			originalPackage := DeliveryUnit{
				Lpn: "ORIGINAL-LPN",
				Items: []Item{
					{Sku: "SKU001"},
					{Sku: "SKU002"},
				},
			}

			// Obtener su DocID usando el LPN original
			originalDocID := originalPackage.DocID(ctx1)

			// Actualizar el LPN
			updatedPackage, changed := originalPackage.UpdateIfChanged(DeliveryUnit{
				Lpn: "NEW-LPN",
			})
			Expect(changed).To(BeTrue())
			Expect(updatedPackage.Lpn).To(Equal("NEW-LPN"))

			// Obtener el nuevo DocID - debería ser diferente porque usa el LPN nuevo
			updatedDocID := updatedPackage.DocID(ctx1)
			Expect(updatedDocID).ToNot(Equal(originalDocID))
			Expect(updatedDocID).To(Equal(HashByTenant(ctx1, "NEW-LPN")))
		})

		It("should correctly handle packages transitioning from LPN to no LPN", func() {
			// Crear un paquete con LPN
			originalPackage := DeliveryUnit{
				Lpn: "HAS-LPN",
				Items: []Item{
					{Sku: "SKU001"},
					{Sku: "SKU002"},
				},
			}

			// Obtener su DocID usando el LPN original
			originalDocID := originalPackage.DocID(ctx1)
			Expect(originalDocID).To(Equal(HashByTenant(ctx1, "HAS-LPN")))

			// Actualizar estableciendo LPN en vacío
			updatedPackage, changed := originalPackage.UpdateIfChanged(DeliveryUnit{
				Lpn: "",
			})

			// Comportamiento esperado: el LPN no se actualiza a vacío
			Expect(changed).To(BeFalse())
			Expect(updatedPackage.Lpn).To(Equal("HAS-LPN"))

			// El DocID debería seguir basándose en el LPN original
			updatedDocID := updatedPackage.DocID(ctx1)
			Expect(updatedDocID).To(Equal(originalDocID))
		})

		It("should correctly handle packages transitioning from no LPN to having LPN", func() {
			originalPackage := DeliveryUnit{
				Lpn:            "",
				noLPNReference: "REF001",
				Items: []Item{
					{Sku: "SKU001"},
					{Sku: "SKU002"},
				},
			}

			expectedInputs := []string{"REF001", "SKU001", "SKU002"} // ordenados
			originalDocID := originalPackage.DocID(ctx1)
			Expect(originalDocID).To(Equal(HashByTenant(ctx1, expectedInputs...)))

			updatedPackage, changed := originalPackage.UpdateIfChanged(DeliveryUnit{
				Lpn: "NEW-LPN",
			})
			Expect(changed).To(BeTrue())
			Expect(updatedPackage.Lpn).To(Equal("NEW-LPN"))

			updatedDocID := updatedPackage.DocID(ctx1)
			Expect(updatedDocID).ToNot(Equal(originalDocID))
			Expect(updatedDocID).To(Equal(HashByTenant(ctx1, "NEW-LPN")))
		})

	})

	It("should generate same ID regardless of item SKU order", func() {
		pkg1 := DeliveryUnit{
			Lpn:            "",
			noLPNReference: "REF001",
			Items: []Item{
				{Sku: "SKU001"},
				{Sku: "SKU002"},
			},
		}
		pkg2 := DeliveryUnit{
			Lpn:            "",
			noLPNReference: "REF001",
			Items: []Item{
				{Sku: "SKU002"},
				{Sku: "SKU001"}, // orden diferente
			},
		}

		id1 := pkg1.DocID(ctx1)
		id2 := pkg2.DocID(ctx1)

		Expect(id1).To(Equal(id2), "DocID debe ser igual si los SKUs son los mismos aunque estén en diferente orden")
	})

})
