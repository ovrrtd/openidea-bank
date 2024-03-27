package repository

import (
	"context"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/ovrrtd/openidea-bank/internal/helper/errorer"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type S3Repository interface {
	UploadFile(ctx context.Context, filename string, file multipart.File) (string, int, error)
}

func NewS3Repository(logger zerolog.Logger) S3Repository {
	creds := credentials.NewStaticCredentialsProvider(os.Getenv("S3_ID"), os.Getenv("S3_SECRET_KEY"), "")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds), config.WithRegion("ap-southeast-1"))
	if err != nil {
		logger.Error().Err(err).Msg("Failed to load config")
	}

	return &S3RepositoryImpl{
		logger:      logger,
		awsS3Client: s3.NewFromConfig(cfg),
	}
}

type S3RepositoryImpl struct {
	logger      zerolog.Logger
	awsS3Client *s3.Client
}

func (s *S3RepositoryImpl) UploadFile(ctx context.Context, filename string, file multipart.File) (string, int, error) {
	uploader := manager.NewUploader(s.awsS3Client)
	result, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(filename),
		ACL:    "public-read",
		Body:   file,
	})

	if err != nil {
		s.logger.Debug().Msg(err.Error())
		return "", http.StatusInternalServerError, errors.Wrap(errorer.ErrInternalServer, err.Error())
	}

	return result.Location, http.StatusOK, nil
}
