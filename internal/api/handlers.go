package api

import (
	"errors"
	"io"

	"github.com/h2non/filetype"
	"github.com/valyala/fasthttp"
)

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

	text, err := api.image.Solve(bytes)
	if err != nil {
		internalServerError(ctx, err)
		return
	}

	ok(ctx, text)
}

func (api *api) generate(ctx *fasthttp.RequestCtx) {
	image, err := api.captcha.Image()
	if err != nil {
		internalServerError(ctx, err)
		return
	}

	ok(ctx, image.Item)
}
