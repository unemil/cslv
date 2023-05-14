package api

import (
	"io"
	"log"

	"github.com/valyala/fasthttp"
)

func (api *api) solve(ctx *fasthttp.RequestCtx) {
	header, err := ctx.FormFile("file")
	if err != nil {
		log.Println(err)
		return
	}

	file, err := header.Open()
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}

	ctx.WriteString(api.image.Solve(bytes))
}

func (api *api) generate(ctx *fasthttp.RequestCtx) {
	image, err := api.captcha.Image()
	if err != nil {
		log.Println(err)
		return
	}

	if _, err := image.WriteTo(ctx); err != nil {
		log.Println(err)
		return
	}
}
