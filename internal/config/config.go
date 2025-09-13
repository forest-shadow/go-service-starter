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
	Env     Environment `mapstructure:"env" validate:"required,oneof=local dev stage prod"`
	Logging Logging     `mapstructure:"logging"`
}

const (
	envPrefix        = "APP"
	envKeyEnv        = envPrefix + "_ENV"
	configFilePrefix = "env."
	configType       = "yaml"
)

var configPaths = []string{".", "./configs", "/etc/app"}

func GetConfig() (*Config, error) {
	var env Environment
	if appEnv := os.Getenv(envKeyEnv); appEnv == "" {
		log.Printf("%s is not set, defaulting to '%s'", envKeyEnv, EnvLocal.String())
		env = EnvLocal
	} else {
		env = Environment(appEnv)
		// Validate environment value early
		if !env.IsValid() {
			return nil, fmt.Errorf("invalid environment '%s', must be one of: local, dev, stage, prod", appEnv)
		}
	}

	v := viper.New()

	v.SetConfigName(configFilePrefix + env.String())
	v.SetConfigType(configType)

	// Add multiple config paths
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		var nf viper.ConfigFileNotFoundError
		// file found, but error parsing/access permission
		if !errors.As(err, &nf) {
			return nil, fmt.Errorf("error reading config file '%s': %w", configFilePrefix+env.String()+"."+configType, err)
		}
		log.Printf("config file '%s' not found, using defaults and environment variables", configFilePrefix+env.String()+"."+configType)
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

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("logging.level", LoggingLevelInfo.String())
	v.SetDefault("logging.format", LoggingFormatJSON.String())
}

func (c *Config) Validate() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	return nil
}
