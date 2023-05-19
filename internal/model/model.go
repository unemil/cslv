package model

import "github.com/mojocn/base64Captcha"

type Captcha struct {
	ID     string
	Item   base64Captcha.Item
	Answer string
}

type Analysis struct {
	ID       string  `json:"id"`
	Image    string  `json:"image"`
	Answer   string  `json:"answer"`
	Solution string  `json:"solution"`
	Rate     float32 `json:"rate"`
}
