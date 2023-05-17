package captcha

import "github.com/mojocn/base64Captcha"

type Captcha struct {
	ID     string
	Item   base64Captcha.Item
	Answer string
}
