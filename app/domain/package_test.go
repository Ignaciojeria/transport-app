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
			package1 := Package{
				Lpn: "PKG001",
			}
			package2 := Package{
				Lpn: "PKG002",
			}

			Expect(package1.DocID(ctx1, "REF001")).To(Equal(HashByTenant(ctx1, "PKG001")))
			Expect(package1.DocID(ctx1, "REF001")).ToNot(Equal(package2.DocID(ctx1, "REF001")))
			Expect(package1.DocID(ctx1, "REF001")).ToNot(Equal(package1.DocID(ctx2, "REF001")))
		})

		It("should generate ID based on reference, index and sorted items when Lpn is empty", func() {
			pkg := Package{
				Lpn:   "",
				Index: 1,
				Items: []Item{
					{Sku: "SKU002"},
					{Sku: "SKU001"},
				},
			}
			expectedInputs := []string{"REF001", "index:1", "SKU001", "SKU002"} // ordenados
			Expect(pkg.DocID(ctx1, "REF001")).To(Equal(HashByTenant(ctx1, expectedInputs...)))
		})

		It("should generate different IDs for packages with same items but different references", func() {
			pkg := Package{
				Lpn: "",
				Items: []Item{
					{Sku: "SKU001"},
					{Sku: "SKU002"},
				},
			}

			id1 := pkg.DocID(ctx1, "REF001")
			id2 := pkg.DocID(ctx1, "REF002")

			Expect(id1).ToNot(Equal(id2))
		})

		It("should generate different IDs for packages with same reference but different items", func() {
			pkg1 := Package{
				Lpn: "",
				Items: []Item{
					{Sku: "SKU001"},
					{Sku: "SKU002"},
				},
			}

			pkg2 := Package{
				Lpn: "",
				Items: []Item{
					{Sku: "SKU001"},
					{Sku: "SKU003"}, // Diferente SKU
				},
			}

			id1 := pkg1.DocID(ctx1, "REF001")
			id2 := pkg2.DocID(ctx1, "REF001")

			Expect(id1).ToNot(Equal(id2))
		})

		It("should generate consistent IDs for same reference and items", func() {
			pkg1 := Package{
				Lpn: "",
				Items: []Item{
					{Sku: "SKU001"},
					{Sku: "SKU002"},
				},
			}

			pkg2 := Package{
				Lpn: "",
				Items: []Item{
					{Sku: "SKU001"},
					{Sku: "SKU002"},
				},
			}

			id1 := pkg1.DocID(ctx1, "REF001")
			id2 := pkg2.DocID(ctx1, "REF001")

			Expect(id1).To(Equal(id2))
		})

		It("should handle the case with no LPN and no items", func() {
			pkg := Package{
				Lpn:   "",
				Items: []Item{},
				Index: 1,
			}

			expected := HashByTenant(ctx1, "REF001", "index:1")
			Expect(pkg.DocID(ctx1, "REF001")).To(Equal(expected))
		})

	})

	Describe("SearchPackageByLpn", func() {
		var packages []Package

		BeforeEach(func() {
			packages = []Package{
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
			Expect(result).To(Equal(Package{}))
		})

		It("should return first matching package when multiple matches exist", func() {
			// Agregar un duplicado con el mismo Lpn
			duplicatePackages := append(packages, Package{Lpn: "PKG001"})

			result := SearchPackageByLpn(duplicatePackages, "PKG001")
			Expect(result.Lpn).To(Equal("PKG001"))
		})

		It("should handle empty package slice", func() {
			result := SearchPackageByLpn([]Package{}, "PKG001")
			Expect(result).To(Equal(Package{}))
		})

		It("should handle case-sensitive lpn comparison", func() {
			// Agregar paquete con LPN en minúsculas
			packagesWithCase := append(packages, Package{Lpn: "pkg001"})

			// Debería devolver el original (PKG001) y no el nuevo (pkg001)
			result := SearchPackageByLpn(packagesWithCase, "PKG001")
			Expect(result.Lpn).To(Equal("PKG001"))

			// Debería encontrar el lpn en minúsculas
			result = SearchPackageByLpn(packagesWithCase, "pkg001")
			Expect(result.Lpn).To(Equal("pkg001"))
		})
	})

	Describe("UpdateIfChanged", func() {
		var basePackage Package

		BeforeEach(func() {
			basePackage = Package{
				Lpn: "PKG-TEST",
				Dimensions: Dimensions{
					Length: 10.0,
					Width:  20.0,
					Height: 30.0,
					Unit:   "cm",
				},
				Weight: Weight{
					Value: 5.0,
					Unit:  "kg",
				},
				Insurance: Insurance{
					UnitValue: 1000.0,
					Currency:  "USD",
				},
				Items: []Item{
					{
						Sku:         "ITEM001",
						Description: "Item 1 Description",
						Quantity:    Quantity{QuantityNumber: 2, QuantityUnit: "unit"},
						Weight: Weight{
							Value: 1.0,
							Unit:  "kg",
						},
					},
					{
						Sku:         "ITEM002",
						Description: "Item 2 Description",
						Quantity:    Quantity{QuantityNumber: 1, QuantityUnit: "unit"},
						Weight: Weight{
							Value: 0.5,
							Unit:  "kg",
						},
					},
				},
			}
		})

		It("should update Lpn", func() {
			newPackage := Package{
				Lpn: "PKG-UPDATED",
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			Expect(changed).To(BeTrue())
			Expect(updated.Lpn).To(Equal("PKG-UPDATED"))
			// Verificar que otros campos se mantienen igual
			Expect(updated.Dimensions).To(Equal(basePackage.Dimensions))
			Expect(updated.Weight).To(Equal(basePackage.Weight))
			Expect(updated.Insurance).To(Equal(basePackage.Insurance))
			Expect(updated.Items).To(Equal(basePackage.Items))
		})

		It("should update Dimensions", func() {
			newPackage := Package{
				Dimensions: Dimensions{
					Length: 15.0,
					Width:  25.0,
					Height: 35.0,
					Unit:   "mm",
				},
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			Expect(changed).To(BeTrue())
			Expect(updated.Dimensions).To(Equal(newPackage.Dimensions))
			// Verificar que otros campos se mantienen igual
			Expect(updated.Lpn).To(Equal(basePackage.Lpn))
			Expect(updated.Weight).To(Equal(basePackage.Weight))
			Expect(updated.Insurance).To(Equal(basePackage.Insurance))
			Expect(updated.Items).To(Equal(basePackage.Items))
		})

		It("should update Weight", func() {
			newPackage := Package{
				Weight: Weight{
					Value: 7.5,
					Unit:  "lb",
				},
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			Expect(changed).To(BeTrue())
			Expect(updated.Weight).To(Equal(newPackage.Weight))
			// Verificar que otros campos se mantienen igual
			Expect(updated.Lpn).To(Equal(basePackage.Lpn))
			Expect(updated.Dimensions).To(Equal(basePackage.Dimensions))
			Expect(updated.Insurance).To(Equal(basePackage.Insurance))
			Expect(updated.Items).To(Equal(basePackage.Items))
		})

		It("should update Insurance", func() {
			newPackage := Package{
				Insurance: Insurance{
					UnitValue: 2000.0,
					Currency:  "EUR",
				},
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			Expect(changed).To(BeTrue())
			Expect(updated.Insurance).To(Equal(newPackage.Insurance))
			// Verificar que otros campos se mantienen igual
			Expect(updated.Lpn).To(Equal(basePackage.Lpn))
			Expect(updated.Dimensions).To(Equal(basePackage.Dimensions))
			Expect(updated.Weight).To(Equal(basePackage.Weight))
			Expect(updated.Items).To(Equal(basePackage.Items))
		})

		It("should update Items", func() {
			newPackage := Package{
				Items: []Item{
					{
						Sku:         "ITEM003",
						Description: "New Item 3",
						Quantity:    Quantity{QuantityNumber: 3, QuantityUnit: "unit"},
						Weight: Weight{
							Value: 1.5,
							Unit:  "kg",
						},
					},
					{
						Sku:         "ITEM004",
						Description: "New Item 4",
						Quantity:    Quantity{QuantityNumber: 4, QuantityUnit: "box"},
						Weight: Weight{
							Value: 2.0,
							Unit:  "kg",
						},
					},
				},
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			Expect(changed).To(BeTrue())
			Expect(updated.Items).To(Equal(newPackage.Items))
			// Verificar que otros campos se mantienen igual
			Expect(updated.Lpn).To(Equal(basePackage.Lpn))
			Expect(updated.Dimensions).To(Equal(basePackage.Dimensions))
			Expect(updated.Weight).To(Equal(basePackage.Weight))
			Expect(updated.Insurance).To(Equal(basePackage.Insurance))
		})

		It("should not update fields when new values are empty", func() {
			newPackage := Package{
				Lpn: "",
				// Todos los demás campos están vacíos o con sus valores por defecto
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			// Nada debería cambiar
			Expect(changed).To(BeFalse())
			Expect(updated).To(Equal(basePackage))
		})

		It("should update multiple fields at once", func() {
			newPackage := Package{
				Lpn: "PKG-MULTI-UPDATE",
				Weight: Weight{
					Value: 8.0,
					Unit:  "oz",
				},
				Items: []Item{
					{
						Sku:         "ITEM005",
						Description: "New Item for Multi-update",
						Quantity:    Quantity{QuantityNumber: 5, QuantityUnit: "pallet"},
						Weight: Weight{
							Value: 3.0,
							Unit:  "kg",
						},
					},
				},
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			Expect(changed).To(BeTrue())
			Expect(updated.Lpn).To(Equal("PKG-MULTI-UPDATE"))
			Expect(updated.Weight).To(Equal(newPackage.Weight))
			Expect(updated.Items).To(Equal(newPackage.Items))
			// Estos campos no deberían cambiar
			Expect(updated.Dimensions).To(Equal(basePackage.Dimensions))
			Expect(updated.Insurance).To(Equal(basePackage.Insurance))
		})

		It("should handle empty Items array", func() {
			// Primero confirmar que tenemos items iniciales
			Expect(basePackage.Items).ToNot(BeEmpty())

			// Intentar actualizar con un array vacío
			newPackage := Package{
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
			newPackage := Package{
				Items: []Item{
					{
						Sku:         "NEW-SINGLE-ITEM",
						Description: "Single Item Update",
						Quantity:    Quantity{QuantityNumber: 10, QuantityUnit: "unit"},
					},
				},
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			// Debería reemplazar todos los items con el nuevo array
			Expect(changed).To(BeTrue())
			Expect(len(updated.Items)).To(Equal(1))
			Expect(updated.Items[0].Sku).To(Equal("NEW-SINGLE-ITEM"))
		})
	})

	Describe("Integration scenarios", func() {
		It("should correctly handle packages with LPN changes", func() {
			// Crear un paquete con LPN
			originalPackage := Package{
				Lpn: "ORIGINAL-LPN",
				Items: []Item{
					{Sku: "SKU001"},
					{Sku: "SKU002"},
				},
			}

			// Obtener su DocID usando el LPN original
			originalDocID := originalPackage.DocID(ctx1, "REF001")

			// Actualizar el LPN
			updatedPackage, changed := originalPackage.UpdateIfChanged(Package{
				Lpn: "NEW-LPN",
			})
			Expect(changed).To(BeTrue())
			Expect(updatedPackage.Lpn).To(Equal("NEW-LPN"))

			// Obtener el nuevo DocID - debería ser diferente porque usa el LPN nuevo
			updatedDocID := updatedPackage.DocID(ctx1, "REF001")
			Expect(updatedDocID).ToNot(Equal(originalDocID))
			Expect(updatedDocID).To(Equal(HashByTenant(ctx1, "NEW-LPN")))
		})

		It("should correctly handle packages transitioning from LPN to no LPN", func() {
			// Crear un paquete con LPN
			originalPackage := Package{
				Lpn: "HAS-LPN",
				Items: []Item{
					{Sku: "SKU001"},
					{Sku: "SKU002"},
				},
			}

			// Obtener su DocID usando el LPN original
			originalDocID := originalPackage.DocID(ctx1, "REF001")
			Expect(originalDocID).To(Equal(HashByTenant(ctx1, "HAS-LPN")))

			// Actualizar estableciendo LPN en vacío
			updatedPackage, changed := originalPackage.UpdateIfChanged(Package{
				Lpn: "",
			})

			// Comportamiento esperado: el LPN no se actualiza a vacío
			Expect(changed).To(BeFalse())
			Expect(updatedPackage.Lpn).To(Equal("HAS-LPN"))

			// El DocID debería seguir basándose en el LPN original
			updatedDocID := updatedPackage.DocID(ctx1, "REF001")
			Expect(updatedDocID).To(Equal(originalDocID))
		})

		It("should correctly handle packages transitioning from no LPN to having LPN", func() {
			originalPackage := Package{
				Lpn:   "",
				Index: 1,
				Items: []Item{
					{Sku: "SKU001"},
					{Sku: "SKU002"},
				},
			}

			expectedInputs := []string{"REF001", "index:1", "SKU001", "SKU002"}
			originalDocID := originalPackage.DocID(ctx1, "REF001")
			Expect(originalDocID).To(Equal(HashByTenant(ctx1, expectedInputs...)))

			updatedPackage, changed := originalPackage.UpdateIfChanged(Package{
				Lpn: "NEW-LPN",
			})
			Expect(changed).To(BeTrue())
			Expect(updatedPackage.Lpn).To(Equal("NEW-LPN"))

			updatedDocID := updatedPackage.DocID(ctx1, "REF001")
			Expect(updatedDocID).ToNot(Equal(originalDocID))
			Expect(updatedDocID).To(Equal(HashByTenant(ctx1, "NEW-LPN")))
		})

	})

	It("should generate same ID regardless of item SKU order", func() {
		pkg1 := Package{
			Lpn:   "",
			Index: 1,
			Items: []Item{
				{Sku: "SKU001"},
				{Sku: "SKU002"},
			},
		}
		pkg2 := Package{
			Lpn:   "",
			Index: 1,
			Items: []Item{
				{Sku: "SKU002"},
				{Sku: "SKU001"}, // orden diferente
			},
		}

		id1 := pkg1.DocID(ctx1, "REF001")
		id2 := pkg2.DocID(ctx1, "REF001")

		Expect(id1).To(Equal(id2), "DocID debe ser igual si los SKUs son los mismos aunque estén en diferente orden")
	})

})
