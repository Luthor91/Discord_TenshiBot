PROJECT_DIR := $(CURDIR)
EXEC := Bot_Tenshi

# Definir les variables d'extension et de commande en fonction du système d'exploitation
ifeq ($(OS),Windows_NT)
	EXT := .exe
	BUILD_CMD := go build -o $(EXEC)$(EXT)
	RUN_CMD := .\$(EXEC)$(EXT)
	CREATE_DB_CMD := psql -U postgres -c "CREATE DATABASE $(DB_NAME)"
else
	EXT :=
	BUILD_CMD := go build -o $(EXEC)
	RUN_CMD := ./$(EXEC)
	CREATE_DB_CMD := psql -U postgres -c "CREATE DATABASE $(DB_NAME)"
endif

# Charger les variables d'environnement depuis le fichier .env
ifneq (,$(wildcard .env))
	include .env
endif

# Cible pour preparer les modules Go
setup:
	cd $(PROJECT_DIR) && go mod tidy

# Cible pour construire le projet
build:
	cd $(PROJECT_DIR) && $(BUILD_CMD)

# Cible pour executer le projet
run:
	cd $(PROJECT_DIR) && $(RUN_CMD)

# Cible pour creer la base de donnees
create_db:
ifeq ($(OS),Windows_NT)
	@if "$(DB_NAME)"=="" ( \
		echo Erreur : La variable d'environnement DB_NAME n'est pas definie. ;\
		exit 1 ;\
	) else ( \
		echo Creation de la base de donnees : $(DB_NAME) ;\
		$(CREATE_DB_CMD) ;\
		echo Base de donnees $(DB_NAME) creee avec succès. ;\
	)
else
	@if [ -z "$(DB_NAME)" ]; then \
		echo "Erreur : La variable d'environnement DB_NAME n'est pas definie."; \
		exit 1; \
	else \
		echo "Creation de la base de donnees : $(DB_NAME)"; \
		$(CREATE_DB_CMD); \
		echo "Base de donnees $(DB_NAME) creee avec succès."; \
	fi
endif

# Cible pour nettoyer le projet et le cache
clean:
	cd $(PROJECT_DIR) && rm -f $(EXEC)$(EXT) && go clean -modcache

# Cible pour tout construire et executer sur le système d'exploitation detecte
all: setup build run

.PHONY: setup build run create_db clean all