MODELS = WenQuanYiMicroHei DenneThreedee

MAX_ITERATIONS_WENQUANYIMICROHEI = 1500
MAX_ITERATIONS_DENNETHREEDEE = 10000

build:
	cd templates/ui && \
	npm install && \
	npm run build
	
	go mod tidy
	go mod vendor
	go build -mod=vendor -o app cmd/main.go

run:
	go run cmd/main.go

test:
	go test -v ./...

train:
	cd training/tesstrain && \
		if [ "$(MODEL)" = "WenQuanYiMicroHei" ]; then \
			TESSDATA_PREFIX=../tesseract/tessdata make training MODEL_NAME=$(MODEL) START_MODEL=eng TESSDATA=../tesseract/tessdata MAX_ITERATIONS=$(MAX_ITERATIONS_WENQUANYIMICROHEI); \
			cp data/WenQuanYiMicroHei.traineddata /usr/share/tesseract-ocr/4.00/tessdata/WenQuanYiMicroHei.traineddata; \
		elif [ "$(MODEL)" = "DenneThreedee" ]; then \
			TESSDATA_PREFIX=../tesseract/tessdata make training MODEL_NAME=$(MODEL) START_MODEL=eng TESSDATA=../tesseract/tessdata MAX_ITERATIONS=$(MAX_ITERATIONS_DENNETHREEDEE); \
			cp data/DenneThreedee.traineddata /usr/share/tesseract-ocr/4.00/tessdata/DenneThreedee.traineddata; \
		fi