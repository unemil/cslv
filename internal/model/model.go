package model

import "github.com/mojocn/base64Captcha"

type Captcha struct {
	ID     string             `json:"id"`
	Item   base64Captcha.Item `json:"-"`
	Answer string             `json:"answer"`
}

type Analysis struct {
	ID       string  `json:"id"`
	Answer   string  `json:"answer"`
	Solution string  `json:"solution"`
	Rate     float32 `json:"rate"`
}

type AnalyzeInfo struct {
	Image []byte
	Data  Captcha
}
