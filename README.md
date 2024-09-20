# TenshiBot

TenshiBot est un bot Discord de modération et d'interaction basé sur l'API de Riot Games, conçu pour les communautés de joueurs. Il permet de récupérer des informations sur les matchs, les champions et d'interagir avec les utilisateurs.

## Prérequis

Avant d'installer TenshiBot, assurez-vous d'avoir les éléments suivants :

- **Go** (testé avec la version 1.23.1)
- **Base de données** : une base de données Postgresql
- **Make** 

## Installation

### 1. Cloner le dépôt

Clonez le dépôt GitHub sur votre machine locale :

```bash
git clone https://github.com/Luthor91/Discord_TenshiBot.git
cd Discord_TenshiBot
```

### 2. Installer les dépendances

```bash
make create_db
```

## Exécution

```bash
make all
```
