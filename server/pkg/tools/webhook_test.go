package tools

import (
	"testing"
)

func TestDingTalk(t *testing.T) {
	ding := DingTalk{
		MsgType: "text",
	}

	if err := ding.SendMsg(`业务受到法律
			http://www.baidu.com`,
		"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=8eed3dca-27d5-46bf-ba4d-67c8ef69054a",
		"17377042972", "17377042972", "17377042972",
	); err != nil {
		t.Error("异常", err.Error())
	}
}

func TestTencent(t *testing.T) {
	ten := Tencent{
		MsgType:   "image",
		ImgBase64: "",
		ImgMd5:    "",
	}
	if err := ten.SendMsg(ten.ImgBase64, ten.ImgMd5,
		"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=8eed3dca-27d5-46bf-ba4d-67c8ef69054a",
	); err != nil {
		t.Error("异常", err.Error())
	}
}
