#!/bin/bash

# Fonction pour afficher les erreurs
error() {
    echo "Erreur: $1"
    exit 1
}

# Fonction pour demander les tokens à l'utilisateur
ask_user_input() {
    read -p "Veuillez entrer votre $1: " value
    if [ -z "$value" ]; then
        echo "Erreur: $1 ne peut pas être vide."
        exit 1
    fi
    echo "$value"
}

# Fonction pour définir et persister une variable d'environnement
set_env_variable() {
    if grep -q "export $1=" ~/.bashrc; then
        sed -i "s|export $1=.*|export $1=\"$2\"|" ~/.bashrc
    else
        echo "export $1=\"$2\"" >> ~/.bashrc
    fi
    export $1="$2"
}

# Mise à jour des paquets
echo "Mise à jour des paquets..."
sudo apt update || error "Impossible de mettre à jour les paquets"

# Installation de Golang (v1.23.1 minimum)
echo "Installation de Golang..."
GO_VERSION="1.23.1"
if ! go version >/dev/null 2>&1 || [[ $(go version | awk '{print $3}') < "go$GO_VERSION" ]]; then
    wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz || error "Échec du téléchargement de Go"
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz || error "Erreur d'installation de Go"
    rm go${GO_VERSION}.linux-amd64.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    source ~/.bashrc
fi

# Vérification de l'installation de Go
go version || error "Installation de Go échouée"

# Installation de PostgreSQL
echo "Installation de PostgreSQL..."
sudo apt install -y postgresql postgresql-contrib || error "Erreur d'installation de PostgreSQL"

# Installation de Make
echo "Installation de Make..."
sudo apt install -y make || error "Erreur d'installation de Make"

# Installation de Git
echo "Installation de Git..."
sudo apt install -y git || error "Erreur d'installation de Git"

# Récupérer les tokens auprès de l'utilisateur
TOKEN=$(ask_user_input "TOKEN")
API_RIOT=$(ask_user_input "API_RIOT")
API_LOL_ESPORT=$(ask_user_input "API_LOL_ESPORT")

# Définir les variables d'environnement de manière permanente
echo "Configuration des variables d'environnement..."

set_env_variable "TOKEN" "$TOKEN"
set_env_variable "API_RIOT" "$API_RIOT"
set_env_variable "API_LOL_ESPORT" "$API_LOL_ESPORT"
set_env_variable "DB_NAME" "manchoux_db"
set_env_variable "DB_ADMIN_USER" "postgres"
set_env_variable "DB_ADMIN_PASSWORD" "root"
set_env_variable "DB_HOST" "localhost"
set_env_variable "DB_PORT" "5432"
set_env_variable "DB_SSL_MODE" "disable"
set_env_variable "DB_USER" "postgres"
set_env_variable "DB_PASSWORD" "root"
set_env_variable "PREFIX" "?"
set_env_variable "LOL_PATCH_VERSION" "12.6.1"
set_env_variable "LOL_REGION" "euw1"
set_env_variable "LOL_SERVER" "euw1"

# Recharger le fichier .bashrc pour appliquer immédiatement les modifications
source ~/.bashrc

# Clonage du dépôt Git
echo "Clonage du dépôt Discord_TenshiBot..."
git clone https://github.com/Luthor91/Discord_TenshiBot.git || error "Erreur de clonage du dépôt"

# Se déplacer dans le dépôt
cd Discord_TenshiBot || error "Le dépôt n'existe pas"

# Exécution de make reboot
echo "Exécution de 'make reboot'..."
make reboot || error "Échec de l'exécution de make reboot"

echo "Recherche du fichier d'update."

# Recherche du fichier watch_update.sh dans le répertoire courant et ses sous-répertoires
FILE_PATH=$(find . -type f -name "watch_update.sh")
# Vérification si le fichier a été trouvé
if [ -z "$FILE_PATH" ]; then
    echo "Le fichier watch_update.sh n'a pas été trouvé."
    exit 1
else
    echo "Le fichier watch_update.sh a été trouvé à l'emplacement : $FILE_PATH"
    echo "Exécution du script..."
    chmod +x "$FILE_PATH"
    "$FILE_PATH"
fi

echo "Installation terminée avec succès."