BINARY := file
PKG := .
OUT := ./$(BINARY)
EXT := .exe
RM := cmd /C del /Q
RUN_CMD := .\$(BINARY)$(EXT)

.PHONY: all build run install clean fmt vet

all: build

build:
	go build -o $(OUT)$(EXT) $(PKG)

run: build
	$(RUN_CMD)

install:
	go install $(PKG)

clean:
	-@$(RM) $(OUT)$(EXT)
	go clean

