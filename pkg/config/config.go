package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

var (
	config     *Config
	once       sync.Once
	configPath = "config.json"
)

// LoadConfig loads the configuration from config.json in the current working directory.
// If the file does not exist, it creates one with default values.
func LoadConfig() (*Config, error) {
	var err error
	once.Do(func() {
		cwd, _ := os.Getwd()
		fullPath := filepath.Join(cwd, configPath)
		file, e := os.Open(fullPath)
		if e != nil {
			// File does not exist, create with default config
			config = DefaultConfig()
			data, _ := config.ToJSON()
			_ = os.WriteFile(fullPath, data, 0644)
			return
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		cfg := DefaultConfig()
		if e := decoder.Decode(cfg); e != nil {
			config = DefaultConfig()
			err = e
			return
		}
		config = cfg
	})
	return config, err
}

// GetConfig returns the loaded configuration, loading it if necessary.
func GetConfig() *Config {
	cfg, _ := LoadConfig()
	return cfg
}