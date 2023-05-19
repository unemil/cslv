package api

import (
	"cslv/internal/model"
	"encoding/json"

	"github.com/mojocn/base64Captcha"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type result struct {
	Solution string           `json:"solution,omitempty"` // solve
	Error    string           `json:"error,omitempty"`    // solve, generate, analyze
	Accuracy float32          `json:"accuracy,omitempty"` // analyze
	Analysis []model.Analysis `json:"analysis,omitempty"` // analyze
}

func ok(ctx *fasthttp.RequestCtx, data any) {
	switch data.(type) {
	case string:
		body, err := json.Marshal(result{Solution: data.(string)})
		if err != nil {
			internalServerError(ctx, err)
			return
		}

		log.Debug().Msg(string(body))

		if _, err := ctx.Write(body); err != nil {
			internalServerError(ctx, err)
			return
		}
	case base64Captcha.Item:
		if _, err := data.(base64Captcha.Item).WriteTo(ctx); err != nil {
			internalServerError(ctx, err)
			return
		}
	case result:
		body, err := json.MarshalIndent(data.(result), "", "\t")
		if err != nil {
			internalServerError(ctx, err)
			return
		}

		log.Debug().Msg(string(body))

		if _, err := ctx.Write(body); err != nil {
			internalServerError(ctx, err)
			return
		}
	case nil:
		ctx.SetStatusCode(fasthttp.StatusNoContent)
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func badRequest(ctx *fasthttp.RequestCtx, err error) {
	body, err := json.Marshal(result{Error: err.Error()})
	if err != nil {
		internalServerError(ctx, err)
		return
	}

	log.Debug().Msg(string(body))

	if _, err := ctx.Write(body); err != nil {
		internalServerError(ctx, err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusBadRequest)
}

func internalServerError(ctx *fasthttp.RequestCtx, err error) {
	body, err := json.Marshal(result{Error: err.Error()})
	if err != nil {
		internalServerError(ctx, err)
		return
	}

	log.Debug().Msg(string(body))

	if _, err := ctx.Write(body); err != nil {
		internalServerError(ctx, err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
}
