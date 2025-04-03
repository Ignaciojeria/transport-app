package domain_test

import (
	. "transport-app/app/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DocumentID", func() {

	Describe("IsZero", func() {
		It("should return true when DocumentID is empty", func() {
			id := DocumentID("")
			Expect(id.IsZero()).To(BeTrue())
		})

		It("should return false when DocumentID has value", func() {
			id := DocumentID("abc123")
			Expect(id.IsZero()).To(BeFalse())
		})
	})

	Describe("Equals", func() {
		It("should return true when values are equal", func() {
			id := DocumentID("abc123")
			Expect(id.Equals("abc123")).To(BeTrue())
		})

		It("should return false when values are different", func() {
			id := DocumentID("abc123")
			Expect(id.Equals("xyz789")).To(BeFalse())
		})
	})

	Describe("ShouldUpdate", func() {
		It("should return true when ID is not zero and different from existing", func() {
			id := DocumentID("new-value")
			Expect(id.ShouldUpdate("old-value")).To(BeTrue())
		})

		It("should return false when ID is zero", func() {
			id := DocumentID("")
			Expect(id.ShouldUpdate("anything")).To(BeFalse())
		})

		It("should return false when ID is the same as existing", func() {
			id := DocumentID("same-value")
			Expect(id.ShouldUpdate("same-value")).To(BeFalse())
		})
	})
})
