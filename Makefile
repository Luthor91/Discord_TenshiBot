PROJECT_DIR := $(CURDIR)
EXEC := Bot_Tenshi

setup:
	cd $(PROJECT_DIR) && go mod tidy

build_linux:
	cd $(PROJECT_DIR) && go build -o $(EXEC) main.go

build_windows:
	cd $(PROJECT_DIR) && go build -o $(EXEC).exe main.go

run_linux:
	cd $(PROJECT_DIR) && ./$(EXEC)

run_windows:
	cd $(PROJECT_DIR) && .\$(EXEC)

lall: setup build_linux run_linux

wall: setup build_windows run_windows

.PHONY: setup build run_linux run_windows lall wall