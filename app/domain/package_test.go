package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Package", func() {
	var org1 = Organization{ID: 1, Country: countries.CL}
	var org2 = Organization{ID: 2, Country: countries.AR}

	Describe("DocID", func() {
		It("should generate unique ID based on Organization and Lpn", func() {
			package1 := Package{
				Organization: org1,
				Lpn:          "PKG001",
			}
			package2 := Package{
				Organization: org1,
				Lpn:          "PKG002",
			}
			package3 := Package{
				Organization: org2,
				Lpn:          "PKG001",
			}

			Expect(package1.DocID()).To(Equal(Hash(org1, "PKG001")))
			Expect(package1.DocID()).ToNot(Equal(package2.DocID()))
			Expect(package1.DocID()).ToNot(Equal(package3.DocID()))
		})

		It("should return empty DocumentID if Lpn is empty", func() {
			pkg := Package{
				Organization: org1,
				Lpn:          "",
			}
			Expect(pkg.DocID()).To(Equal(Hash(org1, "")))
		})
	})

	Describe("SearchPackageByLpn", func() {
		var packages []Package

		BeforeEach(func() {
			packages = []Package{
				{Lpn: "PKG001", Organization: org1},
				{Lpn: "PKG002", Organization: org1},
				{Lpn: "PKG003", Organization: org2},
			}
		})

		It("should find package by lpn when it exists", func() {
			result := SearchPackageByLpn(packages, "PKG002")
			Expect(result.Lpn).To(Equal("PKG002"))
			Expect(result.Organization).To(Equal(org1))
		})

		It("should return empty package when lpn doesn't exist", func() {
			result := SearchPackageByLpn(packages, "NONEXISTENT")
			Expect(result).To(Equal(Package{}))
		})

		It("should return first matching package when multiple matches exist", func() {
			// Agregar un duplicado con el mismo Lpn
			duplicatePackages := append(packages, Package{Lpn: "PKG001", Organization: org2})

			result := SearchPackageByLpn(duplicatePackages, "PKG001")
			Expect(result.Lpn).To(Equal("PKG001"))
			Expect(result.Organization).To(Equal(org1)) // Debería devolver la primera coincidencia
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
				Lpn:          "PKG-TEST",
				Organization: org1,
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

			updated := basePackage.UpdateIfChanged(newPackage)

			Expect(updated.Lpn).To(Equal("PKG-UPDATED"))
			// Verificar que otros campos se mantienen igual
			Expect(updated.Organization).To(Equal(basePackage.Organization))
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

			updated := basePackage.UpdateIfChanged(newPackage)

			Expect(updated.Dimensions).To(Equal(newPackage.Dimensions))
			// Verificar que otros campos se mantienen igual
			Expect(updated.Lpn).To(Equal(basePackage.Lpn))
			Expect(updated.Organization).To(Equal(basePackage.Organization))
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

			updated := basePackage.UpdateIfChanged(newPackage)

			Expect(updated.Weight).To(Equal(newPackage.Weight))
			// Verificar que otros campos se mantienen igual
			Expect(updated.Lpn).To(Equal(basePackage.Lpn))
			Expect(updated.Organization).To(Equal(basePackage.Organization))
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

			updated := basePackage.UpdateIfChanged(newPackage)

			Expect(updated.Insurance).To(Equal(newPackage.Insurance))
			// Verificar que otros campos se mantienen igual
			Expect(updated.Lpn).To(Equal(basePackage.Lpn))
			Expect(updated.Organization).To(Equal(basePackage.Organization))
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

			updated := basePackage.UpdateIfChanged(newPackage)

			Expect(updated.ItemReferences).To(Equal(newPackage.ItemReferences))
			// Verificar que otros campos se mantienen igual
			Expect(updated.Lpn).To(Equal(basePackage.Lpn))
			Expect(updated.Organization).To(Equal(basePackage.Organization))
			Expect(updated.Dimensions).To(Equal(basePackage.Dimensions))
			Expect(updated.Weight).To(Equal(basePackage.Weight))
			Expect(updated.Insurance).To(Equal(basePackage.Insurance))
		})

		It("should not update fields when new UnitValues are empty", func() {
			newPackage := Package{
				Lpn: "",
				// Todos los demás campos están vacíos o con sus valores por defecto
			}

			updated := basePackage.UpdateIfChanged(newPackage)

			// Nada debería cambiar
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

			updated := basePackage.UpdateIfChanged(newPackage)

			Expect(updated.Lpn).To(Equal("PKG-MULTI-UPDATE"))
			Expect(updated.Weight).To(Equal(newPackage.Weight))
			Expect(updated.ItemReferences).To(Equal(newPackage.ItemReferences))
			// Estos campos no deberían cambiar
			Expect(updated.Organization).To(Equal(basePackage.Organization))
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

			updated := basePackage.UpdateIfChanged(newPackage)

			// Las referencias deberían mantenerse sin cambios, ya que el array vacío
			// no debería sobrescribir los valores existentes según la lógica del método
			Expect(updated.ItemReferences).To(Equal(basePackage.ItemReferences))
		})
	})
})
