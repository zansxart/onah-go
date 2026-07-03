package config

import (
	"encoding/json"
	"os"
)

// Config represents the system settings loaded from config.json
type Config struct {
	OwnerNumber        string            `json:"owner_number"`
	OwnerName          string            `json:"owner_name"`
	BotName            string            `json:"bot_name"`
	Prefixes           []string          `json:"prefixes"`
	DatabasePath       string            `json:"database_path"`
	LimitDefault       int               `json:"limit_default"`
	PairingCodeEnabled bool              `json:"pairing_code_enabled"`
	PairingNumber      string            `json:"pairing_number"`
	ApiKeys            map[string]string `json:"api_keys"`
	Messages           map[string]string `json:"messages"`
}

// ActiveConfig is the globally accessible loaded configuration
var ActiveConfig Config

// LoadConfig parses configuration from a JSON file
func LoadConfig(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &ActiveConfig)
	if err != nil {
		return err
	}
	return nil
}
