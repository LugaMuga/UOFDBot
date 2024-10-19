package config

import (
	"github.com/goccy/go-yaml"
	"log"
	"os"
	"path/filepath"
)

const ConfigFilePath string = `/UOFDBot/configs/config.yml`
const ConnectionTypeChannel string = `CHANNEL`
const ConnectionTypeWebhook string = `WEBHOOK`

var Config struct {
	BotName               string `yaml:"bot_name"`
	BotToken              string `yaml:"bot_token"`
	BotTimeout            int    `yaml:"bot_timeout"`
	BotTimeLayout         string `yaml:"bot_time_layout"`
	BotDefaultLanguage    string `yaml:"bot_default_language"`
	BotResetMinPercentage int    `yaml:"bot_reset_min_percantage"`
	ConnectionType        string `yaml:"connection_type"`
	WebhookHost           string `yaml:"webhook_host"`
	WebhookPort           string `yaml:"webhook_port"`
	WebhookCertPath       string `yaml:"webhook_cert_path"`
	WebhookKeyPath        string `yaml:"webhook_key_path"`
	ServerAddress         string `yaml:"server_adress"`
}

func LoadConfig() {
	configPath := os.Getenv("UOFD_CONFIG_FILE_PATH")
	if configPath == "" {
		workDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		configPath = filepath.Dir(workDir) + ConfigFilePath
	}
	file, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(`Config successfully loaded!`)
}
