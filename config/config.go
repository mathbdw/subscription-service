package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
)

// Build information -ldflags .
const (
	version    string = "dev"
	commitHash string = "-"
)

// Database - contains all parameters database connection.
type Database struct {
	Host            string        `yaml:"host" env:"PG_HOST,required"`
	Port            uint16        `yaml:"port" env:"PG_PORT,required"`
	User            string        `yaml:"user" env:"PG_USER,required"`
	Password        string        `yaml:"password" env:"PG_PASSWORD,required"`
	Migrations      string        `yaml:"migrations"`
	Name            string        `yaml:"name"`
	SslMode         string        `yaml:"sslmode"`
	Driver          string        `yaml:"driver"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`
	MaxIdleConns    int           `yaml:"maxIdleConns"`
	ConnMaxIdleTime time.Duration `yaml:"connMaxIdleTime"`
	ConnMaxLifetime time.Duration `yaml:"connMaxLifetime"`
}

// Rest - contains parameter rest json connection.
type Rest struct {
	Host            string        `yaml:"host" env:"REST_HOST,required"`
	Port            uint16        `yaml:"port"`
	Prefork         bool          `yaml:"prefork"`
	ReadTimeout     time.Duration `yaml:"readTimeout"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
	Swagger         bool          `yaml:"swagger"`
}

// Project - contains all parameters project information.
type Project struct {
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Version     string `yaml:"version"`
	CommitHash  string
	Debug       bool `yaml:"debug"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	Project  Project  `yaml:"project"`
	Rest     Rest     `yaml:"rest"`
	Database Database `yaml:"database"`
}

// ReadConfigYML - read configurations from file and init instance Config.
func ReadConfigYML(filePath string) (*Config, error) {
	var cfg *Config
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	cfg.Project.Version = version
	cfg.Project.CommitHash = commitHash

	return cfg, nil
}
