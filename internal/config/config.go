package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	LocalEnv = "local"
	ProdEnv  = "prod"
)

// Config is a structure that holds the application configuration.
type Config struct {
	Env         string      `json:"env"`
	DB          DB          `json:"db"`
	Cache       Cache       `json:"cache"`
	JwtSettings JwtSettings `json:"jwt_settings"`
	HTTPServer  HTTPServer  `json:"http_server"`
}

func (c Config) String() string {
	return fmt.Sprintf("{Env: %s, DB: %s, Cache: %s, JwtSettings: %s, HTTPServer: %s}",
		c.Env, c.DB, c.Cache, c.JwtSettings, c.HTTPServer)
}

// MustLoad reads the configuration from the file specified from the command line 'config' argument
// or from the CONFIG_PATH environment variable.
// The command line argument has a higher priority than the environment variable.
func MustLoad(args []string, getenv func(string) (string, bool)) *Config {
	configPath, err := getConfigPath(args, getenv)
	if err != nil {
		panic(err)
	}

	cfg := new(Config)
	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read config: %s", err)
	}
	err = json.Unmarshal(content, cfg)
	if err != nil {
		log.Fatalf("failed to unmarshal config: %s", err)
	}

	return cfg
}

// getConfigPath reads the config path either from command line
// or from the CONFIG_PATH environment variable.
// The command line argument has a higher priority than the environment variable.
func getConfigPath(args []string, getenv func(string) (string, bool)) (string, error) {
	var path string
	flagSet := flag.NewFlagSet("config-path", flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)

	pathFlag := flagSet.String("config", "", "")

	_ = flagSet.Parse(args)

	if pathFlag == nil || *pathFlag == "" {
		pathEnv, ok := getenv("CONFIG_PATH")
		if !ok {
			return "", errors.New("config path not provided")
		}
		path = pathEnv
	} else {
		path = *pathFlag
	}

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("specified config file does not exist: %s", path)
	}

	return path, nil
}
