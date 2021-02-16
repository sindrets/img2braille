BIN_NAME = img2braille

SRC = $(shell git ls-files '*.go')
OUT_DIR = build
OUT = $(OUT_DIR)/$(BIN_NAME)

.PHONY: all\
		build\
		pre-build\
		clean

all: build

$(OUT): $(SRC)
	go build -o $(OUT_DIR)/$(BIN_NAME) $(SRC)

build: pre-build $(OUT)
	@echo Build DONE

pre-build:
	mkdir -p $(OUT_DIR)

clean:
	rm -rf $(OUT_DIR)
	@echo Clean DONE
