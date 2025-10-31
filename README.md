# Projet Puissance 4 (R5.08)

Implémentation d'un jeu de Puissance 4 en Go, utilisant la bibliothèque graphique Ebiten. Ce projet inclut un mode deux joueurs (local) et un mode joueur contre une intelligence artificielle (IA) basée sur l'algorithme Minimax.

Ce dépôt contient le code source complet ainsi que toute la documentation requise (cahier des charges, documentation technique, plan de test) dans le cadre du projet.

## Contexte du Projet

Ce projet a été réalisé dans le cadre de l'unité d'enseignement **R5.08 - Qualité de développement** (Semestre 5).

L'objectif était de développer un jeu vidéo simple en groupe, en respectant un ensemble de bonnes pratiques de développement logiciel. Les livrables requis pour ce projet incluent :

- Le **code source** complet et commenté.
- Un **cahier des charges** décrivant les fonctionnalités.
- Une **documentation technique** détaillant l'architecture.
- Un **plan de test** fonctionnel.
- Des **tests unitaires** pour valider la logique métier.
- La gestion du projet via **GitHub** (revues de code, branches, etc.).

## Fonctionnalités Principales

- **Jeu de Puissance 4 complet** : Grille de 6x7 avec détection de victoire (horizontale, verticale, diagonale) et de match nul.
- **Deux Modes de Jeu** :
  1. **Joueur vs Joueur** : Mode local à deux joueurs sur le même ordinateur.
  2. **Joueur vs IA** : Jouez contre l'ordinateur.
- **Intelligence Artificielle** :
  - Basée sur un algorithme **Minimax avec élagage Alpha-Bêta**.
  - **Difficulté variable** : L'utilisateur peut choisir un niveau de difficulté (1-9) au lancement, ce qui impacte la profondeur de recherche de l'IA.
- **Interface Graphique (UI)** :
  - Interface visuelle simple et réactive construite avec Ebiten.
  - **Animation de chute** des pions avec simulation de gravité.
  - **Indicateurs visuels** : Un "hibou" indique la colonne sélectionnée, un "fantôme" montre le coup de l'IA.
  - **Suivi des scores** (Victoires vs Défaites).
  - Bouton "Rejouer" après la fin d'une partie.

## Technologies Utilisées

- **Langage** : Go (Golang)
- **Bibliothèque Graphique** : Ebiten v2
- **Tests** : Framework de test standard de Go (`package testing`)

## Installation et Lancement

### Prérequis

- Avoir [Go](https://go.dev/doc/install) (version 1.25.3 ou supérieure) installé sur votre machine.

### Lancement

1.  Clonez ce dépôt :

    ```sh
    git clone https://github.com/AbassHammed/c4.git
    cd c4
    ```

2.  Exécutez le jeu depuis le répertoire racine du projet :

    ```sh
    go run -x ./main.go
    ```

### Compilation (Build)

Pour créer un fichier exécutable autonome :

```sh
# Se placer dans le dossier contenant main.go
cd c4

# Compiler
go build

# Vous pouvez maintenant exécuter le binaire généré (./c4 ou c4.exe)
```

## Structure du Projet

L'architecture du projet est conçue pour séparer clairement la logique métier (le "backend" du jeu) de l'interface utilisateur (le "frontend").

```
/
├── c4c/
│   ├── main.go             # Point d'entrée de l'application
│   │
│   ├── game/               # (Backend) Logique de jeu pure
│   │   ├── board.go        # Structure du plateau, détection de victoire
│   │   ├── game_manager.go # Machine à états (tours, état du jeu)
│   │   ├── ai.go           # Logique de l'IA (Minimax Alpha-Beta)
│   │   └── *_test.go       # Tests unitaires pour la logique métier
│   │
│   ├── ui/                 # (Frontend) Interface graphique
│   │   └── game.go         # Boucle de jeu (Update/Draw), gestion des entrées
│   │
│   ├── images/             # Ressources graphiques (embarquées dans le binaire)
│   │   ├── bg.go           # ... (fichiers .go générés à partir des .png)
│   │
│   └── go.mod              # Dépendances (Ebiten)
│
└── Documentation/
    ├── Cahier_des_Charges.pdf
    ├── Documentation_Technique.pdf
    └── Plan_de_Test.xlsx
```

## Auteurs

- ABASS Hammed
- BOUBRIT Maryam
- Hamoudi Mohieddine

## 📜 Licence

Ce projet est distribué sous la licence Apache-2.0. Voir le fichier `LICENSE` pour plus de détails.
