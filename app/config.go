package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	DefaultConfigName = "local_config"
	DefaultConfigType = "yaml"
	LocalConfigPath   = "local_configuration"
)
const (
	DBDriver   = "DB.Driver"
	DBUser     = "DB.DBUser"
	DBPassword = "DB.DBPassword"
	DBName     = "DB.DBName"
	Port       = "Port"
)

var ErrInvalidConfig = errors.New("invalid config")

var requiredConfig = []string{
	DBDriver,
	DBUser,
	DBPassword,
	DBName,
	Port,
}

type DatabaseOptions struct {
	Driver,
	User,
	Password,
	DbName string
}

type Configuration struct {
	ContactEmail,
	Port string
	Database DatabaseOptions
}

func GetConfiguration() (*Configuration, error) {
	viper.SetConfigName(DefaultConfigName)
	viper.SetConfigType(DefaultConfigType)
	homePath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	localConfigPath := fmt.Sprintf("%v/%v", homePath, LocalConfigPath)
	viper.AddConfigPath(LocalConfigPath)
	viper.AddConfigPath(localConfigPath)
	err = viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %s.%s: %w", DefaultConfigName, DefaultConfigType, ErrInvalidConfig)
	}

	if err := validateConfig(requiredConfig); err != nil {
		return nil, err
	}

	return &Configuration{
		Database: DatabaseOptions{
			Driver:   viper.GetString(DBDriver),
			User:     viper.GetString(DBUser),
			Password: viper.GetString(DBPassword),
			DbName:   viper.GetString(DBName),
		},
	}, nil
}

func validateConfig(requiredParams []string) error {
	var missingParams []string

	for _, p := range requiredParams {
		if ok := viper.IsSet(p); !ok {
			missingParams = append(missingParams, p)
		}
	}

	if len(missingParams) > 0 {
		return fmt.Errorf("required parameters are missing: %q: %w", missingParams, ErrInvalidConfig)
	}

	return nil
}
