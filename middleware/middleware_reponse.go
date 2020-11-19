package middleware

import (
	"github.com/LaYa-op/laya"
	"github.com/LaYa-op/laya/response"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"net/http"
	"strconv"
	"strings"
)

func (*Middleware) Response() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Written() {
			return
		}

		params := c.Keys
		if len(params) == 0 {
			return
		}

		lang := GetLang(c.GetHeader("Accept-Language"))
		resp := GetResponse(params, lang)
		c.JSON(http.StatusOK, resp)
	}
}

func GetLang(lang string) string {
	if lang == "" {
		if laya.I18nConf.Open {
			lang = laya.I18nConf.DefaultLang
		} else {
			lang = language.English.String()
		}
	}

	rs := []rune(lang)
	lang = string(rs[0:2])
	return lang
}

func GetResponse(params map[string]interface{}, lang string) interface{} {
	var resp response.Response
	for name, value := range params {
		if !strings.HasPrefix(name, "$.") {
			continue
		}
		lastOne := strings.Split(name, ".")[len(strings.Split(name, "."))-1]
		switch lastOne {
		case "code":
			resp.Code = value.(int)
		case "response":
			op, _ := value.(response.Response)
			resp = op
		}
		resp.Location = name
		resp.Msg = GetMessage(lang, strconv.Itoa(resp.Code))
	}
	return resp
}

func GetMessage(lang string, msg string) string {
	loc := i18n.NewLocalizer(laya.I18nBundle, lang)

	return loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID: msg,
		DefaultMessage: &i18n.Message{
			ID:    msg,
			Other: "The translation could not be found.",
		},
	})
}