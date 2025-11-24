# Go 製リンクチェッカー
[![en](https://img.shields.io/badge/lang-en-red.svg)](https://github.com/sekika/linkchecker/blob/main/README.md)

`linkchecker` は、指定した URL またはローカル HTML ファイル内のリンクをチェックするための、Go で書かれたコマンドラインツールです。

主な特徴は次のとおりです：

* **配慮ある並行処理**：複数ワーカーによる並行チェックを行いながら、同一ホストへのリクエスト間隔を `-wait` で厳密に制御し、意図せぬ DoS や過剰な負荷を避けます。
* **柔軟なリンク取得元**：リモート URL とローカル HTML ファイルの両方をクロール可能です。
* **動作のカスタマイズ性**：HTTP タイムアウト、User-Agent、内部リンク無視、特定ホストの無視設定（ignore ファイル）などを柔軟に調整できます。

## インストール

システムに Go がインストールされていれば、以下のコマンドでツールをインストールできます：

```bash
go install github.com/sekika/linkchecker/cmd/linkchecker@latest
```

## 使い方

インストール後は、`linkchecker` コマンドでツールを実行できます。

### 基本的な使用方法

`-u` フラグで対象とする URL もしくはローカル HTML ファイルを指定します。

```bash
# Web サイト内のリンクをチェック
linkchecker -u https://example.com/page.html

# ローカルファイル内のリンクをチェック
linkchecker -u path/to/local/file.html
```

### 結果のフィルタ（NG のみ表示）

linkchecker は各リンクに対して `[OK]` または `[NG]` を出力するため、`grep` で NG のみ抽出できます：

```bash
linkchecker -u https://example.com/page.html | grep "\[NG\]"
```

### オプション一覧

| フラグ            | 説明                           | デフォルト値                        |
| -------------- | ---------------------------- | ----------------------------- |
| `-u`           | 対象 URL またはローカル HTML ファイル（必須） | ""                            |
| `-no-internal` | 同一ホスト・同一ドメインの内部リンクをチェックしない   | false                         |
| `-ignore`      | 無視するホスト名・ドメインを列挙したファイルのパス    | ""                            |
| `-timeout`     | HTTP リクエストのタイムアウト（秒）         | 10                            |
| `-wait`        | 同一ホストへのリクエスト間隔（秒）。クロール速度を制御  | 3                             |
| `-user-agent`  | HTTP リクエストに使用する User-Agent   | github.com/sekika/linkchecker |

### 使用例

内部リンクを除外し、タイムアウトを 5 秒に設定：

```bash
linkchecker -u https://example.com -no-internal -timeout 5
```

## ライブラリとして利用する（中級〜上級向け）

このリポジトリは主にコマンドラインツールを目的としていますが、リンク抽出などのコア機能はパッケージとして利用できます。

### コア機能をインポートする

プログラム内でリンク抽出機能を使用するには、新しい公開パスから `crawler` パッケージをインポートします：

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

    // URL からリンクを抽出
    links, err := crawler.ExtractLinksFromURL(url, timeoutSec, userAgent)
    if err != nil {
        log.Fatalf("Error extracting links: %v", err)
    }

    fmt.Printf("Found %d links on %s\n", len(links), url)

    // ワーカー実行の例（注：RunWorkers は絶対 URL のリストを必要とします）
    // crawler.RunWorkers(links, url, false, make(map[string]bool), timeoutSec, 3, userAgent)
}
```

## コード解析

* [Go で作るリンクチェッカーの仕組み](https://sekika.github.io/2025/11/19/go-linkchecker/)
