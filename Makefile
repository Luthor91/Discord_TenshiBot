PROJECT_DIR := $(CURDIR)/src
EXEC := Bot_Tenshi
EXT := $(if $(findstring Windows_NT,$(OS)),.exe,)
BUILD_CMD := go build -o $(EXEC)$(EXT)
RUN_CMD := ./$(EXEC)$(EXT)
CREATE_DB_CMD := psql -U postgres -c "CREATE DATABASE $(DB_NAME)"
DROP_DB_CMD := psql -U postgres -c "DROP DATABASE IF EXISTS $(DB_NAME)"

# Charger les variables d'environnement depuis le fichier .env
ifneq (,$(wildcard $(PROJECT_DIR)/.env))
  include $(PROJECT_DIR)/.env
  export
endif

# Vérification de la variable DB_NAME
check_db_name:
	@if [ -z "$(DB_NAME)" ]; then \
		echo "Erreur : La variable d'environnement DB_NAME n'est pas définie."; \
		exit 1; \
	fi

# Création de la base de données
create_db: check_db_name
	@echo "Création de la base de données : $(DB_NAME)"
	@$(CREATE_DB_CMD)
	@echo "Base de données $(DB_NAME) créée avec succès."

# Suppression de la base de données
delete_db: check_db_name
	@echo "Suppression de la base de données : $(DB_NAME)"
	@$(DROP_DB_CMD)
	@echo "Base de données $(DB_NAME) supprimée avec succès."

# Préparation des modules Go
setup:
	cd $(PROJECT_DIR) && go mod tidy

# Construction du projet
build:
	cd $(PROJECT_DIR) && $(BUILD_CMD)

# Exécution du projet
run:
	cd $(PROJECT_DIR) && $(RUN_CMD)

# Nettoyage des fichiers de build
clean:
	cd $(PROJECT_DIR) && rm -f $(EXEC)$(EXT) && go clean -modcache

# Cible pour tout détruire, recréer et exécuter
reboot: clean delete_db create_db setup build run

# Cible pour préparer, construire et exécuter sans rien détruire
exec: setup build run

.PHONY: check_db_name create_db delete_db setup build run clean reboot exec