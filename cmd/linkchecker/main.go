/*
linkchecker is a command-line tool written in Go for checking link health
within a specified URL or a local HTML file. It manages request intervals
to the same host to prevent accidental overload.

Usage:

	linkchecker -u <URL or file> [options]

For detailed usage, options, and library information, please refer to the
module root's README.md file.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sekika/linkchecker/pkg/crawler"
)

func main() {
	urlFlag := flag.String("u", "", "URL or local HTML")
	noInternal := flag.Bool("no-internal", false, "Check no internal link")
	ignoreHostsFile := flag.String("ignore", "", "File of ignoring host")
	timeoutSec := flag.Int("timeout", 10, "Timeout second of HTTP request")
	waitSec := flag.Int("wait", 3, "Waiting second for the same host")
	userAgent := flag.String("user-agent", "github.com/sekika/linkchecker", "User-Agent for HTTP-request")
	flag.Parse()
	log.SetOutput(os.Stdout)

	if *urlFlag == "" {
		fmt.Println("Specify URL or filename")
		os.Exit(1)
	}

	baseURL := *urlFlag
	var links []string
	var err error

	if isLocalFile(baseURL) {
		links, err = crawler.ExtractLinksFromFile(baseURL)
		if err != nil {
			log.Fatalf("Failed to extract links from file: %v", err)
		}
		baseURL = "file://" + baseURL
	} else {
		links, err = crawler.ExtractLinksFromURL(baseURL, *timeoutSec, *userAgent)
		if err != nil {
			log.Fatalf("Failed to extract links from URL: %v", err)
		}
	}

	ignoreHosts := make(map[string]bool)
	if *ignoreHostsFile != "" {
		ignoreHosts, err = crawler.LoadIgnoreHosts(*ignoreHostsFile)
		if err != nil {
			log.Fatalf("Failed to load ignore hosts file: %v", err)
		}
	}

	crawler.RunWorkers(links, baseURL, *noInternal, ignoreHosts, *timeoutSec, *waitSec, *userAgent)
}

func isLocalFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
