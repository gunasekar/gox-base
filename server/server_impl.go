package goxServer

import (
	"fmt"
	"github.com/devlibx/gox-base"
	"github.com/devlibx/gox-base/config"
	"github.com/devlibx/gox-base/errors"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type serverImpl struct {
	server *http.Server
	gox.CrossFunction
}

func (s *serverImpl) Start(handler http.Handler, applicationConfig *config.App) error {
	if applicationConfig == nil {
		return errors.New("application config is nil")
	}

	// Setup default values
	applicationConfig.SetupDefaults()

	// Setup server
	rootHandler := negroni.Classic()
	rootHandler.Use(s.setupTimeLogging())
	rootHandler.UseHandler(handler)

	// Setup http server
	s.server = &http.Server{
		Handler:      rootHandler,
		Addr:         fmt.Sprintf("0.0.0.0:%d", applicationConfig.HttpPort),
		WriteTimeout: time.Duration(applicationConfig.RequestWriteTimeoutMs) * time.Millisecond,
		ReadTimeout:  time.Duration(applicationConfig.RequestReadTimeoutMs) * time.Millisecond,
		IdleTimeout:  time.Duration(applicationConfig.IdleTimeoutMs) * time.Millisecond,
	}

	return s.server.ListenAndServe()
}

func (s *serverImpl) setupTimeLogging() negroni.HandlerFunc {
	logger := s.Logger().Named("negroni")
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		start := time.Now()
		next(rw, r)
		end := time.Now()
		logger.Info("",
			zap.String("remoteAddr", r.RemoteAddr),
			zap.String("source", r.Header.Get("X-FORWARDED-FOR")),
			zap.Int64("duration", end.Sub(start).Milliseconds()),
		)
	}
}
