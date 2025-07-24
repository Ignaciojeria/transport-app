package tidbrepository

import (
	"context"
	"transport-app/app/domain"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FindDefaultTenantByEmail", func() {

	It("should return empty tenant account when not found", func() {
		ctx := context.Background()
		findDefaultTenant := NewFindDefaultTenantByEmail(connection)

		tenantAccount, err := findDefaultTenant(ctx, "nonexistent@example.com")
		Expect(err).ToNot(HaveOccurred())
		Expect(tenantAccount).To(Equal(domain.TenantAccount{}))
	})

	It("should return tenant account when exists and not invited", func() {
		ctx := context.Background()

		// Create test tenant
		tenant, tenantCtx, err := CreateTestTenant(ctx, connection)
		Expect(err).ToNot(HaveOccurred())

		// Create test account
		account := domain.Account{
			Email: "test@example.com",
		}
		err = NewUpsertAccount(connection, nil)(ctx, account)
		Expect(err).ToNot(HaveOccurred())

		// Create tenant account link (not invited)
		tenantAccount := domain.TenantAccount{
			Tenant:  tenant,
			Account: account,
			Role:    "admin",
			Status:  "active",
			Invited: false,
		}
		err = NewSaveTenantAccount(connection, nil)(tenantCtx, tenantAccount)
		Expect(err).ToNot(HaveOccurred())

		// Test find function
		findDefaultTenant := NewFindDefaultTenantByEmail(connection)
		result, err := findDefaultTenant(ctx, "test@example.com")
		Expect(err).ToNot(HaveOccurred())
		Expect(result.Account.Email).To(Equal("test@example.com"))
		Expect(result.Tenant.ID).To(Equal(tenant.ID))
		Expect(result.Role).To(Equal("admin"))
		Expect(result.Status).To(Equal("active"))
		Expect(result.Invited).To(BeFalse())
	})

	It("should not return tenant account when invited", func() {
		ctx := context.Background()

		// Create test tenant
		tenant, tenantCtx, err := CreateTestTenant(ctx, connection)
		Expect(err).ToNot(HaveOccurred())

		// Create test account
		account := domain.Account{
			Email: "invited@example.com",
		}
		err = NewUpsertAccount(connection, nil)(ctx, account)
		Expect(err).ToNot(HaveOccurred())

		// Create tenant account link (invited)
		tenantAccount := domain.TenantAccount{
			Tenant:  tenant,
			Account: account,
			Role:    "user",
			Status:  "pending",
			Invited: true,
		}
		err = NewSaveTenantAccount(connection, nil)(tenantCtx, tenantAccount)
		Expect(err).ToNot(HaveOccurred())

		// Test find function - should not find invited account
		findDefaultTenant := NewFindDefaultTenantByEmail(connection)
		result, err := findDefaultTenant(ctx, "invited@example.com")
		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal(domain.TenantAccount{}))
	})

	It("should fail if database has no required tables", func() {
		ctx := context.Background()
		findDefaultTenant := NewFindDefaultTenantByEmail(noTablesContainerConnection)

		_, err := findDefaultTenant(ctx, "test@example.com")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("account_tenants"))
	})
})
