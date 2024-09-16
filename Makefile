PROJECT_DIR := $(CURDIR)
EXEC := Bot_Tenshi

setup:
	cd $(PROJECT_DIR) && go mod tidy

build:
	cd $(PROJECT_DIR) && go build -o $(EXEC) main.go

run:
	cd $(PROJECT_DIR) && ./$(EXEC)

all: setup build run

.PHONY: setup build run all