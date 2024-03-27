package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func (s *service) UploadImage(ctx context.Context, file *multipart.FileHeader) (string, int, error) {
	s.log.Debug().Msgf("file size: %d", file.Size)
	if file.Size >= 2_000_000 || file.Size <= 10_000 {
		return "", http.StatusBadRequest, errors.Wrap(errors.New("file size must between 10KB and 2MB"), "file size must between 10KB and 2MB")
	}

	if !strings.HasSuffix(strings.ToLower(file.Filename), ".jpg") && !strings.HasSuffix(strings.ToLower(file.Filename), ".jpeg") {
		return "", http.StatusBadRequest, errors.Wrap(errors.New("file type not allowed"), "file type not allowed")
	}

	src, err := file.Open()
	if err != nil {
		return "", http.StatusBadRequest, err
	}
	defer src.Close()

	imageUrl, code, err := s.s3Repo.UploadFile(ctx, fmt.Sprintf("%d-%s", time.Now().UnixMilli(), file.Filename), src)

	return imageUrl, code, err
}
