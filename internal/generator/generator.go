package generator

import (
	"math/rand"

	"github.com/mojocn/base64Captcha"
)

type service struct {
	driver base64Captcha.Driver
}

func New() *service {
	return &service{}
}

func (s *service) Image() (base64Captcha.Item, error) {
	driverString := &base64Captcha.DriverString{
		Length: rand.Intn(6-4) + 4,
		Width:  200,
		Height: 80,
		Source: "1234567890qwertyuioplkjhgfdsazxcvbnm",
		Fonts:  []string{"wqy-microhei.ttc"}, // "chromohv.ttf", "DENNEthree-dee.ttf"
	}

	s.driver = driverString.ConvertFonts()

	_, content, _ := s.driver.GenerateIdQuestionAnswer()
	item, err := s.driver.DrawCaptcha(content)
	if err != nil {
		return nil, err
	}

	return item, nil
}
