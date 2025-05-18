package tidbrepository

import (
	"context"
	"transport-app/app/adapter/out/tidbrepository/table"
	"transport-app/app/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FindAccountByEmail", func() {

	It("should return empty operator when account not found", func() {
		ctx := context.Background()
		findAccount := NewFindAccountByEmail(connection)

		account, err := findAccount(ctx, "nonexistent@example.com")
		Expect(err).ToNot(HaveOccurred())
		Expect(account).To(Equal(domain.Account{}))
	})

	It("should return operator when account exists", func() {
		ctx := context.Background()

		// Create test account
		account := domain.Account{
			Email: "test@example.com",
		}
		err := NewUpsertAccount(connection)(ctx, account)
		Expect(err).ToNot(HaveOccurred())

		findAccount := NewFindAccountByEmail(connection)
		accountDomain, err := findAccount(ctx, "test@example.com")
		Expect(err).ToNot(HaveOccurred())
		Expect(accountDomain.Email).To(Equal("test@example.com"))

		// Verify using document_id
		var dbAccount table.Account
		err = connection.DB.WithContext(ctx).
			Table("accounts").
			Where("document_id = ?", accountDomain.DocID()).
			First(&dbAccount).Error
		Expect(err).ToNot(HaveOccurred())
		Expect(dbAccount.Email).To(Equal("test@example.com"))
	})

	It("should fail if database has no accounts table", func() {
		ctx := context.Background()
		findAccount := NewFindAccountByEmail(noTablesContainerConnection)

		_, err := findAccount(ctx, "test@example.com")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("accounts"))
	})
})
