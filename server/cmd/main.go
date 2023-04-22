package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"server/config"
	"server/internal/repo"
	"server/internal/scheduler"
	"server/internal/server/web"
	"syscall"
	"time"
)

const Version = "v.1.0"

func main() {
	logger := logrus.New()
	if logger == nil {
		logger.Fatal("New logger failed")
	}

	conf := config.NewConfig()
	if conf == nil {
		logger.Fatal("New config failed")
	}

	logger.SetLevel(conf.Logger.Level)

	credentials := config.NewCredentials()
	if credentials == nil {
		logger.Fatal("New credentials failed")
	}

	DBConnect, err := repo.NewDBConnect(credentials)
	if err != nil {
		logger.WithError(err).Fatal("New db connect failed")
	}

	logger.Error("*************************************************")
	logger.Errorf("*** Start Web Crawler - Version: %v", Version)

	repositoryRegistry := InitRepoRegistry(DBConnect, conf, logger)

	scheduler, err := scheduler.NewScheduler(scheduler.InitSchedulerDep{
		Conf:               conf,
		Logger:             logger,
		RepositoryRegistry: repositoryRegistry,
	})
	if err != nil {
		logger.Error(err)
	}
	scheduler.Start()

	srv := web.NewServer(web.Options{
		ReadTimeout:     time.Second * 3,
		WriteTimeout:    time.Second * 30,
		ShutdownTimeout: time.Second * 10,
	}, web.Dependencies{
		Logger: logger,
		Config: conf,
	})
	go func() {
		err := srv.Start()
		if err != nil {
			logger.WithError(err).Fatal("failed to start server")
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	logger.Error("Server Stopped")

	err = srv.Stop()
	if err != nil {
		logger.WithError(err).Fatalf("failed to stop server: %v", err)
	}
	time.Sleep(time.Second * 3)
	logger.Error("Server Exited Properly")
}

func InitRepoRegistry(DBConnect *sqlx.DB, conf *config.Config, logger logrus.FieldLogger) *repo.Registry {
	return &repo.Registry{}
}
