package api

import (
	"bytes"
	"cslv/internal/generator/captcha"
	"cslv/internal/service/image"
	"io"
	"mime/multipart"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func TestSolve(t *testing.T) {
	api := New(&Config{
		Host:         ":80",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	}, image.New(), captcha.New())

	testCases := []struct {
		name               string
		file               string
		requestBody        []byte
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:               "ValidRequest",
			file:               "../../tests/captcha.png",
			requestBody:        nil,
			expectedStatusCode: fasthttp.StatusOK,
			expectedResponse:   `{"text": "qtjed"}`,
		},
		{
			name:               "InvalidRequest",
			file:               "../../tests/captcha.txt",
			requestBody:        nil,
			expectedStatusCode: fasthttp.StatusBadRequest,
			expectedResponse:   `{"error": "file is not image!"}`,
		},
	}

	ml := fasthttputil.NewInmemoryListener()
	defer ml.Close()

	go func() {
		s := &fasthttp.Server{
			Handler: api.router,
		}

		if err := s.Serve(ml); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	}()

	time.Sleep(time.Second)

	for _, tc := range testCases {
		buf := new(bytes.Buffer)
		writer := multipart.NewWriter(buf)

		file, err := os.Open(tc.file)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		defer file.Close()

		formFile, err := writer.CreateFormFile("file", tc.file)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if _, err = io.Copy(formFile, file); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		writer.Close()

		tc.requestBody = buf.Bytes()

		t.Run(tc.name, func(t *testing.T) {
			req := fasthttp.AcquireRequest()
			defer fasthttp.ReleaseRequest(req)

			resp := fasthttp.AcquireResponse()
			fasthttp.ReleaseResponse(resp)

			req.Header.SetHost(api.config.Host)
			req.Header.SetMethod("POST")
			req.SetRequestURI("/api/v1/captcha/solve")
			req.Header.SetContentType(writer.FormDataContentType())
			req.SetBody(buf.Bytes())

			c := &fasthttp.Client{
				Dial: func(addr string) (net.Conn, error) {
					return ml.Dial()
				},
			}

			if err := c.Do(req, resp); err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode())
			assert.JSONEq(t, tc.expectedResponse, string(resp.Body()))
		})
	}
}
