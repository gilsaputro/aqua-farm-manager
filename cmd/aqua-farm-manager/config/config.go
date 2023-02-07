package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config struct to hold the configuration data for server
type Config struct {
	Vault       Vault    `yaml:"vault"`
	Redis       Redis    `yaml:"redis"`
	Postgres    Postgres `yaml:"postgres"`
	ES          ES       `yaml:"es"`
	FarmHandler Handler  `yaml:"farm_handler"`
	StatHandler Handler  `yaml:"stat_handler"`
}

// Vault struct to hold the configuration data for vault
type Vault struct {
	Token     string `yaml:"token"`
	VaultHost string `yaml:"vault_host"`
}

// Postgres struct to hold the configuration data for postgres
type Postgres struct {
	Config string `yaml:"postgres_config"`
}

// ES struct to hold the configuration data for ES
type ES struct {
	Host string `yaml:"host"`
}

// Redis struct to hold the configuration data for redis
type Redis struct {
	RedisPassword    string `yaml:"redis_password"`
	RedisHost        string `yaml:"redis_host"`
	MaxIdleInSec     int64  `yaml:"max_idle_in_sec"`
	IdleTimeoutInSec int64  `yaml:"idle_timeout_in_sec"`
}

// Handler struct to hold the configuration data for handler
type Handler struct {
	TimeoutInSec       int `yaml:"timeout_in_sec"`
	BackupTimeInMinute int `yaml:"backup_time_in_minute"`
}

// GetConfig is func to load config and replace it by secret value
func GetConfig(values map[string]string) (Config, error) {
	var cfg Config
	// Read the YAML file into a byte slice
	configPath := filepath.Join("config", "config.yaml")
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}
	// Replace yaml value with secret
	for k, v := range values {
		configData = []byte(strings.Replace(string(configData), fmt.Sprintf("<%v>", k), v, -1))
	}
	// Unmarshal the YAML into a Config struct
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
