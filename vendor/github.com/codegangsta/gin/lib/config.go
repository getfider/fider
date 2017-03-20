package gin

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Laddr    string `json:"laddr"`
	Port     int    `json:"port"`
	ProxyTo  string `json:"proxy_to"`
	KeyFile  string `json:"key_file"`
	CertFile string `json:"cert_file"`
}

func LoadConfig(path string) (*Config, error) {
	configFile, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("Unable to read configuration file %s", path)
	}

	config := new(Config)

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse configuration file %s", path)
	}

	return config, nil
}
