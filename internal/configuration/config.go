package configuration

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
	DBDriver string = "DB.Driver"
	Port     string = "Port"
)

var ErrInvalidConfig = errors.New("invalid config")

var requiredConfig = []string{
	DBDriver,
	Port,
}

type Configuration struct {
	ContactEmail string
	DBDriver     string
	Port         string
}

func New() (*Configuration, error) {
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
		DBDriver: viper.GetString(DBDriver),
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
