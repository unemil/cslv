package api

import (
	"context"
	"cslv/internal/model"
)

type Service interface {
	Generate() (model.Captcha, error)
	Solve(file []byte) (string, error)
	Analyze(ctx context.Context, count int) ([]model.Analysis, float32, error)
}
