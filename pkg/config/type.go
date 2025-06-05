package config

import "encoding/json"

type MongoDBConfig struct {
	URI      string `json:"uri"`
	Database string `json:"database"`
}

type HTTPServerConfig struct {
	Port        int    `json:"port"`
	EnableSSL   bool   `json:"enable_ssl"`
	SSLCertFile string `json:"ssl_cert_file"`
	SSLKeyFile  string `json:"ssl_key_file"`
}

type Config struct {
	MongoDB    MongoDBConfig    `json:"mongodb"`
	HTTPServer HTTPServerConfig `json:"http_server"`
	AdminSecret string          `json:"admin_secret"`
}

func DefaultConfig() *Config {
	return &Config{
		MongoDB: MongoDBConfig{
			URI:      "mongodb://root:root@localhost:27017/",
			Database: "goauthx",
		},
		HTTPServer: HTTPServerConfig{
			Port:        8080,
			EnableSSL:   false,
			SSLCertFile: "",
			SSLKeyFile:  "",
		},
		AdminSecret: "changeme_admin_secret",
	}
}

// Marshal default config to JSON (for reference)
func (c *Config) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}