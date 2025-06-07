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

type SMTPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	MongoDB     MongoDBConfig    `json:"mongodb"`
	HTTPServer  HTTPServerConfig `json:"http_server"`
	AdminSecret string           `json:"admin_secret"`
	Name        string           `json:"name"`
	SMTP        SMTPConfig       `json:"smtp"`
	JWTSecret   string           `json:"jwt_secret"`
}

func DefaultConfig() *Config {
	return &Config{
		MongoDB: MongoDBConfig{
			URI:      "mongodb://user:password@localhost:27017/",
			Database: "goauthx",
		},
		HTTPServer: HTTPServerConfig{
			Port:        5001,
			EnableSSL:   false,
			SSLCertFile: "",
			SSLKeyFile:  "",
		},
		AdminSecret: "your_admin_secret",
		Name:        "GoAuthX",
		SMTP: SMTPConfig{
			Host:     "smtp.example.com",
			Port:     465,
			Username: "user@example.com",
			Password: "your_smtp_password",
		},
		JWTSecret: "your_jwt_secret",
	}
}

// Marshal default config to JSON (for reference)
func (c *Config) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}
