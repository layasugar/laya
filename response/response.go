package response

import (
	"github.com/LaYa-op/laya/i18n"
	"strconv"
	"strings"
)

// 现做如下约定
// 1. 成功返回--{1."success",{},{}} code必须为1
// 2. 失败返回--{0,"系统发生错误！",{},{}} code 必须为0,code=0前端按照msg进行提示

type Response struct {
	RespData
	i18n.I18ner
}

type RespData struct {
	Code     int
	Msg      string
	Data     interface{}
	WithData interface{}
	Page     *PageRes `json:"Page,omitempty"`
}

type PageRes struct {
	CurPage int // 当前页
	Size    int // 每页条数
	Total   int // 总条数
}

// Get response information
// al is header's [Accept-Language]
func (resp *Response) GetResponse(params map[string]interface{}, al string) interface{} {
	for name, value := range params {
		if !strings.HasPrefix(name, "$.") {
			continue
		}
		lastOne := strings.Split(name, ".")[len(strings.Split(name, "."))-1]
		switch lastOne {
		case "code":
			resp.RespData.Code = value.(int)
		case "response":
			resp.RespData = value.(RespData)
		}
		resp.Msg = resp.I18ner.GetMessage(al, strconv.Itoa(resp.Code))
	}
	return resp
}
