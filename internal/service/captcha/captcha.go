package captcha

import (
	"context"
	"cslv/internal/model"
	"image"
	"math/rand"
	"regexp"
	"strings"
	"sync"

	"github.com/mojocn/base64Captcha"
	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
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
			Fonts:  []string{"wqy-microhei.ttc", "DENNEthree-dee.ttf"},
		},
	}
}

func (s *service) Generate() (model.Captcha, error) {
	s.driver = s.driverString.ConvertFonts()
	s.driverString.Length = rand.Intn(2) + 4
	s.driverString.NoiseCount = s.driverString.Length / 2

	id, content, answer := s.driver.GenerateIdQuestionAnswer()
	item, err := s.driver.DrawCaptcha(content)
	if err != nil {
		return model.Captcha{}, err
	}

	return model.Captcha{
		ID:     id,
		Item:   item,
		Answer: answer,
	}, nil
}

func (s *service) Solve(file []byte) (string, error) {
	img, err := gocv.IMDecode(file, gocv.IMReadColor)
	if err != nil {
		return "", err
	}

	// Preprocessing
	resize := gocv.NewMat()
	gocv.Resize(img, &resize, image.Pt(img.Cols()*2, img.Rows()*2), 0, 0, gocv.InterpolationLinear)

	grayscale := gocv.NewMat()
	gocv.CvtColor(resize, &grayscale, gocv.ColorBGRToGray)

	filter := gocv.NewMat()
	gocv.BilateralFilter(grayscale, &filter, 5, 15, 15)

	// Segmentation
	threshold := gocv.NewMat()
	gocv.Threshold(filter, &threshold, 200, 255, gocv.ThresholdBinary)

	// Classification
	client := gosseract.NewClient()
	defer client.Close()

	buffer, err := gocv.IMEncode(gocv.PNGFileExt, threshold)
	if err != nil {
		return "", err
	}
	defer buffer.Close()

	if err := client.SetImageFromBytes(buffer.GetBytes()); err != nil {
		return "", err
	}

	client.SetLanguage("WenQuanYiMicroHei", "DenneThreedee")

	text, err := client.Text()
	if err != nil {
		return "", err
	}

	text = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(text, "")
	text = strings.ToLower(text)

	return text, nil
}

func (s *service) Analyze(ctx context.Context, info []model.AnalyzeInfo) ([]model.Analysis, float32, error) {
	var (
		accuracy float32 = 0
		analysis         = make([]model.Analysis, 0, 1)

		wg = new(sync.WaitGroup)
		mu = new(sync.Mutex)
	)

	for _, i := range info {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i model.AnalyzeInfo) error {
			defer wg.Done()

			solution, err := s.Solve(i.Image)
			if err != nil {
				return err
			}

			contains := 0
			copy := solution
			for j := range i.Data.Answer {
				if ok := strings.Contains(copy, string(i.Data.Answer[j])); ok {
					contains++
					copy = strings.Replace(copy, string(i.Data.Answer[j]), "", 1)
				}
			}

			var rate float32 = float32(contains) / float32(len(i.Data.Answer))

			mu.Lock()
			accuracy += rate
			mu.Unlock()

			analysis = append(analysis, model.Analysis{
				ID:       i.Data.ID,
				Answer:   i.Data.Answer,
				Solution: solution,
				Rate:     rate * 100,
			})

			return nil
		}(wg, i)
	}
	wg.Wait()

	accuracy = (accuracy / float32(len(info))) * 100

	return analysis, accuracy, nil
}
