# TenshiBot

TenshiBot est un bot Discord de modération et d'interaction basé sur l'API de Riot Games, conçu pour les communautés de joueurs. Il permet de récupérer des informations sur les matchs, les champions et d'interagir avec les utilisateurs.

## Prérequis

Avant d'installer TenshiBot, assurez-vous d'avoir les éléments suivants :

- **Go** (testé avec la version 1.23.1)
- **Base de données** : une base de données Postgresql
- **Make** 

## Installation

Clonez le dépôt GitHub sur votre machine locale :

```bash
git clone https://github.com/Luthor91/Discord_TenshiBot.git
cd Discord_TenshiBot
```

## Exécution

Pour la première exécution ou si vous voulez relancer le programme en faisant un reset de la base de données, il est recommandé d'effectuer la commande ci-dessous, elle permet de ré-installer les dépendances de zéro.

```bash
make reboot
```

Pour les exécutions suivantes, il est conseillé d'effectuer la commande ci-dessous.

```bash
make exec
```

## Structure

Le code suit une structure avec des Controller pour gérer les interractions avec la base de données, des Services qui appelleront les Controllers et ajouteront une logique en plus si nécessaire.

Les commandes discord sont regroupé dans le dossier commands.

