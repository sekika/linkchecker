# linkchecker

`linkchecker` is a command-line tool written in Go for checking links found within a specified URL or a local HTML file.

Its key features include:

* **Concurrency with Respect:** It checks links concurrently using multiple workers while strictly observing a customizable wait period (`-wait`) between requests to the *same host*, preventing accidental DoS or excessive load on target servers.
* **Flexible Link Source:** It can crawl links from both remote URLs and local HTML files.
* **Customizable Behavior:** You can easily configure the HTTP timeout, User-Agent, and selectively ignore internal links or specific hostnames via an ignore file.

## Installation

If you have Go installed on your system, you can install the tool using the following command:

```bash
go install github.com/sekika/linkchecker/cmd/linkchecker@latest
````

## Usage

After installation, you can run the tool using the `linkchecker` command.

### Basic Use

Specify the target URL or local HTML file path using the `-u` flag.

```bash
# Check links on a website
linkchecker -u https://example.com/page.html

# Check links in a local file
linkchecker -u path/to/local/file.html
```

### Options

| Flag | Description | Default Value |
|---|---|---|
| `-u` | Target URL or local HTML file (Required) | "" |
| `-no-internal` | Do not check internal links (links within the same host/domain) | false |
| `-ignore` | Path to a file containing a list of hosts/domains to ignore | "" |
| `-timeout` | HTTP request timeout in seconds | 10 |
| `-wait` | Wait time in seconds between requests to the same host. Controls the crawling interval. | 3 |
| `-user-agent` | User-Agent string to use for HTTP requests | github.com/sekika/linkchecker |

### Examples

  - Exclude internal links and set the timeout to 5 seconds:

<!-- end list -->

```bash
linkchecker -u https://example.com -no-internal -timeout 5
```

## As a Library (Advanced)

Although this repository's primary focus is the command-line tool, you can use the core functionality by importing the relevant package.

### Importing Core Functionality

To use the link extraction logic programmatically, you can import the `crawler` package from the new public path:

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

    // Extract links from a URL
    links, err := crawler.ExtractLinksFromURL(url, timeoutSec, userAgent)
    if err != nil {
        log.Fatalf("Error extracting links: %v", err)
    }

    fmt.Printf("Found %d links on %s\n", len(links), url)

    // Example of running workers (Note: RunWorkers needs a list of absolute links)
    // crawler.RunWorkers(links, url, false, make(map[string]bool), timeoutSec, 3, userAgent)
}
