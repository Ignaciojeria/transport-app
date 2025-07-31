package tidbrepository

import (
	"context"
	"encoding/json"
	"time"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/database"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

type UpsertWebhook func(ctx context.Context, webhook domain.Webhook, fsmState domain.FSMState) error

func init() {
	ioc.Registry(NewUpsertWebhook, database.NewConnectionFactory, NewSaveFSMTransition)
}

func NewUpsertWebhook(conn database.ConnectionFactory, saveFSMTransition SaveFSMTransition) UpsertWebhook {
	return func(ctx context.Context, webhook domain.Webhook, fsmState domain.FSMState) error {
		db, err := conn()
		if err != nil {
			return err
		}
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		now := time.Now()
		webhook.CreatedAt = now
		webhook.UpdatedAt = now

		// Upsert webhook
		query := `
		INSERT INTO webhooks (
			doc_id, type, url, headers, max_retries, backoff_seconds,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			type = VALUES(type),
			url = VALUES(url),
			headers = VALUES(headers),
			max_retries = VALUES(max_retries),
			backoff_seconds = VALUES(backoff_seconds),
			updated_at = VALUES(updated_at)
		`

		headersJSON := "{}"
		if len(webhook.Headers) > 0 {
			headersBytes, err := json.Marshal(webhook.Headers)
			if err != nil {
				return err
			}
			headersJSON = string(headersBytes)
		}

		_, err = tx.ExecContext(ctx, query,
			webhook.DocID(ctx).String(),
			webhook.Type,
			webhook.URL,
			headersJSON,
			webhook.RetryPolicy.MaxRetries,
			webhook.RetryPolicy.BackoffSeconds,
			webhook.CreatedAt,
			webhook.UpdatedAt,
		)
		if err != nil {
			return err
		}

		// Save FSM transition
		if err := saveFSMTransition(ctx, tx, fsmState); err != nil {
			return err
		}

		return tx.Commit()
	}
}

