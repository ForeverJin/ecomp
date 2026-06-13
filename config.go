package main

import (
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

// Config holds MySQL connection settings, persisted as JSON.
type Config struct {
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
}

func configPath(dataDir string) string {
	return filepath.Join(dataDir, "config.json")
}

// LoadConfig reads config.json from dataDir. Returns nil if file doesn't exist.
func LoadConfig(dataDir string) (*Config, error) {
	data, err := os.ReadFile(configPath(dataDir))
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveConfig writes config.json to dataDir.
func SaveConfig(dataDir string, cfg *Config) error {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(dataDir), data, 0644)
}

// DSN returns the MySQL data source name.
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=utf8mb4",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// TestConnection tries to connect to MySQL with the given parameters.
func TestConnection(host string, port int, user, password, dbname string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
		user, password, host, port, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Try a longer timeout for initial connection test
	db.SetMaxOpenConns(1)
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	_ = tlsConfig // not used but avoids import warning

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}
	return nil
}

// IsConfigured returns true if the essential fields are set.
func (c *Config) IsConfigured() bool {
	return c.DBHost != ""
}
