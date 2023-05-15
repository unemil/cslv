package generator

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/mojocn/base64Captcha"
)

type service struct {
	driver base64Captcha.Driver
}

func New() *service {
	driverString := &base64Captcha.DriverString{
		Length: rand.Intn(6-4) + 4,
		Width:  200,
		Height: 80,
		Source: "1234567890qwertyuioplkjhgfdsazxcvbnm",
		Fonts:  []string{"wqy-microhei.ttc", "DENNEthree-dee.ttf", "chromohv.ttf"}, // "chromohv.ttf", "DENNEthree-dee.ttf"
		// ShowLineOptions: base64Captcha.OptionShowSlimeLine,
	}

	return &service{
		driver: driverString.ConvertFonts(),
	}
}

func (s *service) Image() (base64Captcha.Item, error) {
	_, content, _ := s.driver.GenerateIdQuestionAnswer()
	item, err := s.driver.DrawCaptcha(content)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *service) Dataset() error {
	annotations, err := os.Create("annotations.txt")
	if err != nil {
		return err
	}
	defer annotations.Close()

	for i := 0; i < 10000; i++ {
		id, content, answer := s.driver.GenerateIdQuestionAnswer()
		image, err := s.driver.DrawCaptcha(content)
		if err != nil {
			return err
		}

		file, err := os.Create(fmt.Sprintf("images/%s.png", answer))
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := image.WriteTo(file); err != nil {
			return err
		}

		annotations.WriteString(fmt.Sprintf("images/%s.png %s\n", answer, answer))

		log.Printf("%s: %s\n", id, answer)
	}

	return nil
}
