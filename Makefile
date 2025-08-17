BIN_NAME = img2braille

SRC = $(foreach f,$(shell git ls-files '*.go'),$(abspath $(f)))
OUT_DIR = build
OUT = $(abspath $(OUT_DIR)/$(BIN_NAME))

$(info $(SRC))

.PHONY: all\
		build\
		pre-build\
		install\
		clean

all: build

$(OUT): $(SRC)
	cd src && \
	go build -o $(OUT) $(SRC)

build: pre-build $(OUT)
	@ echo Build DONE

pre-build:
	mkdir -p $(OUT_DIR)

install:
	@ if [[ "$$(id -u)" -eq 0 ]]; then \
		cp -f $(OUT) /usr/bin/$(BIN_NAME); \
		echo "Installed as '/usr/bin/$(BIN_NAME)'." > /dev/stderr; \
	else \
		mkdir -p ~/.local/bin/; \
		cp -f $(OUT) ~/.local/bin/$(BIN_NAME); \
		echo "Installed as '~/.local/bin/$(BIN_NAME)'." > /dev/stderr; \
	fi

clean:
	rm -rf $(OUT_DIR)
	@ echo Clean DONE
