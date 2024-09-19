PROJECT_DIR := $(CURDIR)
EXEC := Bot_Tenshi

# Détecter le système d'exploitation
UNAME_S := $(shell uname -s)

# Définir les variables d'extension et de commande en fonction du système d'exploitation
ifeq ($(UNAME_S),Linux)
	EXT :=
	BUILD_CMD := go build -o $(EXEC)
	RUN_CMD := ./$(EXEC)
else
	EXT := .exe
	BUILD_CMD := go build -o $(EXEC)$(EXT)
	RUN_CMD := .\$(EXEC)$(EXT)
endif

# Charger les variables d'environnement depuis le fichier .env
ifneq (,$(wildcard .env))
	include .env
endif

# Cible pour préparer les modules Go
setup:
	cd $(PROJECT_DIR) && go mod tidy

# Cible pour construire le projet
build:
	cd $(PROJECT_DIR) && $(BUILD_CMD)

# Cible pour exécuter le projet
run:
	cd $(PROJECT_DIR) && $(RUN_CMD)

# Cible pour créer la base de données
create_db:
	@if [ -z "$(DB_NAME)" ]; then \
		echo "Erreur : La variable d'environnement DB_NAME n'est pas définie."; \
		exit 1; \
	else \
		echo "Création de la base de données : $(DB_NAME)"; \
		psql -U postgres -c "CREATE DATABASE $(DB_NAME)"; \
		echo "Base de données $(DB_NAME) créée avec succès."; \
	fi

# Cible pour tout construire et exécuter sur le système d'exploitation détecté
all: setup build run

.PHONY: setup build run all