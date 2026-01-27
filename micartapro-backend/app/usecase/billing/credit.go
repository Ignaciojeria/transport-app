package billing

import (
	"time"

	"github.com/google/uuid"
)

// UserCredits representa el saldo de créditos de un usuario
type UserCredits struct {
	UserID  uuid.UUID
	Balance int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreditTransaction representa una transacción de créditos
type CreditTransaction struct {
	ID            int64
	UserID        uuid.UUID
	Amount        int // Positivo = otorgar, Negativo = consumir
	TransactionType string // 'granted', 'consumed', 'expired', 'refunded'
	Source        string // 'payment.mercadopago', 'agent.usage', etc.
	SourceID      *string // ID del pago o evento que originó la transacción
	Description   *string
	BalanceBefore int
	BalanceAfter  int
	CreatedAt     time.Time
}

// GrantCreditsRequest representa una solicitud para otorgar créditos
type GrantCreditsRequest struct {
	UserID      uuid.UUID
	Amount      int
	Source      string
	SourceID    *string
	Description *string
}

// ConsumeCreditsRequest representa una solicitud para consumir créditos
type ConsumeCreditsRequest struct {
	UserID      uuid.UUID
	Amount      int
	Source      string
	SourceID    *string
	Description *string
}
