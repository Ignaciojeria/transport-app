package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FindAccountByEmail", func() {
	BeforeEach(func() {
		// Clean the accounts table before each test
		err := connection.DB.Exec("DELETE FROM accounts").Error
		Expect(err).ToNot(HaveOccurred())
	})

	It("should return empty operator when account not found", func() {
		ctx := context.Background()
		findAccount := NewFindAccountByEmail(connection)

		operator, err := findAccount(ctx, "nonexistent@example.com")
		Expect(err).ToNot(HaveOccurred())
		Expect(operator).To(Equal(domain.Operator{}))
	})

	It("should return operator when account exists", func() {
		ctx := context.Background()

		// Create test account
		account := table.Account{
			Email:    "test@example.com",
			IsActive: true,
		}
		err := connection.DB.Create(&account).Error
		Expect(err).ToNot(HaveOccurred())

		findAccount := NewFindAccountByEmail(connection)
		operator, err := findAccount(ctx, "test@example.com")
		Expect(err).ToNot(HaveOccurred())
		Expect(operator.Contact.PrimaryEmail).To(Equal("test@example.com"))
	})

	It("should fail if database has no accounts table", func() {
		ctx := context.Background()
		findAccount := NewFindAccountByEmail(noTablesContainerConnection)

		_, err := findAccount(ctx, "test@example.com")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("accounts"))
	})
})
