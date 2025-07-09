package extutils

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"

	"abodemine/domains/arc"
	"abodemine/lib/errors"
)

func RollbackPgxTx(ctx context.Context, tx pgx.Tx, id string) {
	if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
		log.Error().
			Err(err).
			Str("id", id).
			Msg("Failed to rollback transaction.")
	}
}

func PgxExec(r *arc.Request, key string, sql string, args []any) (*pgconn.CommandTag, error) {
	if tx, ok := r.SelectPgxTx(key); ok {
		ct, err := tx.Exec(r.Context(), sql, args...)
		if err != nil {
			return nil, &errors.Object{
				Id:     "db499ba7-a19e-4a10-9cb8-eab9c592ffd3",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to query tx.",
				Cause:  err.Error(),
			}
		}

		return &ct, nil
	}

	if pool, err := r.Dom().SelectPgxPool(key); err == nil {
		ct, err := pool.Exec(r.Context(), sql, args...)
		if err != nil {
			return nil, &errors.Object{
				Id:     "1b4ca908-0121-469e-828b-87eedd6d3177",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to query pool.",
				Cause:  err.Error(),
			}
		}

		return &ct, nil
	}

	return nil, &errors.Object{
		Id:     "9fbac405-21c0-4abb-bd39-825476a6e7db",
		Code:   errors.Code_UNKNOWN,
		Detail: "Missing database handler.",
	}
}

func PgxQuery(r *arc.Request, key string, sql string, args []any) (pgx.Rows, error) {
	if tx, ok := r.SelectPgxTx(key); ok {
		rows, err := tx.Query(r.Context(), sql, args...)
		if err != nil {
			return nil, &errors.Object{
				Id:     "a4099b2e-354c-454a-994f-8a13bd63d432",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to query tx.",
				Cause:  err.Error(),
			}
		}

		return rows, nil
	}

	if pool, err := r.Dom().SelectPgxPool(key); err == nil {
		rows, err := pool.Query(r.Context(), sql, args...)
		if err != nil {
			return nil, &errors.Object{
				Id:     "f6d6ec5e-591c-4c56-86ff-a61e9249fa2b",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to query pool.",
				Cause:  err.Error(),
			}
		}

		return rows, nil
	}

	return nil, &errors.Object{
		Id:     "f9a6d5d8-cf3b-45bb-87ae-847a2cf81de9",
		Code:   errors.Code_UNKNOWN,
		Detail: "Missing database handler.",
	}
}

func PgxQueryRow(r *arc.Request, key string, sql string, args []any) (pgx.Row, error) {
	if tx, ok := r.SelectPgxTx(key); ok {
		return tx.QueryRow(r.Context(), sql, args...), nil
	}

	if pool, err := r.Dom().SelectPgxPool(key); err == nil {
		return pool.QueryRow(r.Context(), sql, args...), nil
	}

	return nil, &errors.Object{
		Id:     "16e646f6-3116-4d4c-bb35-9653092bf08d",
		Code:   errors.Code_UNKNOWN,
		Detail: "Missing database handler.",
	}
}
