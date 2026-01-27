package supabaserepo

import (
	"context"
	"errors"

	"micartapro/app/usecase/billing"

	"github.com/google/uuid"
)

var (
	ErrInsufficientCredits = errors.New("insufficient credits")
	ErrCreditsNotFound     = errors.New("credits not found")
)

type GetUserCredits func(ctx context.Context, userID uuid.UUID) (*billing.UserCredits, error)
type GrantCredits func(ctx context.Context, req billing.GrantCreditsRequest) (*billing.CreditTransaction, error)
type ConsumeCredits func(ctx context.Context, req billing.ConsumeCreditsRequest) (*billing.CreditTransaction, error)
type GetCreditTransactions func(ctx context.Context, userID uuid.UUID, limit int) ([]billing.CreditTransaction, error)
