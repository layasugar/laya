package i18n

import (
	"github.com/BurntSushi/toml"
	i "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io/ioutil"
)

// I18n Internationalization support
type I18n struct {
	Bundle *i.Bundle
	Conf   struct {
		Open        bool   `json:"open"`
		DefaultLang string `json:"defaultLang"`
	}
}

// getMessage Gets the restfulApi to return value translation information
func (i18n *I18n) GetMessage(al string, msg string) string {
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
func (i18n *I18n) Translate(lang string, msg string) string {
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
func (i18n *I18n) InitLang() {
	if i18n.Conf.Open {
		i18n.Bundle = i.NewBundle(language.English)
		i18n.Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		err := i18n.LoadAllFile("./conf/lang/")
		if err != nil {
			panic(err)
		}
	}
}

// Load the file
func (i18n *I18n) LoadAllFile(pathname string) error {
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
func (i18n *I18n) getLang(lang string) string {
	if lang == "" {
		if i18n.Conf.Open {
			lang = i18n.Conf.DefaultLang
		} else {
			lang = language.English.String()
		}
	}

	return string([]rune(lang)[0:2])
}
