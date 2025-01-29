package config

import (
	"time"
)

type Config struct {
	WorkerPoolSize int           `yaml:"worker_pool_size"`
	HTTPTimeout    time.Duration `yaml:"http_timeout"`
	CacheEnabled   bool          `yaml:"cache_enabled"`
	CacheTTL       time.Duration `yaml:"cache_ttl"`
}

func LoadConfig(path string) (*Config, error) {
	// Возвращаем конфиг по умолчанию, пока не реализована загрузка из файла
	return &Config{
		WorkerPoolSize: 4,
		HTTPTimeout:    10 * time.Second,
		CacheEnabled:   true,
		CacheTTL:       5 * time.Minute,
	}, nil
}
