package api

import (
	"bytes"
	"cslv/internal/model"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/h2non/filetype"
	"github.com/valyala/fasthttp"
)

func (api *api) generate(ctx *fasthttp.RequestCtx) {
	image, err := api.captcha.Generate()
	if err != nil {
		internalServerError(ctx, err)
		return
	}

	if err := saveInfo(image); err != nil {
		internalServerError(ctx, err)
		return
	}

	ok(ctx, image.Item)
}

func (api *api) solve(ctx *fasthttp.RequestCtx) {
	header, err := ctx.FormFile("file")
	if err != nil {
		badRequest(ctx, err)
		return
	}

	file, err := header.Open()
	if err != nil {
		internalServerError(ctx, err)
		return
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		internalServerError(ctx, err)
		return
	}

	if !filetype.IsImage(bytes) || filetype.IsExtension(bytes, "gif") {
		badRequest(ctx, errors.New("file is not image!"))
		return
	}

	text, err := api.captcha.Solve(bytes)
	if err != nil {
		internalServerError(ctx, err)
		return
	}

	ok(ctx, text)
}

func (api *api) analyze(ctx *fasthttp.RequestCtx) {
	matches, err := filepath.Glob("captchas/*.png")
	if err != nil {
		internalServerError(ctx, err)
		return
	}

	re := regexp.MustCompile(`^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`)

	info := make([]model.AnalyzeInfo, 0, len(matches))
	for _, match := range matches {
		data, err := getInfo(re.FindStringSubmatch(match)[2])
		if err != nil {
			internalServerError(ctx, err)
			return
		}

		image, err := ioutil.ReadFile(match)
		if err != nil {
			internalServerError(ctx, err)
			return
		}

		buf := new(bytes.Buffer)
		if _, err := buf.Write(image); err != nil {
			internalServerError(ctx, err)
			return
		}

		info = append(info, model.AnalyzeInfo{
			Image: image,
			Data:  data,
		})
	}

	analysis, accuracy, err := api.captcha.Analyze(ctx, info)
	if err != nil {
		internalServerError(ctx, err)
		return
	}

	ok(ctx, result{
		Accuracy: accuracy,
		Analysis: analysis,
	})
}
