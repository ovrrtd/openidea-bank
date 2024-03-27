package service

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/ovrrtd/openidea-bank/internal/helper/common"
	"github.com/ovrrtd/openidea-bank/internal/helper/errorer"
	"github.com/ovrrtd/openidea-bank/internal/helper/jwt"
	"github.com/ovrrtd/openidea-bank/internal/helper/validator"
	"github.com/ovrrtd/openidea-bank/internal/model/entity"
	"github.com/ovrrtd/openidea-bank/internal/model/request"
	"github.com/ovrrtd/openidea-bank/internal/model/response"

	jwtV5 "github.com/golang-jwt/jwt/v5"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Register to register a new user by email and password
func (s *service) Register(ctx context.Context, payload request.Register) (*response.Register, int, error) {
	err := validator.ValidateStruct(&payload)

	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}
	ent := entity.User{
		Name: payload.Name,
	}

	// validate email form
	regex := regexp.MustCompile(common.RegexEmailPattern)
	if !regex.MatchString(payload.Email) {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInvalidEmail, errorer.ErrInvalidEmail.Error())
	}

	exist, code, err := s.userRepo.FindByEmail(ctx, payload.Email)

	if err != nil && code != http.StatusNotFound {
		return nil, code, err
	}
	if exist != nil {
		return nil, http.StatusConflict, errors.Wrap(errorer.ErrEmailExist, errorer.ErrEmailExist.Error())
	}

	ent.Email = payload.Email

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), s.cfg.Salt)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, err.Error())
	}
	ent.Password = string(hashedPassword)

	user, code, err := s.userRepo.Register(ctx, ent)

	if err != nil {
		return nil, code, err
	}

	// TODO: generate access token
	userClaims := common.UserClaims{
		Id: user.ID,
		RegisteredClaims: jwtV5.RegisteredClaims{
			IssuedAt:  jwtV5.NewNumericDate(time.Now()),
			ExpiresAt: jwtV5.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	}
	tokenString, err := jwt.GenerateJwt(userClaims, s.cfg.JwtSecret)

	if err != nil {
		return nil, code, errors.Wrap(err, err.Error())
	}

	return &response.Register{
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: tokenString,
	}, code, nil
}

func (s *service) Login(ctx context.Context, payload request.Login) (*response.Login, int, error) {
	err := validator.ValidateStruct(&payload)

	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInputRequest(err), errorer.ErrInputRequest(err).Error())
	}

	// validate email form
	regex := regexp.MustCompile(common.RegexEmailPattern)
	if !regex.MatchString(payload.Email) {
		return nil, http.StatusBadRequest, errors.Wrap(errorer.ErrInvalidEmail, errorer.ErrInvalidEmail.Error())
	}

	user, code, err := s.userRepo.FindByEmail(ctx, payload.Email)

	if err != nil {
		return nil, code, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return nil, http.StatusBadRequest, errors.Wrap(err, err.Error())
	}

	userClaims := common.UserClaims{
		Id: user.ID,
		RegisteredClaims: jwtV5.RegisteredClaims{
			IssuedAt:  jwtV5.NewNumericDate(time.Now()),
			ExpiresAt: jwtV5.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
		},
	}
	tokenString, err := jwt.GenerateJwt(userClaims, s.cfg.JwtSecret)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(err, err.Error())
	}

	return &response.Login{
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: tokenString,
	}, http.StatusOK, nil
}

func (s *service) GetUserByID(ctx context.Context, id int64) (*response.User, int, error) {
	user, code, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, code, err
	}
	return &response.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, code, nil
}
