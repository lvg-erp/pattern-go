package web

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"server/config"
	"time"
)

// Dependencies зависимости сервера
type Dependencies struct {
	Logger logrus.FieldLogger
	Config *config.Config
}

// Options опции сервера
type Options struct {
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// Server HTTP-сервер
type Server struct {
	logger   logrus.FieldLogger
	srvHTTPS *http.Server
	srvHTTP  *http.Server

	shutdownTimeout time.Duration
	config          *config.Config
	certManager     autocert.Manager
}

func NewServer(opt Options, dep Dependencies) *Server {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("./"),
		HostPolicy: autocert.HostWhitelist(dep.Config.WebServer.HttpsSettings.Domains...),
		Email:      dep.Config.WebServer.HttpsSettings.AdminEmail,
	}

	tLSConfig := certManager.TLSConfig()
	tLSConfig.NextProtos = append([]string{"h2", "http/1.1"}, tLSConfig.NextProtos...)

	s := &Server{
		srvHTTPS: &http.Server{
			Addr:         ":" + dep.Config.WebServer.ListenPorts.Https,
			TLSConfig:    tLSConfig,
			ReadTimeout:  opt.ReadTimeout,
			WriteTimeout: opt.WriteTimeout,
		},

		srvHTTP: &http.Server{
			Addr:         ":" + dep.Config.WebServer.ListenPorts.Http,
			ReadTimeout:  opt.ReadTimeout,
			WriteTimeout: opt.WriteTimeout,
		},

		logger:          dep.Logger,
		shutdownTimeout: opt.ShutdownTimeout,
		config:          dep.Config,
		certManager:     certManager,
	}

	return s
}

func (s *Server) Start() error {

	if s.config.WebServer.HttpsSettings.Enabled {
		s.logger.Error("*** Web server startn with SSL/TLS - Enabled")
		s.logger.Error("******************************************")
		go func() {
			err := s.srvHTTP.ListenAndServe()
			if err != nil {
				s.logger.Errorf("TLS - Error %w", err)
			}
		}()
		err := s.srvHTTPS.ListenAndServeTLS("", "")
		if err != nil {
			s.logger.Errorf("TLS - Error %w", err)
			if err == http.ErrServerClosed {
				return nil
			}
		}
		return err
	} else {
		s.logger.Error("*** Web server start with SSL/TLS - Disabled")
		s.logger.Error("******************************************")
		err := s.srvHTTP.ListenAndServe()
		if err != nil {
			s.logger.Errorf("busy port %w", err)
		}
		return nil
	}
}

func (s *Server) Stop() error {
	var ctx context.Context
	var cancelCtx context.CancelFunc

	if s.shutdownTimeout > 0 {
		ctx, cancelCtx = context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancelCtx()
	} else {
		ctx = context.Background()
	}

	err := s.srvHTTP.Shutdown(ctx)
	if err != nil {
		return err
	}

	err = s.srvHTTP.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
