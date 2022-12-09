package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	configuration *Configuration
	configFileExt = ".yml"
	configType    = "yaml"
)

// Configuration ...
type Configuration struct {
	Web      WebConfiguration
	Database DatabaseConfiguration
}

// WebConfiguration is the required parameters to set up a web server
type WebConfiguration struct {
	Port     string `default:"8080"`
	Timeout  int    `default:"24"`
	Username string `default:"brick"`
	Password string `default:"brick1337"`
}

// DatabaseConfiguration is the required parameters to set up a DB instance
type DatabaseConfiguration struct {
	Name    string `default:"db.sqlite3"`
	LogMode bool   `default:"false"`
}

// Init initializes the configuration manager
func Init(configPath, configName string) (*Configuration, error) {

	configFilePath := filepath.Join(configPath, configName) + configFileExt

	// initialize viper configuration
	initializeConfig(configPath, configName)

	// Bind environment variables
	bindEnvs()

	// Set default values
	setDefaults()

	// Read or create configuration file
	if err := readConfiguration(configFilePath); err != nil {
		return nil, err
	}

	// Auto read env variables
	viper.AutomaticEnv()

	// Unmarshal config file to struct
	if err := viper.Unmarshal(&configuration); err != nil {
		return nil, err
	}

	return configuration, nil
}

// read configuration from file
func readConfiguration(configFilePath string) error {
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		// if file does not exist, simply create one
		if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
			os.Create(configFilePath)
		} else {
			return err
		}
		// let's write defaults
		if err := viper.WriteConfig(); err != nil {
			return err
		}
	}
	return nil
}

// initialize the configuration manager
func initializeConfig(configPath, configName string) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
}

func bindEnvs() {
	viper.BindEnv("web.port", "BRK_WEB_PORT")
	viper.BindEnv("web.timeout", "BRK_WEB_TIMEOUT")
	viper.BindEnv("web.username", "BRK_WEB_USERNAME")
	viper.BindEnv("web.password", "BRK_WEB_PASSWORD")

	viper.BindEnv("db.name", "BRK_DB_NAME")
	viper.BindEnv("db.logmode", "BRK_DB_LOG_MODE")
}

func setDefaults() {
	// Web defaults
	viper.SetDefault("web.port", "9090")
	viper.SetDefault("web.timeout", 24)
	viper.SetDefault("web.username", "username")
	viper.SetDefault("web.password", "password")

	// Database defaults
	viper.SetDefault("db.name", "db.sqlite3")
	viper.SetDefault("db.logmode", false)
}
