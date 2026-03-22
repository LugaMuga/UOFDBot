package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

const ConfigFilePath string = `/UOFDBot/configs/config.yml`
const ConnectionTypeChannel string = `CHANNEL`
const ConnectionTypeWebhook string = `WEBHOOK`

type ProxyConfig struct {
	Type     string `yaml:"type"` // "http", "socks5", ""
	Address  string `yaml:"address"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

var Config struct {
	BotName               string      `yaml:"bot_name"`
	BotToken              string      `yaml:"bot_token"`
	BotTimeout            int         `yaml:"bot_timeout"`
	BotTimeLayout         string      `yaml:"bot_time_layout"`
	BotDefaultLanguage    string      `yaml:"bot_default_language"`
	BotResetMinPercentage int         `yaml:"bot_reset_min_percantage"`
	ConnectionType        string      `yaml:"connection_type"`
	WebhookHost           string      `yaml:"webhook_host"`
	WebhookPort           string      `yaml:"webhook_port"`
	WebhookCertPath       string      `yaml:"webhook_cert_path"`
	WebhookKeyPath        string      `yaml:"webhook_key_path"`
	ServerAddress         string      `yaml:"server_adress"`
	Proxy                 ProxyConfig `yaml:"proxy"`
}

func LoadConfig() {
	configPath := os.Getenv("UOFD_CONFIG_FILE_PATH")
	if configPath == "" {
		workDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		configPath = filepath.Join(filepath.Dir(workDir), ConfigFilePath)
	}
	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		log.Fatalf("Error parsing config yaml: %v", err)
	}
	log.Println(`Config successfully loaded!`)
}
