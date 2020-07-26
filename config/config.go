package config

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/hcl"
	"github.com/kelseyhightower/envconfig"
)

// ConfigFile name
const ConfigFile = "amadeus.conf"

// Config is an amadeus configuration struct.
type Config struct {
	Address     string `hcl:"address" envconfig:"AMADEUS_ADDRESS"`
	BaseURL     string `hcl:"base_url" envconfig:"AMADEUS_BASE_URL"`
	Title       string `hcl:"title" envconfig:"AMADEUS_TITLE"`
	Token       string `hcl:"token" envconfig:"AMADEUS_TOKEN"`
	ProjectPath string `hcl:"project_path" envconfig:"PROJECT_PATH"`
}

// GetConfig loads configuration from environment variables or from file.
func GetConfig(useEnv bool) (cfg *Config, err error) {
	if useEnv {
		cfg, err = ReadEnv()
	} else {
		cfg, err = ReadFile(ConfigFile)
	}
	return cfg, err
}

// ReadFile reads an amadeus config from file.
func ReadFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	cfg := &Config{}
	err = hcl.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal hcl: %v", err)
	}

	prepare(cfg)
	return cfg, nil
}

// ReadEnv reads an amadeus config from environment variables.
func ReadEnv() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %v", err)
	}
	prepare(cfg)
	return cfg, nil
}

func prepare(cfg *Config) {
	cfg.BaseURL = strings.TrimSuffix(cfg.BaseURL, "/")
}

// Init generates an initial config string.
func Init() (string, error) {
	buf := new(bytes.Buffer)
	err := tpl.Execute(buf, map[string]interface{}{
		"token_secret": GenKeyHex(32),
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenKeyHex generates a crypto-random key with byte length byteLen
// and hex-encodes it to a string.
func GenKeyHex(byteLen int) string {
	bytes := make([]byte, byteLen)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

var tpl = template.Must(template.New("initial-config").Parse(strings.TrimSpace(`
address 	= "0.0.0.0:8080"
base_url 	= "https://amadeus.com"
title    	= "amadeus"
token 	 	= "{{.token_secret}}"
`)))
