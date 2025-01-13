# Application de Gestion des Associations Scolaires

### Une application développée en Go pour gérer les associations au sein d'une école.

## Équipe de développement

* Riadh Benchouche
* Riad Ahmed Yahia
* Hicham Aissaoui
* HSU WANG Kevin

## Prérequis

* Docker
* Go (version recommandée : dernière version stable)

## Installation

### Clonez le repository :

```bash
git clone git@github.com:riadh-benchouche/challenge-s4.git
cd challenge-s4
```

### Configurez les variables d'environnement :

```bash
 cp .env.example .env
```

⚠️ N'oubliez pas de modifier les valeurs dans le fichier .env selon votre environnement

### Lancez l'application avec Docker :

```bash
docker-compose up -d
```

### Documentation API

La documentation Swagger de l'API est accessible à l'adresse :

`http://localhost:3000/swagger/ui`

### Utilisation

Assurez-vous que tous les conteneurs Docker sont en cours d'exécution :

`docker-compose ps`

### Développement

Pour lancer l'application en mode développement :

```bash
docker-compose up -d --build
```

### Arrêt de l'application

Pour arrêter les conteneurs Docker :

```bash
docker-compose down
```

## Sans Docker (développement local)

```bash
go run main.go
```
