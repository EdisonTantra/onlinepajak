package config

import "time"

// Config represents the data structure for configuration in Apps.
type Config struct {
	Env        string      `json:"env"`
	Service    *Service    `json:"service"`
	ServerHTTP *ServerHTTP `json:"http"`
	Postgres   *Postgresql `json:"postgresql"`
	Telemetry  *Telemetry  `json:"telemetry"`
	DJPClient  *DJPClient  `json:"djp"`
}

// Service represents the data structure for service general info in Apps.
type Service struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Version  string `json:"version"`
}

// Logger represents the data structure for logger configuration in Apps.
type Logger struct {
	Level   string `json:"level"`
	Encoder string `json:"encoder"`
}

// ServerHTTP represents the data structure for HTTP server configuration.
type ServerHTTP struct {
	Address           string        `json:"address"`
	ReadTimeout       time.Duration `json:"readTimeout"`
	WriteTimeout      time.Duration `json:"writeTimeout"`
	ReadHeaderTimeout time.Duration `json:"readHeaderTimeout"`
}

// Postgresql represents the connection for Postgresql Database configuration.
type Postgresql struct {
	Host        string        `json:"host"`
	Port        int           `json:"port"`
	Database    string        `json:"db"`
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	SSLMode     bool          `json:"sslMode"`
	MaxOpenConn int           `json:"maxOpenConn"`
	MaxIdleConn int           `json:"maxIdleConn"`
	MaxIdleTime time.Duration `json:"maxIdleTime"`
}

// Telemetry represents the connection for OpenTelemetry configuration.
type Telemetry struct {
	CollectorURL string `json:"collectorURL"`
	SecureMode   bool   `json:"secureMode"`
}

// DJPClient represents the connection for DJPClient configuration.
type DJPClient struct {
	BaseURL string `json:"baseURL"`
}
