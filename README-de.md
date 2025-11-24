# Link Checker für Go
[![en](https://img.shields.io/badge/lang-en-red.svg)](https://github.com/sekika/linkchecker/blob/main/README.md)

`linkchecker` ist ein in Go geschriebenes Kommandozeilenwerkzeug zum Überprüfen von Links, die in einer angegebenen URL oder in einer lokalen HTML-Datei enthalten sind.

Die wichtigsten Funktionen sind:

* **Rücksichtsvolle Nebenläufigkeit:** Links werden mit mehreren Workern parallel überprüft, während ein konfigurierbarer Warteintervall (`-wait`) zwischen Anfragen an *denselben Host* strikt eingehalten wird. Dies verhindert eine unbeabsichtigte Überlastung oder ein versehentliches DoS.
* **Flexible Linkquellen:** Das Tool kann sowohl entfernte URLs als auch lokale HTML-Dateien durchsuchen.
* **Anpassbares Verhalten:** HTTP-Timeout, User-Agent und das Ignorieren interner Links oder bestimmter Hostnamen können flexibel eingestellt werden.

## Installation

Wenn Go auf Ihrem System installiert ist, können Sie das Tool mit folgendem Befehl installieren:

```bash
go install github.com/sekika/linkchecker/cmd/linkchecker@latest
```

## Verwendung

Nach der Installation kann das Tool mit dem Befehl `linkchecker` ausgeführt werden.

### Grundlegende Nutzung

Geben Sie die Ziel-URL oder den Pfad zu einer lokalen HTML-Datei mit der Option `-u` an.

```bash
# Links auf einer Website überprüfen
linkchecker -u https://example.com/page.html

# Links in einer lokalen Datei überprüfen
linkchecker -u path/to/local/file.html
```

### Ergebnisse filtern (nur Fehler anzeigen)

Da linkchecker für jeden Link `[OK]` oder `[NG]` ausgibt, können Sie fehlgeschlagene Links einfach mit `grep` herausfiltern:

```bash
linkchecker -u https://example.com/page.html | grep "\[NG\]"
```

### Optionen

| Option         | Beschreibung                                                                                  | Standardwert                  |
| -------------- | --------------------------------------------------------------------------------------------- | ----------------------------- |
| `-u`           | Ziel-URL oder lokale HTML-Datei (erforderlich)                                                | ""                            |
| `-no-internal` | Interne Links (gleicher Host/das gleiche Domain) nicht prüfen                                 | false                         |
| `-ignore`      | Pfad zu einer Datei mit Host-/Domainnamen, die ignoriert werden sollen                        | ""                            |
| `-timeout`     | HTTP-Timeout in Sekunden                                                                      | 10                            |
| `-wait`        | Wartezeit in Sekunden zwischen Anfragen an denselben Host. Steuert die Crawl-Geschwindigkeit. | 3                             |
| `-user-agent`  | User-Agent für HTTP-Anfragen                                                                  | github.com/sekika/linkchecker |

### Beispiele

Interne Links ausschließen und das Timeout auf 5 Sekunden setzen:

```bash
linkchecker -u https://example.com -no-internal -timeout 5
```

## Verwendung als Bibliothek (fortgeschritten)

Obwohl dieses Repository hauptsächlich für das Kommandozeilenwerkzeug gedacht ist, können Sie die Kernfunktionalität auch programmatisch nutzen.

### Kernfunktionalität importieren

Um die Link-Extraktion in einem Go-Programm zu verwenden, importieren Sie das `crawler`-Paket vom neuen öffentlichen Pfad:

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

    // Links aus einer URL extrahieren
    links, err := crawler.ExtractLinksFromURL(url, timeoutSec, userAgent)
    if err != nil {
        log.Fatalf("Fehler beim Extrahieren der Links: %v", err)
    }

    fmt.Printf("Es wurden %d Links auf %s gefunden\n", len(links), url)

    // Beispiel für das Ausführen von Workern
    // Hinweis: RunWorkers benötigt eine Liste absoluter Links
    // crawler.RunWorkers(links, url, false, make(map[string]bool), timeoutSec, 3, userAgent)
}
```

### Codeanalyse

* [How the Go-based link checker works (englisch)](https://sekika.github.io/2025/11/21/go-linkchecker/)
