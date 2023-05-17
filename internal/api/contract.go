package api

import "cslv/internal/generator/captcha"

type Service interface {
	Solve(file []byte) (string, error)
}

type Generator interface {
	Image() (captcha.Captcha, error)
}
