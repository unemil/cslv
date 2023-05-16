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
	resize := gocv.NewMat()
	gocv.Resize(img, &resize, image.Pt(img.Cols()*2, img.Rows()*2), 0, 0, gocv.InterpolationLinear)

	grayscale := gocv.NewMat()
	gocv.CvtColor(resize, &grayscale, gocv.ColorBGRToGray)

	filter := gocv.NewMat()
	gocv.BilateralFilter(grayscale, &filter, 5, 15, 15)

	// gocv.IMWrite("preprocessing.png", filter)

	// Segmentation
	threshold := gocv.NewMat()
	gocv.Threshold(filter, &threshold, 210, 255, gocv.ThresholdBinary)

	// gocv.IMWrite("segmentation.png", threshold)

	// Classification
	client := gosseract.NewClient()
	defer client.Close()

	buffer, err := gocv.IMEncode(gocv.PNGFileExt, threshold)
	if err != nil {
		log.Println(err)
		return "gocv.IMEncode error"
	}
	defer buffer.Close()

	if err := client.SetImageFromBytes(buffer.GetBytes()); err != nil {
		log.Println(err)
		return "client.SetImageFromBytes error"
	}

	client.SetLanguage("WenQuanYiMicroHei", "ChromosomeHeavy", "DenneThreedee")
	client.SetPageSegMode(gosseract.PSM_AUTO)

	result, err := client.Text()
	if err != nil {
		log.Println(err)
		return "client.Text error"
	}

	return strings.ReplaceAll(result, " ", "")
}
