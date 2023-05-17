package api

import (
	"encoding/json"

	"github.com/mojocn/base64Captcha"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type result struct {
	Text  string `json:"text,omitempty"`
	Error string `json:"error,omitempty"`
}

func ok(ctx *fasthttp.RequestCtx, data any) {
	switch data.(type) {
	case string:
		body, err := json.Marshal(result{Text: data.(string)})
		if err != nil {
			internalServerError(ctx, err)
			return
		}

		if _, err := ctx.Write(body); err != nil {
			internalServerError(ctx, err)
			return
		}
	case base64Captcha.Item:
		if _, err := data.(base64Captcha.Item).WriteTo(ctx); err != nil {
			internalServerError(ctx, err)
			return
		}
	case nil:
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func badRequest(ctx *fasthttp.RequestCtx, err error) {
	log.Error().Msg(err.Error())

	body, err := json.Marshal(result{Error: err.Error()})
	if err != nil {
		internalServerError(ctx, err)
		return
	}

	if _, err := ctx.Write(body); err != nil {
		internalServerError(ctx, err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusBadRequest)
}

func internalServerError(ctx *fasthttp.RequestCtx, err error) {
	log.Error().Msg(err.Error())

	body, err := json.Marshal(result{Error: err.Error()})
	if err != nil {
		internalServerError(ctx, err)
		return
	}

	if _, err := ctx.Write(body); err != nil {
		internalServerError(ctx, err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
}
