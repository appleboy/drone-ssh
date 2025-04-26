# drone-ssh

> [English](./README.md) | **繁體中文** | [简体中文](./README.zh-cn.md)

![sshlog](images/ssh.png)

<!-- 圖片說明：SSH 日誌畫面，圖片內容與原文相同 -->

[![GitHub tag](https://img.shields.io/github/tag/appleboy/drone-ssh.svg)](https://github.com/appleboy/drone-ssh/releases)
[![GoDoc](https://godoc.org/github.com/appleboy/drone-ssh?status.svg)](https://godoc.org/github.com/appleboy/drone-ssh)
[![Lint and Testing](https://github.com/appleboy/drone-ssh/actions/workflows/testing.yml/badge.svg?branch=master)](https://github.com/appleboy/drone-ssh/actions/workflows/testing.yml)
[![codecov](https://codecov.io/gh/appleboy/drone-ssh/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-ssh)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-ssh)](https://goreportcard.com/report/github.com/appleboy/drone-ssh)
[![Docker Pulls](https://img.shields.io/docker/pulls/appleboy/drone-ssh.svg)](https://hub.docker.com/r/appleboy/drone-ssh/)

Drone 外掛程式，可透過 SSH 在遠端主機執行指令。使用方式與可用選項請參考[官方文件](http://plugins.drone.io/appleboy/drone-ssh/)。

**注意：請將 Drone 的 image config 路徑更新為 `appleboy/drone-ssh`。`plugins/ssh` 已不再維護。**

![demo](./images/demo2017.05.10.gif)

<!-- 圖片說明：SSH 指令執行示意動畫，內容與原文相同 -->

## 重大變更

`v1.5.0`：將指令逾時參數改為 `Duration` 格式。設定範例如下：

```diff
pipeline:
  scp:
    image: ghcr.io/appleboy/drone-ssh
    settings:
      host:
        - example1.com
        - example2.com
      username: ubuntu
      password:
        from_secret: ssh_password
      port: 22
-     command_timeout: 120
+     command_timeout: 2m
      script:
        - echo "Hello World"
```

## 建置或下載執行檔

可於[發行頁面](https://github.com/appleboy/drone-ssh/releases)下載預先編譯的執行檔，支援以下作業系統：

- Windows amd64/386
- Linux arm/amd64/386
- macOS (Darwin) amd64/386

若已安裝 `Go`，可執行：

```sh
go install github.com/appleboy/drone-ssh@latest
```

或使用下列指令自行建置執行檔：

```sh
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go test -cover ./...

go build -v -a -tags netgo -o release/linux/amd64/drone-ssh .
```

## Docker

可使用下列指令建置 Docker 映像檔：

```sh
make docker
```

## 使用方式

於工作目錄下執行：

```sh
docker run --rm \
  -e PLUGIN_HOST=foo.com \
  -e PLUGIN_USERNAME=root \
  -e PLUGIN_KEY="$(cat ${HOME}/.ssh/id_rsa)" \
  -e PLUGIN_SCRIPT=whoami \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  ghcr.io/appleboy/drone-ssh
```

## 以檔案路徑掛載金鑰

請確認已於專案設定中啟用 `trusted` 模式（適用於 [drone 0.8 版本](https://0-8-0.docs.drone.io/)）。

![trusted mode](./images/trust.png)

<!-- 圖片說明：Drone 專案 trusted 模式設定畫面 -->

於 `.drone.yml` 設定檔的 `volumes` 區段掛載私鑰：

```diff
pipeline:
  ssh:
    image: ghcr.io/appleboy/drone-ssh
    host: xxxxx.com
    username: deploy
+   volumes:
+     - /root/drone_rsa:/root/ssh/drone_rsa
    key_path: /root/ssh/drone_rsa
    script:
      - echo "test ssh"
```

詳情請參考 [issue comment](https://github.com/appleboy/drone-ssh/issues/51#issuecomment-336732928)。

## 設定說明

更多範例與完整設定選項請參考 [DOCS.md](./DOCS.md)

設定選項來源如下：

0. 內建 drone-ssh 預設值。詳見 [main.go CLI Flags](https://github.com/appleboy/drone-ssh/blob/6d9d6acc6aef1f9166118c6ba8bd214d3a582bdb/main.go#L39)。
1. 由 `PLUGIN_ENV_FILE` 環境變數指定的 dotenv 檔案。
2. `.drone.yml` Drone 設定檔。

後面的來源會覆蓋前面的設定，例如 `.env` 檔案中的 `PORT` 會覆蓋 main.go 的預設值。
