package captcha

import (
	"math/rand"

	"github.com/mojocn/base64Captcha"
)

type generator struct {
	driver base64Captcha.Driver

	driverString *base64Captcha.DriverString
}

func New() *generator {
	return &generator{
		driverString: &base64Captcha.DriverString{
			Width:  200,
			Height: 80,
			Source: "1234567890qwertyuioplkjhgfdsazxcvbnm",
			Fonts:  []string{"wqy-microhei.ttc", "chromohv.ttf", "DENNEthree-dee.ttf"},
		},
	}
}

func (g *generator) Image() (Captcha, error) {
	g.driver = g.driverString.ConvertFonts()
	g.driverString.Length = rand.Intn(2) + 4
	g.driverString.NoiseCount = g.driverString.Length / 2

	id, content, answer := g.driver.GenerateIdQuestionAnswer()
	item, err := g.driver.DrawCaptcha(content)
	if err != nil {
		return Captcha{}, err
	}

	return Captcha{
		ID:     id,
		Item:   item,
		Answer: answer,
	}, nil
}
