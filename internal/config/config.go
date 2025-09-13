package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Env     Environment  `mapstructure:"env" validate:"required,oneof=local dev stage prod"`
}

const (
	envPrefix = "APP"
	envKeyEnv = envPrefix + "_ENV"
	configFilePrefix = "env."
	configType = "yaml"
)

var configPaths = []string{".", "./configs", "/etc/app"}

func GetConfig() (*Config, error) {
	var env Environment
	if appEnv := os.Getenv(envKeyEnv); appEnv == "" {
		log.Printf("%s is not set, defaulting to '%s'", envKeyEnv, EnvLocal.String())
		env = EnvLocal
	}

	v := viper.New()	

	v.SetConfigName(configFilePrefix + env.String())
	v.SetConfigType(configType)
	
	// Add multiple config paths
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	if err := v.ReadInConfig(); err != nil {
		var nf viper.ConfigFileNotFoundError
		// file found, but error parsing/access permission
    	if !errors.As(err, &nf) {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	} else {
		log.Printf("loaded config file: %s", v.ConfigFileUsed())
	}

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Set env from variable if not in config
	if config.Env == "" {
		config.Env = env
	}

	// Validate config
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, err
	}

	return &config, nil
}

var validate = validator.New()

func (c *Config) Validate() error {
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	return nil
}
