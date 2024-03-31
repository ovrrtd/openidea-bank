package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/mux"
	database "github.com/ovrrtd/openidea-bank/db"
	mw "github.com/ovrrtd/openidea-bank/internal/delivery/middleware"
	"github.com/ovrrtd/openidea-bank/internal/delivery/restapi"
	"github.com/ovrrtd/openidea-bank/internal/repository"
	"github.com/ovrrtd/openidea-bank/internal/service"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	_ "github.com/lib/pq"
)

var (
	APP_PORT = "8080"
)

func Server() error {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := zerolog.New(os.Stdout)
	// db, err := newMongoDB(ConfigMongoDB{Host: cfg.DB.Host})
	db, err := database.NewDBDefaultSql()
	if err != nil {
		logger.Info().Msg(fmt.Sprintf("Postgres connection error: %s", err.Error()))
		return err
	}
	err = db.Ping()
	if err != nil {
		logger.Info().Msg(fmt.Sprintf("Postgres ping error: %s", err.Error()))
		return err
	}
	defer db.Close()

	// repository init
	userRepo := repository.NewUserRepository(logger, db)
	balanceRepo := repository.NewBalanceRepository(logger, db)
	s3Repo := repository.NewS3Repository(logger)

	salt, err := strconv.Atoi(os.Getenv("BCRYPT_SALT"))
	if err != nil {
		salt = 8
	}
	// service registry
	service := service.New(
		service.Config{Salt: salt, JwtSecret: os.Getenv("JWT_SECRET")},
		logger,
		userRepo,
		s3Repo,
		balanceRepo,
	)

	// middleware init
	md := mw.New(logger, service)

	// restapi init
	rest := restapi.New(logger, md, service)

	router := mux.NewRouter()

	// add restapi route
	rest.MakeRoute(router)

	errs := make(chan error)
	go func() {
		logger.Log().Msg(fmt.Sprintf("start server on port %s", APP_PORT))
		errs <- http.ListenAndServe(":8080", md.RemoveTrailingSlash(router))
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	return <-errs
}
