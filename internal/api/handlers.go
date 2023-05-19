package api

import (
	"errors"
	"io"

	"github.com/h2non/filetype"
	"github.com/valyala/fasthttp"
)

func (api *api) generate(ctx *fasthttp.RequestCtx) {
	image, err := api.captcha.Generate()
	if err != nil {
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
	count, err := ctx.QueryArgs().GetUint("count")
	if err != nil {
		badRequest(ctx, errors.New("count is not string!"))
		return
	}

	analysis, accuracy, err := api.captcha.Analyze(ctx, count)
	if err != nil {
		internalServerError(ctx, err)
		return
	}

	ok(ctx, result{
		Accuracy: accuracy,
		Analysis: analysis,
	})
}
