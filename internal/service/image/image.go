package image

import (
	"image"
	"strings"

	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
)

type service struct{}

func New() *service {
	return &service{}
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
	gocv.Threshold(filter, &threshold, 210, 255, gocv.ThresholdBinary)

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

	client.SetLanguage("WenQuanYiMicroHei", "ChromosomeHeavy", "DenneThreedee")
	client.SetPageSegMode(gosseract.PSM_AUTO)

	text, err := client.Text()
	if err != nil {
		return "", err
	}

	text = strings.ReplaceAll(strings.ReplaceAll(text, "  ", " "), " ", "")
	text = strings.ToLower(text)

	return text, nil
}
