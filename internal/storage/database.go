package storage

import (
	"{{cookiecutter.project_slug}}/internal/config"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	db *gorm.DB
}

// DBConn databese connection
func DBConn(cfg *config.DatabaseConfiguration) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	logmode := viper.GetBool("db.logmode")
	loglevel := logger.Silent
	if logmode {
		loglevel = logger.Info
	}

	newDBLogger := logger.New(
		log.New(getWriter(), "\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  loglevel,    // Log level (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	db, err = gorm.Open(sqlite.Open(cfg.Name), &gorm.Config{Logger: newDBLogger})
	if err != nil {
		return nil, fmt.Errorf("could not open sqlite3 database: %v", err)
	}

	return db, err
}

// New opens a database according to configuration.
func New(db *gorm.DB) *Database {
	return &Database{
		db: db,
	}
}

func getWriter() io.Writer {
	file, err := os.OpenFile("{{cookiecutter.project_slug}}.db.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return os.Stdout
	} else {
		return file
	}
}
