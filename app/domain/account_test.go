package domain

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account", func() {
	Describe("UUID", func() {
		It("should generate deterministic UUID for same email", func() {
			account := Account{
				Email: "usuario@ejemplo.com",
				Role:  "admin",
			}

			uuid1 := account.UUID()
			uuid2 := account.UUID()

			Expect(uuid1).ToNot(Equal(uuid.Nil))
			Expect(uuid1).To(Equal(uuid2))
		})

		It("should generate different UUIDs for different emails", func() {
			account1 := Account{Email: "usuario1@ejemplo.com", Role: "admin"}
			account2 := Account{Email: "usuario2@ejemplo.com", Role: "admin"}

			uuid1 := account1.UUID()
			uuid2 := account2.UUID()

			Expect(uuid1).ToNot(Equal(uuid2))
		})

		It("should generate same UUID for same email regardless of role", func() {
			account1 := Account{Email: "usuario@ejemplo.com", Role: "admin"}
			account2 := Account{Email: "usuario@ejemplo.com", Role: "user"}

			uuid1 := account1.UUID()
			uuid2 := account2.UUID()

			Expect(uuid1).To(Equal(uuid2))
		})

		It("should generate valid UUID v5", func() {
			account := Account{
				Email: "usuario@ejemplo.com",
				Role:  "admin",
			}

			generatedUUID := account.UUID()

			Expect(generatedUUID.Version()).To(Equal(uuid.Version(5)))
			Expect(generatedUUID.Variant()).To(Equal(uuid.RFC4122))
		})

		It("should be case sensitive for email", func() {
			account1 := Account{Email: "Usuario@Ejemplo.com", Role: "admin"}
			account2 := Account{Email: "usuario@ejemplo.com", Role: "admin"}

			uuid1 := account1.UUID()
			uuid2 := account2.UUID()

			Expect(uuid1).ToNot(Equal(uuid2))
		})

		It("should handle special characters in email", func() {
			account := Account{
				Email: "usuario+test@ejemplo.com",
				Role:  "admin",
			}

			generatedUUID := account.UUID()

			Expect(generatedUUID).ToNot(Equal(uuid.Nil))
			Expect(generatedUUID.Version()).To(Equal(uuid.Version(5)))
		})

		It("should handle empty email", func() {
			account := Account{
				Email: "",
				Role:  "admin",
			}

			generatedUUID := account.UUID()

			Expect(generatedUUID).ToNot(Equal(uuid.Nil))
			Expect(generatedUUID.Version()).To(Equal(uuid.Version(5)))
		})
	})
})
