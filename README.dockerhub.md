# OPT-NC Fiber Eligibility Checker

![Docker Pulls](https://img.shields.io/docker/pulls/rastadidi/optnc-fiber-eligibility)
![Docker Image Size](https://img.shields.io/docker/image-size/rastadidi/optnc-fiber-eligibility)
![Docker Image Version](https://img.shields.io/docker/v/rastadidi/optnc-fiber-eligibility)

**API REST de vérification d'éligibilité à la fibre optique OPT Nouvelle-Calédonie**

🚀 Testez l'éligibilité d'un numéro de téléphone fixe à la fibre optique en Nouvelle-Calédonie via une API REST simple.

## 🎯 Quick Start

### Avec Podman

```bash
# Démarrer l'API
podman run -d -p 8080:8080 --name optnc-eligibilite-fibre  rastadidi/optnc-fiber-eligibility

# Tester
curl http://localhost:8080/health
```

### Avec Docker

```bash
# Démarrer l'API
docker run -d -p 8080:8080 --name opt-api rastadidi/optnc-fiber-eligibility

# Tester
curl http://localhost:8080/health
```

## 📚 Swagger UI

Une fois l'API démarrée, accédez à la documentation interactive :

**http://localhost:8080/swagger/**

## 🔥 Exemples avec HTTPie

[HTTPie](https://httpie.io/) est un client HTTP moderne et intuitif.

### Installation HTTPie

```bash
# Ubuntu/Debian
sudo apt install httpie

# macOS
brew install httpie

# Fedora
sudo dnf install httpie
```

### Vérifier le health de l'API

```bash
http GET :8080/health
```

**Réponse :**
```json
{
    "service": "opt-nc-eligibility-api",
    "status": "ok",
    "timestamp": "2025-12-13T22:20:00Z"
}
```

### Vérifier l'éligibilité (GET)

```bash
# Format simple
http GET :8080/api/v1/eligibility phone==257364

# Avec un numéro formaté (avec points)
http GET :8080/api/v1/eligibility phone=="25.73.64"
```

**Réponse (200 OK) :**
```json
{
    "success": true,
    "data": {
        "phone_number": "257364",
        "checked_at": "2025-12-13T22:20:00Z",
        "found": true,
        "adsl": {
            "status": "non-eligible",
            "message": "L'offre souscrite sur votre ligne n'est pas compatible avec l'ADSL."
        },
        "fiber": {
            "status": "non-eligible",
            "available": false,
            "message": "Votre ligne n'est pas encore éligible à la fibre optique."
        },
        "contact_phone": "1016",
        "isp_providers": [
            {"name": "can'l", "url": "http://www.canl.nc/"},
            {"name": "InternetNC", "url": "http://www.internetnc.nc/"},
            {"name": "Lagoon", "url": "http://www.lagoon.nc/"},
            {"name": "MLS", "url": "http://www.mls.nc/"},
            {"name": "Nautile", "url": "http://www.nautile.nc/"}
        ]
    }
}
```

### Vérifier l'éligibilité (POST)

```bash
# Format JSON
http POST :8080/api/v1/eligibility phone_number=257364

# Avec un numéro formaté
http POST :8080/api/v1/eligibility phone_number="25.73.64"
```

### Extraire uniquement les infos fibre avec jq

```bash
# Voir uniquement les informations fibre
http GET :8080/api/v1/eligibility phone==257364 | jq '.data.fiber'
```

**Sortie :**
```json
{
  "status": "non-eligible",
  "available": false,
  "message": "Votre ligne n'est pas encore éligible à la fibre optique."
}
```

### Vérification rapide de disponibilité

```bash
# Retourne true ou false
http GET :8080/api/v1/eligibility phone==257364 | jq '.data.fiber.available'
```

**Sortie :**
```
false
```

### Format personnalisé avec jq

```bash
# Format compact lisible
http GET :8080/api/v1/eligibility phone==257364 | \
  jq -r '"\(.data.phone_number) -> Fibre: \(.data.fiber.status) (disponible: \(.data.fiber.available))"'
```

**Sortie :**
```
257364 -> Fibre: non-eligible (disponible: false)
```

## 📡 Endpoints disponibles

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| `GET` | `/health` | Health check de l'API |
| `GET` | `/api/v1/eligibility?phone=257364` | Vérifier l'éligibilité (GET) |
| `POST` | `/api/v1/eligibility` | Vérifier l'éligibilité (POST) |
| `GET` | `/swagger/` | Documentation Swagger UI |

## 📋 Codes HTTP

| Code | Description | Cas d'usage |
|------|-------------|-------------|
| `200 OK` | Succès | Numéro trouvé et vérifié |
| `400 Bad Request` | Erreur validation | Format invalide ou paramètre manquant |
| `404 Not Found` | Introuvable | Numéro inexistant dans la base OPT |
| `405 Method Not Allowed` | Méthode invalide | Seuls GET et POST acceptés |

## 🔧 Options de démarrage

### Port personnalisé

```bash
# Avec Podman
podman run -d -p 9090:8080 --name opt-api rastadidi/optnc-fiber-eligibility

# L'API sera accessible sur le port 9090
http GET :9090/health
```

### Mode CLI (sans serveur)

```bash
# Vérifier un numéro en mode CLI
podman run --rm rastadidi/optnc-fiber-eligibility 257364

# Format JSON
podman run --rm rastadidi/optnc-fiber-eligibility 257364 --json
```

### Variables d'environnement

```bash
podman run -d -p 8080:8080 \
  -e CHROMIUM_PATH=/usr/bin/chromium-browser \
  --name opt-api \
  rastadidi/optnc-fiber-eligibility
```

## 🛠️ Gestion du conteneur

```bash
# Voir les logs
podman logs opt-api

# Suivre les logs en temps réel
podman logs -f opt-api

# Arrêter le conteneur
podman stop opt-api

# Redémarrer le conteneur
podman restart opt-api

# Supprimer le conteneur
podman rm opt-api
```

## 📊 Informations techniques

- **Base :** Alpine Linux 3.19
- **Langage :** Go 1.21
- **Taille :** 675 MB
- **Architecture :** amd64
- **Port :** 8080
- **User :** Non-root (app)

## 🔗 Liens

- **GitHub :** https://github.com/adriens/optnc-fiber-eligibility
- **Documentation complète :** https://github.com/adriens/optnc-fiber-eligibility#readme
- **Issues :** https://github.com/adriens/optnc-fiber-eligibility/issues

## 📝 License

MIT License

---

**Made with ❤️ in New Caledonia 🇳🇨**
