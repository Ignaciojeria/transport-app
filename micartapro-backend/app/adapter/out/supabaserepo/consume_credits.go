package supabaserepo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"micartapro/app/shared/configuration"
	"micartapro/app/usecase/billing"

	ioc "github.com/Ignaciojeria/ioc"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
)

func init() {
	ioc.Register(NewConsumeCredits)
}

func NewConsumeCredits(supabase *supabase.Client, conf configuration.Conf) ConsumeCredits {
	return func(ctx context.Context, req billing.ConsumeCreditsRequest) (*billing.CreditTransaction, error) {
		// Construir los parámetros para la función
		params := map[string]interface{}{
			"p_user_id": req.UserID.String(),
			"p_amount":  req.Amount,
			"p_source":  req.Source,
		}

		if req.SourceID != nil {
			params["p_source_id"] = *req.SourceID
		}
		if req.Description != nil {
			params["p_description"] = *req.Description
		}

		// Llamar a la función consume_credits usando HTTP POST directo
		rpcURL := fmt.Sprintf("%s/rest/v1/rpc/consume_credits", conf.SUPABASE_PROJECT_URL)
		requestBody, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("error marshaling RPC params: %w", err)
		}

		httpReq, err := http.NewRequestWithContext(ctx, "POST", rpcURL, bytes.NewBuffer(requestBody))
		if err != nil {
			return nil, fmt.Errorf("error creating HTTP request: %w", err)
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("apikey", conf.SUPABASE_BACKEND_API_KEY)
		httpReq.Header.Set("Authorization", "Bearer "+conf.SUPABASE_BACKEND_API_KEY)
		httpReq.Header.Set("Prefer", "return=representation")

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			return nil, fmt.Errorf("error calling consume_credits RPC: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading RPC response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("RPC call failed with status %d: %s", resp.StatusCode, string(body))
		}

		// La función retorna un array con un solo elemento
		var results []struct {
			Success       bool   `json:"success"`
			BalanceBefore int64  `json:"balance_before"`
			BalanceAfter  int64  `json:"balance_after"`
			TransactionID *int64 `json:"transaction_id"`
		}

		if err := json.Unmarshal(body, &results); err != nil {
			return nil, fmt.Errorf("error unmarshaling consume_credits result: %w", err)
		}

		if len(results) == 0 {
			return nil, fmt.Errorf("consume_credits function returned no results")
		}

		result := results[0]

		if !result.Success {
			return nil, ErrInsufficientCredits
		}

		if result.TransactionID == nil {
			return nil, fmt.Errorf("transaction_id is null")
		}

		// Obtener la transacción completa
		var transaction []struct {
			ID              int64     `json:"id"`
			UserID          string    `json:"user_id"`
			Amount          int       `json:"amount"`
			TransactionType string    `json:"transaction_type"`
			Source          string    `json:"source"`
			SourceID        *string   `json:"source_id"`
			Description     *string   `json:"description"`
			BalanceBefore   int       `json:"balance_before"`
			BalanceAfter    int       `json:"balance_after"`
			CreatedAt       time.Time `json:"created_at"`
		}

		var data []byte
		data, _, err = supabase.From("credit_transactions").
			Select("*", "", false).
			Eq("id", fmt.Sprintf("%d", *result.TransactionID)).
			Execute()

		if err != nil {
			return nil, fmt.Errorf("error querying credit_transaction: %w", err)
		}

		if err := json.Unmarshal(data, &transaction); err != nil {
			return nil, fmt.Errorf("error unmarshaling credit_transaction: %w", err)
		}

		if len(transaction) == 0 {
			return nil, fmt.Errorf("credit_transaction not found")
		}

		parsedUserID, err := uuid.Parse(transaction[0].UserID)
		if err != nil {
			return nil, fmt.Errorf("error parsing user_id: %w", err)
		}

		return &billing.CreditTransaction{
			ID:              transaction[0].ID,
			UserID:          parsedUserID,
			Amount:          transaction[0].Amount,
			TransactionType: transaction[0].TransactionType,
			Source:          transaction[0].Source,
			SourceID:        transaction[0].SourceID,
			Description:     transaction[0].Description,
			BalanceBefore:   transaction[0].BalanceBefore,
			BalanceAfter:    transaction[0].BalanceAfter,
			CreatedAt:       transaction[0].CreatedAt,
		}, nil
	}
}
