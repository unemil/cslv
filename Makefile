MODELS = WenQuanYiMicroHei DenneThreedee ChromosomeHeavy

MAX_ITERATIONS_WENQUANYIMICROHEI = 100
MAX_ITERATIONS_DENNETHREEDEE = 100
MAX_ITERATIONS_CHROMOSOMEHEAVY = 100

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
		for model in $(MODELS); do \
			echo ; \
			if [ "$$model" = "WenQuanYiMicroHei" ]; then \
				TESSDATA_PREFIX=../tesseract/tessdata make training MODEL_NAME=$$model START_MODEL=eng TESSDATA=../tesseract/tessdata MAX_ITERATIONS=$(MAX_ITERATIONS_WENQUANYIMICROHEI) && \
				cp data/WenQuanYiMicroHei.traineddata /usr/share/tesseract-ocr/4.00/tessdata/WenQuanYiMicroHei.traineddata; \
			elif [ "$$model" = "DenneThreedee" ]; then \
				TESSDATA_PREFIX=../tesseract/tessdata make training MODEL_NAME=$$model START_MODEL=eng TESSDATA=../tesseract/tessdata MAX_ITERATIONS=$(MAX_ITERATIONS_DENNETHREEDEE) && \
				cp data/DenneThreedee.traineddata /usr/share/tesseract-ocr/4.00/tessdata/DenneThreedee.traineddata; \
			elif [ "$$model" = "ChromosomeHeavy" ]; then \
				TESSDATA_PREFIX=../tesseract/tessdata make training MODEL_NAME=$$model START_MODEL=eng TESSDATA=../tesseract/tessdata MAX_ITERATIONS=$(MAX_ITERATIONS_CHROMOSOMEHEAVY) && \
				cp data/ChromosomeHeavy.traineddata /usr/share/tesseract-ocr/4.00/tessdata/ChromosomeHeavy.traineddata; \
			else \
				echo "Unknown model $$model"; \
			fi; \
			echo ; \
		done