PROJECT_DIR := $(CURDIR)
EXEC := .\Tenshi.exe

setup:
	cd $(PROJECT_DIR) && go mod tidy

build:
	cd $(PROJECT_DIR) && go build

run:
	cd $(PROJECT_DIR) && $(EXEC)

all: setup build run

.PHONY: setup build run all