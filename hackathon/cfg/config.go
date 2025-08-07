package cfg

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"strconv"
	"strings"
)

type Config struct {
	AppEnv      string `json:"APP_ENV"`
	EnableCache string `json:"CACHE"`
	HttpPort    string `json:"HTTP_PORT"`
	HttpHost    string `json:"HTTP_HOST"`
	JWTSecret   string `json:"JWT_SECRET"`
	MaxFileSize string `json:"MAX_FILE_SIZE"`

	DB
	Cache
}

func (c Config) GetMaxFileSize() (int64, error) {
	sizeStr := strings.TrimSpace(strings.ToUpper(c.MaxFileSize))
	if size, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
		return size, nil
	}

	var multiplier int64 = 1
	var numStr string

	switch {
	case strings.HasSuffix(sizeStr, "GB"):
		multiplier = 1024 * 1024 * 1024
		numStr = strings.TrimSuffix(sizeStr, "GB")
	case strings.HasSuffix(sizeStr, "MB"):
		multiplier = 1024 * 1024
		numStr = strings.TrimSuffix(sizeStr, "MB")
	case strings.HasSuffix(sizeStr, "KB"):
		multiplier = 1024
		numStr = strings.TrimSuffix(sizeStr, "KB")
	case strings.HasSuffix(sizeStr, "B"):
		multiplier = 1
		numStr = strings.TrimSuffix(sizeStr, "B")
	default:
		return 0, fmt.Errorf("invalid format, supported units: B, KB, MB, GB")
	}

	// Parse the numeric part (allow decimals for convenience)
	if strings.Contains(numStr, ".") {
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %s", numStr)
		}
		return int64(num * float64(multiplier)), nil
	} else {
		num, err := strconv.ParseInt(numStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number: %s", numStr)
		}
		return num * multiplier, nil
	}
}

type Cache struct {
	CacheHost string `json:"CACHE_HOST"`
	CachePort string `json:"CACHE_PORT"`
	CachePass string `json:"CACHE_PASS"`
	CacheDB   string `json:"CACHE_DB"`
}

type DB struct {
	DBDriver string `json:"DB_DRIVER"`
	DBHost   string `json:"DB_HOST"`
	DBPort   string `json:"DB_PORT"`
	DBUser   string `json:"DB_USER"`
	DBPass   string `json:"DB_PASS"`
	DBName   string `json:"DB_DATABASE"`
}

func LoadConfig() Config {
	var config Config
	data, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(jsonStr, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
