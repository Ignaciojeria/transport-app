package domain

import (
	"github.com/biter777/countries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Operator", func() {

	Describe("DocID", func() {
		It("should generate DocID based on Organization and PrimaryEmail", func() {
			org := Organization{ID: 1, Country: countries.CL}
			email := "operador@ejemplo.com"
			operator := Operator{
				Organization: org,
				Contact:      Contact{PrimaryEmail: email},
			}

			expected := Hash(org, email)
			Expect(operator.DocID()).To(Equal(expected))
		})
	})

	Describe("UpdateIfChanged", func() {
		It("should update Role if the new one is not empty", func() {
			original := Operator{
				Role: RolePlanner,
			}
			newOperator := Operator{
				Role: RoleDispatcher,
			}

			updated := original.UpdateIfChanged(newOperator)
			Expect(updated.Role).To(Equal(RoleDispatcher))
		})

		It("should not update Role if the new one is empty", func() {
			original := Operator{
				Role: RolePlanner,
			}
			newOperator := Operator{
				Role: "",
			}

			updated := original.UpdateIfChanged(newOperator)
			Expect(updated.Role).To(Equal(RolePlanner))
		})
	})
})

var _ = Describe("Role", func() {

	DescribeTable("IsValid",
		func(role Role, expected bool) {
			Expect(role.IsValid()).To(Equal(expected))
		},
		Entry("valid admin", RoleAdmin, true),
		Entry("valid driver", RoleDriver, true),
		Entry("valid planner", RolePlanner, true),
		Entry("valid dispatcher", RoleDispatcher, true),
		Entry("valid seller", RoleSeller, true),
		Entry("invalid role", Role("unknown"), false),
		Entry("empty role", Role(""), false),
	)

	It("should return the string value of the role", func() {
		Expect(RoleAdmin.String()).To(Equal("admin"))
		Expect(Role("").String()).To(Equal(""))
	})
})
