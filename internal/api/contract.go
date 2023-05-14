package api

import "github.com/mojocn/base64Captcha"

type Service interface {
	Solve(file []byte) string
}

type Generator interface {
	Image() (base64Captcha.Item, error)
}
