package locale

import (
	"github.com/LugaMuga/UOFDBot/internal/config"
	"github.com/goccy/go-yaml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"os"
	"path/filepath"
)

const russianLangFile = `ru.yml`
const englishLangFile = `en.yml`

var i18nBundle *i18n.Bundle
var DefaultLang language.Tag
var ruLocalizer *i18n.Localizer
var enLocalizer *i18n.Localizer

func LoadI18N() {
	var err error
	DefaultLang, err = language.Parse(config.Config.BotDefaultLanguage)
	if err != nil {
		log.Fatal(err)
	}
	i18nBundle = i18n.NewBundle(DefaultLang)
	i18nBundle.RegisterUnmarshalFunc("yml", yaml.Unmarshal)
	langDirPath := os.Getenv("UOFD_LANG_DIR_PATH")
	if langDirPath == "" {
		workDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		langDirPath = filepath.Dir(workDir) + "/UOFDBot/lang"
	}
	i18nBundle.MustLoadMessageFile(langDirPath + "/" + russianLangFile)
	i18nBundle.MustLoadMessageFile(langDirPath + "/" + englishLangFile)
	log.Println(`Locales successfully loaded!`)

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

func Loc(languageTag language.Tag, messageId string, args ...interface{}) string {
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

func LocPlural(languageTag language.Tag, messageId string, pluralVal interface{}, args ...interface{}) string {
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
