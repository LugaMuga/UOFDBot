package main

import (
	"github.com/goccy/go-yaml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
)

const russianLangFile = `./lang/ru.yml`
const englishLangFile = `./lang/en.yml`

var i18nBundle *i18n.Bundle
var defaultLang language.Tag
var ruLocalizer *i18n.Localizer
var enLocalizer *i18n.Localizer

func LoadI18N() {
	var err error
	defaultLang, err = language.Parse(config.BotDefaultLanguage)
	if err != nil {
		log.Fatal(err)
	}
	i18nBundle = i18n.NewBundle(defaultLang)
	i18nBundle.RegisterUnmarshalFunc("yml", yaml.Unmarshal)
	i18nBundle.MustLoadMessageFile(russianLangFile)
	i18nBundle.MustLoadMessageFile(englishLangFile)

	ruLocalizer = i18n.NewLocalizer(i18nBundle, language.Russian.String())
	enLocalizer = i18n.NewLocalizer(i18nBundle, language.English.String())
}

func getLocalizer(languageTag language.Tag) *i18n.Localizer {
	switch languageTag.String() {
	case language.Russian.String():
		return ruLocalizer
	case language.English.String():
		return enLocalizer
	default:
		return nil
	}
}

func loc(languageTag language.Tag, messageId string, args ...interface{}) string {
	localizer := getLocalizer(languageTag)
	data := make(map[string]interface{})
	for index, arg := range args {
		data[toCharStr(index)] = arg
	}
	msg := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageId,
		TemplateData: data,
	})
	return msg
}

func locPlural(languageTag language.Tag, messageId string, pluralVal interface{}, args ...interface{}) string {
	localizer := getLocalizer(languageTag)
	data := make(map[string]interface{})
	for index, arg := range args {
		data[toCharStr(index)] = arg
	}
	msg := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageId,
		TemplateData: data,
		PluralCount:  pluralVal,
	})
	return msg
}

func toCharStr(i int) string {
	return string('a' + i)
}
