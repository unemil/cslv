package generator

import (
	"math/rand"

	"github.com/mojocn/base64Captcha"
)

type service struct {
	driver base64Captcha.Driver

	driverString *base64Captcha.DriverString
}

func New() *service {
	return &service{
		driverString: &base64Captcha.DriverString{
			Width:  200,
			Height: 80,
			Source: "1234567890qwertyuioplkjhgfdsazxcvbnm",
			Fonts:  []string{"wqy-microhei.ttc", "chromohv.ttf", "DENNEthree-dee.ttf"},
		},
	}
}

func (s *service) Image() (base64Captcha.Item, error) {
	s.driver = s.driverString.ConvertFonts()
	s.driverString.Length = rand.Intn(6-4) + 4
	s.driverString.NoiseCount = s.driverString.Length / 2

	_, content, _ := s.driver.GenerateIdQuestionAnswer()
	item, err := s.driver.DrawCaptcha(content)
	if err != nil {
		return nil, err
	}

	return item, nil
}
