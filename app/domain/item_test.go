package domain

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CompareItemReferences", func() {
	It("should return true for identical item references", func() {
		oldRefs := []ItemReference{
			{Sku: "ITEM001", Quantity: 2},
			{Sku: "ITEM002", Quantity: 1},
		}
		newRefs := []ItemReference{
			{Sku: "ITEM001", Quantity: 2},
			{Sku: "ITEM002", Quantity: 1},
		}
		Expect(compareItemReferences(oldRefs, newRefs)).To(BeTrue())
	})

	It("should return false for different lengths", func() {
		oldRefs := []ItemReference{
			{Sku: "ITEM001", Quantity: 2},
		}
		newRefs := []ItemReference{}
		Expect(compareItemReferences(oldRefs, newRefs)).To(BeFalse())
	})

	It("should return false for different Sku", func() {
		oldRefs := []ItemReference{
			{Sku: "ITEM001", Quantity: 2},
		}
		newRefs := []ItemReference{
			{Sku: "ITEM999", Quantity: 2},
		}
		Expect(compareItemReferences(oldRefs, newRefs)).To(BeFalse())
	})

	It("should return false for different Quantity", func() {
		oldRefs := []ItemReference{
			{Sku: "ITEM001", Quantity: 2},
		}
		newRefs := []ItemReference{
			{Sku: "ITEM001", Quantity: 5},
		}
		Expect(compareItemReferences(oldRefs, newRefs)).To(BeFalse())
	})

	It("should return true for both empty slices", func() {
		Expect(compareItemReferences([]ItemReference{}, []ItemReference{})).To(BeTrue())
	})

	It("should return false again for different lengths", func() { // test duplicado pero lo mantenemos renombrado
		oldRefs := []ItemReference{
			{Sku: "ITEM001", Quantity: 2},
		}
		newRefs := []ItemReference{}
		Expect(compareItemReferences(oldRefs, newRefs)).To(BeFalse())
	})
})
