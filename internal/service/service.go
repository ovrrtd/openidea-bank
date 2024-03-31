package service

import (
	"context"
	"mime/multipart"

	"github.com/ovrrtd/openidea-bank/internal/model/request"
	"github.com/ovrrtd/openidea-bank/internal/model/response"
	"github.com/ovrrtd/openidea-bank/internal/repository"

	"github.com/rs/zerolog"
)

type Service interface {
	// User
	Register(ctx context.Context, payload request.Register) (*response.Register, int, error)
	Login(ctx context.Context, payload request.Login) (*response.Login, int, error)
	GetUserByID(ctx context.Context, id string) (*response.User, int, error)
	// s3
	UploadImage(ctx context.Context, file *multipart.FileHeader) (string, int, error)

	// Balance
	AddBalance(ctx context.Context, payload request.AddBalance) (int, error)
	CreateTransaction(ctx context.Context, payload request.CreateTransaction) (int, error)
	GetBalances(ctx context.Context, userId string) ([]response.Balance, int, error)
	GetBalancesHistory(ctx context.Context, payload request.GetBalancesHistory) ([]response.GetBalancesHistory, int, error)
}

type Config struct {
	Salt      int
	JwtSecret string
}

type service struct {
	cfg         Config
	log         zerolog.Logger
	userRepo    repository.UserRepository
	s3Repo      repository.S3Repository
	balanceRepo repository.BalanceRepository
}

func New(
	cfg Config,
	logger zerolog.Logger,
	userRepo repository.UserRepository,
	s3Repo repository.S3Repository,
	balanceRepo repository.BalanceRepository,
) Service {
	return &service{
		cfg:         cfg,
		log:         logger,
		userRepo:    userRepo,
		s3Repo:      s3Repo,
		balanceRepo: balanceRepo,
	}
}
