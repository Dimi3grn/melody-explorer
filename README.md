# MelodyExplorer - Découverte Musicale avec Spotify

## Présentation du Projet

MelodyExplorer est une application web de découverte musicale qui utilise l'API Spotify pour permettre aux utilisateurs de rechercher, explorer et sauvegarder leurs artistes, albums et morceaux préférés. Cette application a été développée dans le cadre du projet "Groupie Tracker" pour démontrer l'utilisation d'une API REST et l'implémentation de diverses fonctionnalités de gestion de données.

### Thème
Le thème du projet est la découverte musicale - permettre aux utilisateurs d'explorer la vaste bibliothèque de musique de Spotify à travers une interface simple et conviviale.

### Fonctionnalités Principales
- **Système de Recherche** : Recherche d'artistes, d'albums et de morceaux par mots-clés
- **Système de Filtres** : Filtrage par type (artiste, album, morceau), genre, popularité, année de sortie
- **Système de Pagination** : Navigation à travers les résultats par lots de 10, 20 ou 30 items
- **Système de Favoris** : Ajout et suppression d'éléments à une liste de favoris persistante
- **Détails** : Affichage des informations détaillées sur les artistes, albums et morceaux
- **Catégories** : Exploration de la musique par genres

## Installation et Lancement

### Prérequis
- Go (version 1.16 ou supérieure)
- Compte développeur Spotify

### Configuration de l'API Spotify
1. Créez un compte sur [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/)
2. Créez une nouvelle application
3. Notez votre Client ID et Client Secret
4. Ajoutez `http://localhost:8085/callback` comme URI de redirection

### Installation
1. Clonez ce dépôt
   ```
   git clone https://github.com/Dimi3grn/melody-explorer.git
   cd melody-explorer
   ```

2. Installez les dépendances
   ```
   go mod download
   ```

### Lancement
1. Lancez l'application
   ```
   go run ./cmd/server/main.go
   ```

2. Ouvrez votre navigateur et accédez à `http://localhost:8085`

## Routes Implémentées

### Pages
- `GET /` - Page d'accueil
- `GET /search` - Recherche de musique
- `GET /collection` - Collection musicale avec filtres
- `GET /artist/{id}` - Détails d'un artiste
- `GET /album/{id}` - Détails d'un album
- `GET /track/{id}` - Détails d'un morceau
- `GET /favorites` - Gestion des favoris
- `GET /category/{genre}` - Exploration par genre
- `GET /about` - À propos du projet

### Authentification
- `GET /login` - Connexion via Spotify
- `GET /callback` - Callback après authentification Spotify
- `GET /logout` - Déconnexion

### API
- `POST /api/favorites/add` - Ajouter un élément aux favoris
- `POST /api/favorites/remove` - Supprimer un élément des favoris

## Endpoints Spotify Utilisés

Ce projet utilise les endpoints suivants de l'API Spotify :

- `/search` - Recherche d'artistes, d'albums et de morceaux
- `/artists/{id}` - Obtention des détails d'un artiste
- `/artists/{id}/albums` - Obtention des albums d'un artiste
- `/albums/{id}` - Obtention des détails d'un album
- `/albums/{id}/tracks` - Obtention des morceaux d'un album
- `/tracks/{id}` - Obtention des détails d'un morceau
- `/browse/categories` - Obtention des catégories de musique
- `/browse/new-releases` - Obtention des nouvelles sorties

## Synthèse du Déroulement du Projet

### Comment ai-je décomposé le projet ? Quelles ont été les phases clé ?

J'ai décomposé le projet en plusieurs phases distinctes :

1. **Phase d'analyse et de conception (3 jours)**
   - Exploration de l'API Spotify et ses fonctionnalités
   - Définition des besoins et des fonctionnalités
   - Conception de l'architecture de l'application
   - Réalisation de maquettes pour les différentes pages

2. **Phase de développement du backend (5 jours)**
   - Mise en place de l'authentification Spotify OAuth
   - Création des clients pour interagir avec l'API Spotify
   - Implémentation des modèles de données
   - Développement du système de stockage pour les favoris
   - J'ai appris à structurer mon code en packages pour une meilleure organisation

3. **Phase de développement du frontend (5 jours)**
   - Création des templates HTML
   - Développement des styles CSS
   - Implémentation des fonctionnalités JavaScript pour l'interaction utilisateur

4. **Phase d'intégration et de tests (2 jours)**
   - Intégration des différentes composantes
   - Tests des fonctionnalités
   - Correction des bugs

5. **Phase de finalisation (1 jour)**
   - Optimisations finales
   - Documentation
   - Préparation du rendu

### Comment ai-je réparti les tâches entre les membres de l'équipe ?

Comme ce projet était individuel, j'ai dû gérer toutes les tâches moi-même. Pour m'organiser efficacement, j'ai :

- Utilisé un tableau Kanban personnel pour suivre les tâches
- Défini des objectifs quotidiens clairs
- Maintenu un journal de développement pour suivre ma progression
- Alterné entre les tâches backend et frontend pour éviter la monotonie

### Comment ai-je géré mon temps ? Ai-je défini des priorités ?

Pour gérer efficacement mon temps, j'ai :

1. **Défini des priorités claires**
   - J'ai donné la priorité aux fonctionnalités essentielles (recherche, filtres, pagination, favoris)
   - Les améliorations visuelles et les fonctionnalités bonus ont été traitées une fois les éléments critiques terminés

2. **Créé un calendrier de développement**
   - Allocation de temps spécifique pour chaque phase du projet
   - Points de contrôle réguliers pour évaluer ma progression

3. **Identifié les risques potentiels**
   - L'authentification OAuth était un point critique que j'ai abordé tôt
   - J'ai prévu du temps supplémentaire pour les fonctionnalités complexes

4. **Séparé les journées de travail en blocs**
   - Sessions de codage concentrées (2-3 heures)
   - Pauses régulières pour maintenir la productivité
   - Sessions de revue de code et de test en fin de journée

### Quelle stratégie ai-je adoptée pour me documenter ?

Pour me documenter efficacement, j'ai :

1. **Étudié la documentation officielle**
   - Documentation de l'API Spotify
   - Documentation de Go et ses packages
   - J'ai appris à utiliser les fichiers `.env` pour gérer les variables d'environnement sensibles

2. **Consulté des ressources complémentaires**
   - Tutoriels sur l'authentification OAuth
   - Exemples de code pour l'interaction avec des APIs REST
   - Bonnes pratiques pour l'architecture d'applications web en Go

3. **Expérimenté avec Postman**
   - Tests des endpoints de l'API Spotify
   - Compréhension des formats de données et des paramètres

4. **Maintenu une documentation interne**
   - Commentaires détaillés dans le code
   - Documentation des API internes
   - Notes sur les décisions d'architecture

5. **Participé à des forums et communautés**
   - Stack Overflow pour des questions spécifiques
   - Communauté Go pour des conseils sur les meilleures pratiques

## Conclusion

Le développement de MelodyExplorer a été une expérience enrichissante qui m'a permis d'approfondir mes connaissances en développement web avec Go et en intégration d'API externes. Le projet m'a également donné l'occasion de mettre en pratique des compétences en gestion de projet, notamment en termes de planification, de définition des priorités et de documentation.

Les défis rencontrés, notamment avec l'authentification OAuth et la manipulation des réponses de l'API Spotify, ont été des opportunités d'apprentissage précieuses. Je suis particulièrement satisfait(e) de la mise en œuvre du système de favoris persistants et de l'interface utilisateur intuitive pour la découverte musicale.
