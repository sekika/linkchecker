# Link Checker pour Go
[![en](https://img.shields.io/badge/lang-en-red.svg)](https://github.com/sekika/linkchecker/blob/main/README.md)

`linkchecker` est un outil en ligne de commande écrit en Go permettant de vérifier les liens présents dans une URL donnée ou dans un fichier HTML local.

Ses principales caractéristiques sont les suivantes :

* **Concurrence avec respect des serveurs :** Les liens sont vérifiés de manière concurrente à l’aide de plusieurs workers, tout en respectant strictement un délai personnalisable (`-wait`) entre les requêtes vers un *même hôte*, afin d’éviter toute surcharge accidentelle ou un DoS involontaire.
* **Source de liens flexible :** L’outil peut analyser aussi bien des URLs distantes que des fichiers HTML locaux.
* **Comportement personnalisable :** Vous pouvez configurer facilement le délai d’expiration HTTP, le User-Agent, ignorer les liens internes ou exclure certains domaines via un fichier d’exclusion.

## Installation

Si Go est installé sur votre système, vous pouvez installer l’outil avec la commande suivante :

```bash
go install github.com/sekika/linkchecker/cmd/linkchecker@latest
```

## Utilisation

Après l’installation, l’outil peut être exécuté avec la commande `linkchecker`.

### Utilisation de base

Spécifiez l’URL cible ou le chemin du fichier HTML local à l’aide de l’option `-u`.

```bash
# Vérifier les liens d’un site web
linkchecker -u https://example.com/page.html

# Vérifier les liens d’un fichier local
linkchecker -u path/to/local/file.html
```

### Filtrer les résultats (afficher uniquement les échecs)

Puisque linkchecker affiche `[OK]` ou `[NG]` pour chaque lien vérifié, vous pouvez facilement filtrer les liens en échec à l’aide de `grep` :

```bash
linkchecker -u https://example.com/page.html | grep "\[NG\]"
```

### Options

| Option         | Description                                                                                             | Valeur par défaut             |
| -------------- | ------------------------------------------------------------------------------------------------------- | ----------------------------- |
| `-u`           | URL cible ou fichier HTML local (obligatoire)                                                           | ""                            |
| `-no-internal` | Ne pas vérifier les liens internes (même hôte/domaine)                                                  | false                         |
| `-ignore`      | Chemin d’un fichier contenant la liste des hôtes/domaines à ignorer                                     | ""                            |
| `-timeout`     | Délai d’expiration de la requête HTTP (en secondes)                                                     | 10                            |
| `-wait`        | Temps d’attente (en secondes) entre deux requêtes vers le même hôte. Contrôle la vitesse d’exploration. | 3                             |
| `-user-agent`  | Chaîne User-Agent utilisée pour les requêtes HTTP                                                       | github.com/sekika/linkchecker |

### Exemples

Exclure les liens internes et définir un délai d’expiration de 5 secondes :

```bash
linkchecker -u https://example.com -no-internal -timeout 5
```

## Utilisation comme bibliothèque (avancé)

Bien que ce dépôt soit principalement destiné à l’outil en ligne de commande, vous pouvez réutiliser la logique principale en important le package adéquat.

### Importer les fonctionnalités principales

Pour utiliser l’extraction de liens dans un programme Go, importez le package `crawler` depuis le nouveau chemin public :

```go
package main

import (
    "fmt"
    "log"
    "time"

    "https://github.com/sekika/linkchecker/pkg/crawler"
)

func main() {
    url := "https://example.com"
    timeoutSec := 10
    userAgent := "MyCustomApp/1.0"

    // Extraire les liens depuis une URL
    links, err := crawler.ExtractLinksFromURL(url, timeoutSec, userAgent)
    if err != nil {
        log.Fatalf("Erreur lors de l’extraction des liens : %v", err)
    }

    fmt.Printf("Trouvé %d liens sur %s\n", len(links), url)

    // Exemple d’exécution des workers (note : RunWorkers nécessite une liste de liens absolus)
    // crawler.RunWorkers(links, url, false, make(map[string]bool), timeoutSec, 3, userAgent)
}
```

## Analyse du code

* [Fonctionnement du link checker en Go (en anglais)](https://sekika.github.io/2025/11/21/go-linkchecker/)
