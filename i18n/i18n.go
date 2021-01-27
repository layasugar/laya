package i18n

import (
	"github.com/BurntSushi/toml"
	"github.com/LaYa-op/laya/config"
	i "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io/ioutil"
)

var I18n = &I18ner{}

// I18ner Internationalization support
type I18ner struct {
	Bundle *i.Bundle
}

// getMessage Gets the restfulApi to return value translation information
// al = accept_language
func (i18n *I18ner) GetMessage(al string, msg string) string {
	lang := i18n.getLang(al)
	loc := i.NewLocalizer(i18n.Bundle, lang)

	return loc.MustLocalize(&i.LocalizeConfig{
		MessageID: msg,
		DefaultMessage: &i.Message{
			ID:    msg,
			Other: "The translation could not be found.",
		},
	})
}

// translate Get general translation information
func (i18n *I18ner) Translate(lang string, msg string) string {
	loc := i.NewLocalizer(i18n.Bundle, lang)

	return loc.MustLocalize(&i.LocalizeConfig{
		MessageID: msg,
		DefaultMessage: &i.Message{
			ID:    msg,
			Other: "The translation could not be found.",
		},
	})
}

// initialize i18n
func Init() {
	c := config.GetI18nConfig()
	if c.Open {
		I18n.Bundle = i.NewBundle(language.English)
		I18n.Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		err := I18n.LoadAllFile(c.Path)
		if err != nil {
			panic(err)
		}
	}
}

// Load the file
func (i18n *I18ner) LoadAllFile(pathname string) error {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			_ = i18n.LoadAllFile(pathname + fi.Name() + "\\")
		} else {
			_, err := i18n.Bundle.LoadMessageFile(pathname + fi.Name())
			if err != nil {
				return err
			}
		}
	}
	return err
}

// get language
func (i18n *I18ner) getLang(lang string) string {
	c := config.GetI18nConfig()
	if lang == "" {
		if c.Open {
			lang = c.DefaultLang
		} else {
			lang = language.English.String()
		}
	}

	return string([]rune(lang)[0:2])
}
