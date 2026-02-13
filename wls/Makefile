BINARY := wls
PKG := .
OUT := ./$(BINARY)
EXT := .exe
RM := cmd /C del /Q
RUN_CMD := .\$(BINARY)$(EXT)
WLL_SRC := wll.c
WLL_BIN := ./wll.exe
WLA_SRC := wla.c
WLA_BIN := ./wla.exe

.PHONY: all build build-wll build-wla run install clean fmt vet

all: build build-wll build-wla

build:
	go build -o $(OUT)$(EXT) wls.go

build-wll:
	gcc -O2 -o $(WLL_BIN) $(WLL_SRC)

build-wla:
	gcc -O2 -o $(WLA_BIN) $(WLA_SRC)

run: build
	$(RUN_CMD)

install:
	go install $(PKG)

clean:
	-@$(RM) $(OUT)$(EXT)
	-@$(RM) $(WLL_BIN)
	go clean
