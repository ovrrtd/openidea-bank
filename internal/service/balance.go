package service

import (
	"context"
	"net/http"

	"github.com/ovrrtd/openidea-bank/internal/helper/common"
	"github.com/ovrrtd/openidea-bank/internal/helper/errorer"
	"github.com/ovrrtd/openidea-bank/internal/helper/validator"
	"github.com/ovrrtd/openidea-bank/internal/model/entity"
	"github.com/ovrrtd/openidea-bank/internal/model/request"
	"github.com/ovrrtd/openidea-bank/internal/model/response"
	"github.com/pkg/errors"
)

func (s *service) AddBalance(ctx context.Context, payload request.AddBalance) (int, error) {
	err := validator.ValidateStruct(&payload)
	if err != nil {
		return http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	// validate url proof image with regex
	ok := common.ValidateUrl(payload.ProofImageURL)
	if !ok {
		return http.StatusBadRequest, errors.Wrap(errorer.ErrBadRequest, errorer.ErrBadRequest.Error())
	}

	code, err := s.balanceRepo.UpsertBalance(ctx, entity.UpsertBalance{
		UserID:                  payload.UserID,
		Balance:                 payload.Balance,
		Currency:                payload.Currency,
		ProofImageURL:           payload.ProofImageURL,
		SenderBankAccountNumber: payload.BankAccountNumber,
		SenderBankName:          payload.BankName,
		TransactionID:           "",
	})

	return code, err

}
func (s *service) CreateTransaction(ctx context.Context, payload request.CreateTransaction) (int, error) {
	err := validator.ValidateStruct(&payload)
	if err != nil {
		return http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	code, err := s.balanceRepo.UpsertBalance(ctx, entity.UpsertBalance{
		UserID:                  payload.UserID,
		Balance:                 -payload.Balance,
		Currency:                payload.Currency,
		SenderBankAccountNumber: payload.BankAccountNumber,
		SenderBankName:          payload.BankName,
		TransactionID:           common.GenerateULID(),
	})

	return code, err
}

func (s *service) GetBalances(ctx context.Context, userId string) ([]response.Balance, int, error) {
	entBalance, code, err := s.balanceRepo.GetBalances(ctx, userId)
	if err != nil {
		return nil, code, err
	}
	balances := make([]response.Balance, len(entBalance))

	for i, v := range entBalance {
		balances[i] = response.Balance{
			Balance:  v.Balance,
			Currency: v.Currency,
		}
	}

	return balances, code, nil
}

func (s *service) GetBalancesHistory(ctx context.Context, payload request.GetBalancesHistory) ([]response.GetBalancesHistory, int, error) {
	err := validator.ValidateStruct(&payload)
	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	entBH, code, err := s.balanceRepo.GetBalancesHistory(ctx, entity.GetBalancesHistory{
		UserID: payload.UserID,
		Limit:  payload.Limit,
		Offset: payload.Offset,
	})

	if err != nil {
		return nil, code, err
	}

	gh := make([]response.GetBalancesHistory, len(entBH))
	for i, v := range entBH {
		gh[i] = response.GetBalancesHistory{
			TransactionID: v.TransactionID,
			Balance:       v.Balance,
			Currency:      v.Currency,
			CreatedAt:     v.CreatedAt,
			Source: struct {
				BankAccountNumber string `json:"bankAccountNumber"`
				BankName          string `json:"bankName"`
			}{
				BankAccountNumber: v.SourceBankAccountNumber,
				BankName:          v.SourceBankName,
			},
		}
	}

	return gh, code, nil
}
