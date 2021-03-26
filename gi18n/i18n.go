package gi18n

import (
	"github.com/BurntSushi/toml"
	"github.com/layatips/laya/gconf"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io/ioutil"
)

var GI18n = &I18ner{}

// I18ner Internationalization support
type I18ner struct {
	Bundle *i18n.Bundle
}

// getMessage Gets the restfulApi to return value translation information
// al = accept_language
func (gi18n *I18ner) GetMessage(al string, msg string) string {
	lang := gi18n.getLang(al)
	loc := i18n.NewLocalizer(gi18n.Bundle, lang)

	return loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID: msg,
		DefaultMessage: &i18n.Message{
			ID:    msg,
			Other: "The translation could not be found.",
		},
	})
}

// translate Get general translation information
func (gi18n *I18ner) Translate(lang string, msg string) string {
	loc := i18n.NewLocalizer(gi18n.Bundle, lang)

	return loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID: msg,
		DefaultMessage: &i18n.Message{
			ID:    msg,
			Other: "The translation could not be found.",
		},
	})
}

// initialize gi18n
func Init() {
	c := gconf.GetI18nConfig()
	if c.Open {
		GI18n.Bundle = i18n.NewBundle(language.English)
		GI18n.Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		err := GI18n.LoadAllFile(c.Path)
		if err != nil {
			panic(err)
		}
	}
}

// Load the file
func (gi18n *I18ner) LoadAllFile(pathname string) error {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			_ = gi18n.LoadAllFile(pathname + fi.Name() + "\\")
		} else {
			_, err := gi18n.Bundle.LoadMessageFile(pathname + fi.Name())
			if err != nil {
				return err
			}
		}
	}
	return err
}

// get language
func (gi18n *I18ner) getLang(lang string) string {
	c := gconf.GetI18nConfig()
	if lang == "" {
		if c.Open {
			lang = c.DefaultLang
		} else {
			lang = language.English.String()
		}
	}

	return string([]rune(lang)[0:2])
}
