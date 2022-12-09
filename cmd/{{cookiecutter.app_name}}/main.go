package main

import (
	"{{cookiecutter.app_name}}/internal/config"
	"{{cookiecutter.app_name}}/internal/models"
	"{{cookiecutter.app_name}}/internal/router"
	"{{cookiecutter.app_name}}/internal/storage"
	"{{cookiecutter.app_name}}/pkg/constants"
	"{{cookiecutter.app_name}}/pkg/logger"

	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	runner "github.com/oklog/run"
)

func main() {
	// Set current working directory to make logger and config use the application dir
	err := os.Chdir(filepath.Dir(appFilePath()))
	if err != nil {
		logger.Fatalf("os.Chdir failed error: %v", err)
	}

	cfg, err := config.Init(constants.ConfigPath, constants.ConfigName)
	if err != nil {
		logger.Fatalf("config.Init: %s", err)
	}

	db, err := storage.DBConn(&cfg.Database)
	if err != nil {
		logger.Fatalf("storage.DBConn: %s", err)
	}

	s := storage.New(db)

	httpServer := &http.Server{
		MaxHeaderBytes: 10, // 10 MB
		Addr:           ":" + cfg.Web.Port,
		WriteTimeout:   time.Second * time.Duration(cfg.Web.Timeout),
		ReadTimeout:    time.Second * time.Duration(cfg.Web.Timeout),
		IdleTimeout:    time.Second * 60,
		Handler:        router.New(s),
	}

	var serverGroup runner.Group

	err = func() error {
		serverGroup.Add(func() error {
			logger.Infof("http server started at port: %s", cfg.Web.Port)

			return httpServer.ListenAndServe()
		}, func(err error) {
			logger.Errorf("Error start http server: %s", err)
		})

		msg := fmt.Sprintf("{{cookiecutter.app_name}} is up and running on '%s' in '%s' mode", cfg.Server.Port, cfg.Server.Env)
		fmt.Println(msg)

		return serverGroup.Run()
	}()

	if err != nil {
		logger.Errorf("bad start app: %s", err)
	}
}

// appFilePath returns the file path of the executable that is currently running
func appFilePath() string {
	path, err := os.Executable()
	if err != nil {
		// Fallback to args array which may not always be the full path
		return os.Args[0]
	}
	return path
}
