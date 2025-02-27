# Detection de l'OS (Windows ou Linux)
ifeq ($(OS),Windows_NT)
    RM = powershell.exe -Command "Remove-Item -Force -ErrorAction Ignore"
    PSQL = psql.exe
    SHELL := powershell.exe
    
    # Vérifier si DB_NAME est définie dans l'environnement, sinon utiliser tenshi_db
    DB_NAME ?= tenshi_db
else
    RM = rm -f
    PSQL = psql
    SHELL := /bin/bash
    
    # Vérifier si DB_NAME est définie dans l'environnement, sinon utiliser tenshi_db
    DB_NAME ?= tenshi_db
endif

PROJECT_DIR := $(CURDIR)/src
EXEC := $(PROJECT_DIR)/Bot_Tenshi
EXT := $(if $(findstring Windows_NT,$(OS)),.exe,)
BUILD_CMD := go build -o $(EXEC)$(EXT)
RUN_CMD := $(EXEC)$(EXT)

# Charger les variables d'environnement depuis .env (si present)
ifneq (,$(wildcard $(PROJECT_DIR)/.env))
  include $(PROJECT_DIR)/.env
  export
endif

# Affichage de la base de données utilisée
show_db_name:
	@echo "Utilisation de la base de donnees: $(DB_NAME)"

# Creation de la base de donnees
create_db: show_db_name
	@echo "Creation de la base de donnees : $(DB_NAME)"
	@$(PSQL) -U postgres -c "CREATE DATABASE $(DB_NAME);"
	@echo "Base de donnees $(DB_NAME) creee avec succes."

# Suppression de la base de donnees
delete_db: show_db_name
	@echo "Suppression de la base de donnees : $(DB_NAME)"
	@$(PSQL) -U postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	@echo "Base de donnees $(DB_NAME) supprimee avec succes."

# Preparation des modules Go
setup:
ifeq ($(OS),Windows_NT)
	@cd $(PROJECT_DIR); go env -w GOPROXY=https://proxy.golang.org,direct; go mod tidy
else
	@cd $(PROJECT_DIR) && go env -w GOPROXY=https://proxy.golang.org,direct && go mod tidy
endif

# Construction du projet
build:
ifeq ($(OS),Windows_NT)
	@cd $(PROJECT_DIR); $(BUILD_CMD)
else
	@cd $(PROJECT_DIR) && $(BUILD_CMD)
endif

# Execution du projet
run:
ifeq ($(OS),Windows_NT)
	@cd $(PROJECT_DIR); $(RUN_CMD)
else
	@cd $(PROJECT_DIR) && $(RUN_CMD)
endif

# Nettoyage des fichiers de build
clean:
ifeq ($(OS),Windows_NT)
	@$(SHELL) -Command "cd '$(PROJECT_DIR)'; $(RM) '$(EXEC)$(EXT)'; go clean"
else
	@cd $(PROJECT_DIR); $(RM) $(EXEC)$(EXT); go clean
endif

# Cible pour tout detruire, recreer et executer
reboot: clean delete_db create_db setup build run

# Cible pour preparer, construire et executer sans rien detruire
exec: setup build run

deploy: create_db exec

.PHONY: show_db_name create_db delete_db setup build run clean reboot exec deploy