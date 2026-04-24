# Makefile for donut program
CC = gcc
CFLAGS = -Wall -Wextra -std=c99 -O2
TARGET = wpwd
SRC = wpwd.c

all: $(TARGET)

$(TARGET): $(SRC)
	$(CC) $(CFLAGS) -o $@ $^
	
clean:
	powershell -NoProfile -Command "Remove-Item -Force -ErrorAction SilentlyContinue $(TARGET).exe"

# Phony targets
.PHONY: all clean
