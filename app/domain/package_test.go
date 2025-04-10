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
		It("should generate unique ID based on context and Lpn", func() {
			package1 := Package{
				Lpn: "PKG001",
			}
			package2 := Package{
				Lpn: "PKG002",
			}

			Expect(package1.DocID(ctx1)).To(Equal(Hash(ctx1, "PKG001")))
			Expect(package1.DocID(ctx1)).ToNot(Equal(package2.DocID(ctx1)))
			Expect(package1.DocID(ctx1)).ToNot(Equal(package1.DocID(ctx2)))
		})

		It("should return hashed empty string if Lpn is empty", func() {
			pkg := Package{
				Lpn: "",
			}
			Expect(pkg.DocID(ctx1)).To(Equal(Hash(ctx1, "")))
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
				ItemReferences: []ItemReference{
					{Sku: "ITEM001", Quantity: Quantity{QuantityNumber: 2, QuantityUnit: "unit"}},
					{Sku: "ITEM002", Quantity: Quantity{QuantityNumber: 1, QuantityUnit: "unit"}},
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
			Expect(updated.ItemReferences).To(Equal(basePackage.ItemReferences))
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
			Expect(updated.ItemReferences).To(Equal(basePackage.ItemReferences))
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
			Expect(updated.ItemReferences).To(Equal(basePackage.ItemReferences))
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
			Expect(updated.ItemReferences).To(Equal(basePackage.ItemReferences))
		})

		It("should update ItemReferences", func() {
			newPackage := Package{
				ItemReferences: []ItemReference{
					{Sku: "ITEM003", Quantity: Quantity{QuantityNumber: 3, QuantityUnit: "unit"}},
					{Sku: "ITEM004", Quantity: Quantity{QuantityNumber: 4, QuantityUnit: "box"}},
				},
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			Expect(changed).To(BeTrue())
			Expect(updated.ItemReferences).To(Equal(newPackage.ItemReferences))
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
				ItemReferences: []ItemReference{
					{Sku: "ITEM005", Quantity: Quantity{QuantityNumber: 5, QuantityUnit: "pallet"}},
				},
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			Expect(changed).To(BeTrue())
			Expect(updated.Lpn).To(Equal("PKG-MULTI-UPDATE"))
			Expect(updated.Weight).To(Equal(newPackage.Weight))
			Expect(updated.ItemReferences).To(Equal(newPackage.ItemReferences))
			// Estos campos no deberían cambiar
			Expect(updated.Dimensions).To(Equal(basePackage.Dimensions))
			Expect(updated.Insurance).To(Equal(basePackage.Insurance))
		})

		It("should handle empty ItemReferences array", func() {
			// Primero confirmar que tenemos referencias iniciales
			Expect(basePackage.ItemReferences).ToNot(BeEmpty())

			// Intentar actualizar con un array vacío
			newPackage := Package{
				ItemReferences: []ItemReference{},
			}

			updated, changed := basePackage.UpdateIfChanged(newPackage)

			// Las referencias deberían mantenerse sin cambios, ya que el array vacío
			// no debería sobrescribir los valores existentes según la lógica del método
			Expect(changed).To(BeFalse())
			Expect(updated.ItemReferences).To(Equal(basePackage.ItemReferences))
		})
	})
})
