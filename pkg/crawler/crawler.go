package crawler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// ----------------------
// リンク抽出
// ----------------------

func ExtractLinksFromFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return extractLinks(string(data)), nil
}

func ExtractLinksFromURL(rawURL string, timeoutSec int, userAgent string) ([]string, error) {
	client := &http.Client{Timeout: time.Duration(timeoutSec) * time.Second}
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return extractLinks(string(body)), nil
}

func extractLinks(html string) []string {
	re := regexp.MustCompile(`href="([^"#]+)`)
	matches := re.FindAllStringSubmatch(html, -1)
	var links []string
	for _, m := range matches {
		links = append(links, m[1])
	}
	return links
}

// ----------------------
// HTTP リクエスト
// ----------------------

func FetchHTTP(link string, client *http.Client, userAgent string) error {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}

// ----------------------
// ワーカー処理
// ----------------------

func RunWorkers(
	links []string,
	baseURL string,
	noInternal bool,
	ignoreHosts map[string]bool,
	timeoutSec int,
	waitSec int,
	userAgent string,
) {
	base, _ := url.Parse(baseURL)
	filteredLinks := []string{}

	for _, link := range links {
		if noInternal && isInternalLinkRaw(link) {
			continue
		}

		u, err := base.Parse(link)
		if err != nil {
			continue
		}

		if ignoreHosts[u.Host] {
			continue
		}

		filteredLinks = append(filteredLinks, u.String())
	}

	if len(filteredLinks) == 0 {
		fmt.Println("No links to check.")
		return
	}

	// ホストごとの専用キューと goroutine
	hostQueues := make(map[string]chan string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, link := range filteredLinks {
		u, _ := url.Parse(link)
		host := u.Host

		mu.Lock()
		if _, ok := hostQueues[host]; !ok {
			hostQueues[host] = make(chan string, 100)
			ch := hostQueues[host]

			wg.Add(1)
			go func(host string, ch chan string) {
				defer wg.Done()
				jar, _ := cookiejar.New(nil)
				client := &http.Client{
					Timeout: time.Duration(timeoutSec) * time.Second,
					Jar:     jar,
				}
				for l := range ch {
					err := FetchHTTP(l, client, userAgent)
					if err != nil {
						fmt.Printf("[NG] %s (%v)\n", l, err)
					} else {
						fmt.Printf("[OK] %s\n", l)
					}
					time.Sleep(time.Duration(waitSec) * time.Second) // 同じホストの間隔
				}
			}(host, ch)
		}
		mu.Unlock()

		// リンクをホストキューに送信
		hostQueues[host] <- link
	}

	// 全てのキューを閉じる
	for _, ch := range hostQueues {
		close(ch)
	}

	// 全てのホストワーカーの終了を待つ
	wg.Wait()
}

// ----------------------
// 内部リンク判定
// ----------------------

func isInternalLinkRaw(link string) bool {
	u, err := url.Parse(link)
	if err != nil {
		return false
	}
	return !u.IsAbs()
}

// ----------------------
// ignoreHosts ファイル読み込み
// ----------------------

func LoadIgnoreHosts(filename string) (map[string]bool, error) {
	hosts := make(map[string]bool)
	data, err := os.ReadFile(filename)
	if err != nil {
		return hosts, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			hosts[line] = true
		}
	}
	return hosts, nil
}
