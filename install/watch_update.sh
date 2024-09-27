#!/bin/bash

# Dossier où se trouve ton dépôt cloné
REPO_DIR="$HOME/Discord_TenshiBot"

# L'URL du dépôt GitHub
REPO_URL="https://github.com/Luthor91/Discord_TenshiBot.git"

# Branche à surveiller
BRANCH="main"

# Fonction pour vérifier les nouvelles mises à jour
check_for_update() {
    # Aller dans le dossier du dépôt
    cd "$REPO_DIR" || exit

    # Récupérer les dernières informations du dépôt distant
    git fetch origin $BRANCH

    # Comparer la dernière version locale avec la dernière version distante
    LOCAL_HASH=$(git rev-parse $BRANCH)
    REMOTE_HASH=$(git rev-parse origin/$BRANCH)

    # Si les deux hash sont différents, il y a une mise à jour
    if [ "$LOCAL_HASH" != "$REMOTE_HASH" ]; then
        echo "Nouvelle mise à jour détectée, récupération des changements..."
        git pull origin $BRANCH

        echo "Exécution de 'make exec'..."
        make exec

        echo "Application mise à jour et exécutée avec succès."
    else
        echo "Aucune mise à jour disponible."
    fi
}

# Vérifier que le dépôt est cloné
if [ ! -d "$REPO_DIR" ]; then
    echo "Le dépôt n'est pas cloné. Clonage en cours..."
    git clone -b $BRANCH $REPO_URL "$REPO_DIR" || { echo "Échec du clonage du dépôt."; exit 1; }
fi

# Boucle de surveillance
while true; do
    echo "Vérification des mises à jour..."
    check_for_update

    # Attendre 5 minutes avant de vérifier à nouveau (modifiable)
    sleep 300
done