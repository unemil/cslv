package api

import (
	"cslv/internal/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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

func saveInfo(image model.Captcha) error {
	data, err := json.MarshalIndent(image, "", "\t")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fmt.Sprintf("captchas/%s.json", image.ID), data, 0644); err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("captchas/%s.png", image.ID))
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := image.Item.WriteTo(file); err != nil {
		return err
	}

	return nil
}

func getInfo(id string) (model.Captcha, error) {
	file, err := os.Open(fmt.Sprintf("captchas/%s.json", id))
	if err != nil {
		return model.Captcha{}, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return model.Captcha{}, err
	}

	var info model.Captcha
	if err := json.Unmarshal(data, &info); err != nil {
		return model.Captcha{}, err
	}

	return info, nil
}
