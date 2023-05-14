package image

import (
	"image"
	"log"
	"strings"

	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
)

type service struct{}

func New() *service {
	return &service{}
}

func (s *service) Solve(file []byte) string {
	img, err := gocv.IMDecode(file, gocv.IMReadColor)
	if err != nil || img.Empty() {
		log.Println(err)
		return "gocv.IMDecode error"
	}

	// Preprocessing
	gocv.Resize(img, &img, image.Pt(img.Cols()*5, img.Rows()*5), 0, 0, gocv.InterpolationLinear)
	gocv.CvtColor(img, &img, gocv.ColorBGRToGray)
	gocv.GaussianBlur(img, &img, image.Pt(3, 3), 0, 0, gocv.BorderDefault)

	gocv.IMWrite("preprocessing.png", img)

	// Segmentation
	gocv.Threshold(img, &img, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)

	gocv.IMWrite("segmentation.png", img)

	// Classification
	client := gosseract.NewClient()
	defer client.Close()

	buffer, err := gocv.IMEncode(gocv.PNGFileExt, img)
	if err != nil {
		return "gocv.IMEncode error"
	}
	defer buffer.Close()

	if err := client.SetImageFromBytes(buffer.GetBytes()); err != nil {
		return "client.SetImageFromBytes error"
	}

	client.SetLanguage("eng")
	client.SetPageSegMode(gosseract.PSM_AUTO)

	result, err := client.Text()
	if err != nil {
		return "client.Text error"
	}

	return strings.ReplaceAll(result, " ", "")
}
