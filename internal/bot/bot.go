package bot

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/LugaMuga/UOFDBot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/proxy"
)

var Bot *tgbotapi.BotAPI

func createHTTPClient() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}

	pCfg := config.Config.Proxy

	switch pCfg.Type {
	case "http":
		proxyURL, _ := url.Parse("http://" + pCfg.Address)
		if pCfg.User != "" {
			proxyURL.User = url.UserPassword(pCfg.User, pCfg.Password)
		}
		transport.Proxy = http.ProxyURL(proxyURL)

	case "socks5":
		var auth *proxy.Auth
		if pCfg.User != "" {
			auth = &proxy.Auth{
				User:     pCfg.User,
				Password: pCfg.Password,
			}
		}
		dialer, err := proxy.SOCKS5("tcp", pCfg.Address, auth, proxy.Direct)
		if err != nil {
			panic("Failed to create SOCKS5 dialer: " + err.Error())
		}
		contextDialer := dialer.(proxy.ContextDialer)
		transport.DialContext = contextDialer.DialContext

	default:
		transport.Proxy = http.ProxyFromEnvironment
	}

	return &http.Client{
		Transport: transport,
		Timeout:   time.Duration(config.Config.BotTimeout) * time.Second,
	}
}

func RegisterBot() {
	var err error
	client := createHTTPClient()
	Bot, err = tgbotapi.NewBotAPIWithClient(config.Config.BotToken, tgbotapi.APIEndpoint, client)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized on account %s", Bot.Self.UserName)
	Bot.Request(tgbotapi.DeleteWebhookConfig{
		DropPendingUpdates: true,
	})
}

func SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, _ = Bot.Send(msg)
}
