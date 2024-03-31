package repository

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/ovrrtd/openidea-bank/internal/helper/common"
	"github.com/ovrrtd/openidea-bank/internal/helper/errorer"
	"github.com/ovrrtd/openidea-bank/internal/model/entity"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type BalanceRepository interface {
	UpsertBalance(ctx context.Context, payload entity.UpsertBalance) (int, error)
	GetBalances(ctx context.Context, userId string) ([]entity.Balance, int, error)
	GetBalancesHistory(ctx context.Context, payload entity.GetBalancesHistory) ([]entity.BalanceHistory, int, error)
}

func NewBalanceRepository(logger zerolog.Logger, db *sql.DB) BalanceRepository {
	return &BalanceRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

type BalanceRepositoryImpl struct {
	logger zerolog.Logger
	db     *sql.DB
}

func (r *BalanceRepositoryImpl) UpsertBalance(ctx context.Context, payload entity.UpsertBalance) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	defer tx.Rollback()

	balance := entity.Balance{}
	err = tx.QueryRowContext(ctx, `SELECT id, balance FROM balances WHERE user_id = $1 AND currency = $2`, payload.UserID, payload.Currency).Scan(&balance.ID, &balance.Balance)

	if err != nil && err != sql.ErrNoRows {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	if err == sql.ErrNoRows {
		query := `
		INSERT INTO balances (id, balance, user_id, currency) VALUES ($1, $2, $3, $4)
	`

		_, err = tx.ExecContext(ctx, query,
			common.GenerateULID(),
			payload.Balance,
			payload.UserID,
			payload.Currency,
		)
		if err != nil {
			return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
		}
	} else {
		if balance.Balance+payload.Balance < 0 {
			return http.StatusBadRequest, errors.Wrap(errorer.ErrBadRequest, "balance is not enough")
		}

		query := `
		UPDATE balances SET balance=balance + $1 WHERE id = $2
	`

		_, err = tx.ExecContext(ctx, query,
			payload.Balance,
			balance.ID,
		)

		if err != nil {
			return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
		}
	}

	query := `INSERT INTO BALANCES_HISTORY (id, transaction_id, user_id, balance, currency, proof_image_url,source_bank_account_number, source_bank_name, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err = tx.ExecContext(ctx, query,
		common.GenerateULID(),
		payload.TransactionID,
		payload.UserID,
		payload.Balance,
		payload.Currency,
		payload.ProofImageURL,
		payload.SenderBankAccountNumber,
		payload.SenderBankName,
		time.Now().UnixMilli(),
	)

	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}
	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	return http.StatusOK, nil
}

func (r *BalanceRepositoryImpl) GetBalances(ctx context.Context, userId string) ([]entity.Balance, int, error) {
	var balances []entity.Balance
	query := `SELECT id, balance, currency FROM balances WHERE user_id = $1 ORDER BY balance DESC`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	for rows.Next() {
		var balance entity.Balance
		if err := rows.Scan(&balance.ID, &balance.Balance, &balance.Currency); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
		}
		balances = append(balances, balance)
	}

	if err := rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	return balances, http.StatusOK, nil
}

func (r *BalanceRepositoryImpl) GetBalancesHistory(ctx context.Context, payload entity.GetBalancesHistory) ([]entity.BalanceHistory, int, error) {
	var balances []entity.BalanceHistory
	query := `SELECT id, transaction_id, user_id, balance, currency, proof_image_url,source_bank_account_number, source_bank_name, created_at FROM balances_history WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, payload.UserID, payload.Limit, payload.Offset)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
	}

	for rows.Next() {
		bh := entity.BalanceHistory{}
		if err := rows.Scan(&bh.ID, &bh.TransactionID, &bh.UserID, &bh.Balance, &bh.Currency, &bh.ProofImageURL, &bh.SourceBankAccountNumber, &bh.SourceBankName, &bh.CreatedAt); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalDatabase, err.Error())
		}
		balances = append(balances, bh)
	}

	return balances, http.StatusOK, nil
}
