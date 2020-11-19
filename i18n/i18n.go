package i18n

import (
	"github.com/LaYa-op/laya"
	i "github.com/nicksnyder/go-i18n/v2/i18n"
)

// I18n Internationalization support
type I18n struct {
	lang string
}

// getMessage Gets the restfulApi to return value translation information
func (*I18n) getMessage(lang string, msg string) string {
	loc := i.NewLocalizer(laya.I18nBundle, lang)

	return loc.MustLocalize(&i.LocalizeConfig{
		MessageID: msg,
		DefaultMessage: &i.Message{
			ID:    msg,
			Other: "The translation could not be found.",
		},
	})
}

// translate Get general translation information
func (*I18n) translate(lang string, msg string) string {
	loc := i.NewLocalizer(laya.I18nBundle, lang)

	return loc.MustLocalize(&i.LocalizeConfig{
		MessageID: msg,
		DefaultMessage: &i.Message{
			ID:    msg,
			Other: "The translation could not be found.",
		},
	})
}
