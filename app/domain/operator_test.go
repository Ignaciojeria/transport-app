package domain

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Operator", func() {
	var ctx context.Context

	BeforeEach(func() {
		ctx = buildCtx("org1", "CL")
	})

	Describe("DocID", func() {
		It("should generate DocID based on context and PrimaryEmail", func() {
			email := "operador@ejemplo.com"
			operator := Operator{
				Contact: Contact{PrimaryEmail: email},
			}

			expected := HashByTenant(ctx, email)
			Expect(operator.DocID(ctx)).To(Equal(expected))
		})

		It("should generate different DocIDs for different contexts", func() {
			email := "operador@ejemplo.com"
			operator := Operator{
				Contact: Contact{PrimaryEmail: email},
			}

			ctx1 := buildCtx("org1", "CL")
			ctx2 := buildCtx("org2", "AR")

			Expect(operator.DocID(ctx1)).ToNot(Equal(operator.DocID(ctx2)))
		})

		It("should generate different DocIDs for different emails", func() {
			operator1 := Operator{
				Contact: Contact{PrimaryEmail: "operador1@ejemplo.com"},
			}
			operator2 := Operator{
				Contact: Contact{PrimaryEmail: "operador2@ejemplo.com"},
			}

			Expect(operator1.DocID(ctx)).ToNot(Equal(operator2.DocID(ctx)))
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
