package restapi

import (
	"github.com/ovrrtd/openidea-bank/internal/delivery/middleware"
	"github.com/ovrrtd/openidea-bank/internal/service"

	"github.com/rs/zerolog"
)

type Restapi struct {
	log        zerolog.Logger
	middleware middleware.Middleware
	service    service.Service
}

func New(
	log zerolog.Logger,
	middleware middleware.Middleware,
	s service.Service,
) *Restapi {
	return &Restapi{
		log:        log,
		middleware: middleware,
		service:    s,
	}
}

func (r *Restapi) debugError(err error) {
	if err != nil {
		r.log.Debug().Stack().Err(err).Send()
	}
}
