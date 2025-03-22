package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
	Database struct {
		Provider string `json:"provider"`
		User     string `json:"user"`
		Password string `json:"password"`
		Dbname   string `json:"dbname"`
		Host     string `json:"host"`
		SSLmode  bool   `json:"sslmode"`
	} `json:"database"`
	Storage struct {
		StoreMasterPath string `json:"store_master_path"`
	} `json:"storage"`
	Processing struct {
		NumWorkers int `json:"num_workers"`
	} `json:"processing"`
}

// LoadConfig loads configuration from a JSON file and sets default values if not provided
func LoadConfig(filePath string) (*Config, error) {
	config := &Config{}

	// Set default values
	config.Server.Port = ":8080"
	config.Storage.StoreMasterPath = "./store_master.csv"
	config.Processing.NumWorkers = 4

	file, err := os.Open(filePath)
	if err != nil {
		// If the file doesn't exist, return the config with default values
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
