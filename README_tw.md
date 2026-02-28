<div align="center">
<a href="https://github.com/tiny-craft/tiny-rdm/"><img src="build/appicon.png" width="120"/></a>
</div>
<h1 align="center">Tiny RDM</h1>
<h4 align="center"><a href="/">English</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_zh.md">简体中文</a> | <strong>繁體中文</strong> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_ja.md">日本語</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_ko.md">한국어</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_fr.md">Français</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_es.md">Español</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_pt.md">Português (BR)</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_ru.md">Русский</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_tr.md">Türkçe</a></h4>
<div align="center">

[![License](https://img.shields.io/github/license/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/blob/main/LICENSE)
[![GitHub release](https://img.shields.io/github/release/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/releases)
![GitHub All Releases](https://img.shields.io/github/downloads/tiny-craft/tiny-rdm/total)
[![GitHub stars](https://img.shields.io/github/stars/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/fork)

<strong>Tiny RDM 是一款現代化輕量級的跨平台 Redis 桌面管理工具，支援 Mac、Windows 和 Linux</strong>
</div>

<picture>
 <source media="(prefers-color-scheme: dark)" srcset="screenshots/dark_en.png">
 <source media="(prefers-color-scheme: light)" srcset="screenshots/light_en.png">
 <img alt="screenshot" src="screenshots/dark_en.png">
</picture>

<picture>
 <source media="(prefers-color-scheme: dark)" srcset="screenshots/dark_en2.png">
 <source media="(prefers-color-scheme: light)" srcset="screenshots/light_en2.png">
 <img alt="screenshot" src="screenshots/dark_en2.png">
</picture>

## 功能特性

* 極度輕量，基於 Webview2，無內嵌瀏覽器（感謝 [Wails](https://github.com/wailsapp/wails)）
* 介面精美易用，提供淺色/深色主題（感謝 [Naive UI](https://github.com/tusen-ai/naive-ui) 和 [IconPark](https://iconpark.oceanengine.com)）
* 多國語言支援（[需要更多語言支援？點此貢獻](.github/CONTRIBUTING.md)）
* 更好的連線管理：支援 SSH 隧道/SSL/哨兵模式/叢集模式/HTTP 代理/SOCKS5 代理
* 視覺化鍵值操作，支援 List、Hash、String、Set、Sorted Set 和 Stream 的 CRUD
* 支援多種資料檢視格式及轉碼/解壓方式
* 採用 SCAN 分段載入，可輕鬆處理數百萬鍵列表
* 操作命令執行日誌展示
* 提供命令列模式
* 提供慢日誌展示
* List/Hash/Set/Sorted Set 的分段載入和查詢
* List/Hash/Set/Sorted Set 值的轉碼顯示
* 內建高級編輯器 Monaco Editor
* 支援命令即時監控
* 支援匯入/匯出資料
* 支援發布訂閱
* 支援匯入/匯出連線設定
* 自訂資料展示編碼/解碼（[操作指引](https://redis.tinycraft.cc/guide/custom-decoder/)）

## 安裝

提供 Mac、Windows 和 Linux 安裝包，可[免費下載](https://github.com/tiny-craft/tiny-rdm/releases)。

> 如果在 macOS 上安裝後無法開啟，出現**不受信任**或**移到垃圾桶**的錯誤，執行以下命令後再啟動即可：
> ``` shell
>  sudo xattr -d com.apple.quarantine /Applications/Tiny\ RDM.app
> ```

## 建置專案

### 環境需求

* Go（最新版本）
* Node.js >= 20
* NPM >= 9

### 安裝 Wails

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 取得程式碼

```bash
git clone https://github.com/tiny-craft/tiny-rdm --depth=1
```

### 建置前端

```bash
npm install --prefix ./frontend
```

或

```bash
cd frontend
npm install
```

### 編譯並執行

```bash
wails dev
```

## Docker 部署

除桌面客戶端外，Tiny RDM 還提供 Web 版本，可透過 Docker 快速部署。

### 使用 Docker Compose（推薦）

建立 `docker-compose.yml` 檔案：

```yaml
services:
  tinyrdm:
    image: ghcr.io/tiny-craft/tiny-rdm:latest
    container_name: tinyrdm
    restart: unless-stopped
    ports:
      - "8086:8086"
    environment:
      - ADMIN_USERNAME=admin
      - ADMIN_PASSWORD=tinyrdm
    volumes:
      - ./data:/app/tinyrdm
```

啟動服務：

```bash
docker compose up -d
```

啟動後造訪 `http://localhost:8086`，使用上方設定的帳號密碼登入。

### 使用 Docker 命令

```bash
docker run -d --name tinyrdm \
  -p 8086:8086 \
  -e ADMIN_USERNAME=admin \
  -e ADMIN_PASSWORD=tinyrdm \
  -v ./data:/app/tinyrdm \
  ghcr.io/tiny-craft/tiny-rdm:latest
```

### 環境變數說明

| 變數 | 說明 | 預設值 |
|------|------|--------|
| `ADMIN_USERNAME` | 登入帳號 | - |
| `ADMIN_PASSWORD` | 登入密碼 | - |
| `PORT` | Go 後端監聽埠 | `8088` |

## 關於

### 贊助

如果此專案對您有幫助，歡迎請作者喝杯咖啡 ☕️

* 微信贊賞

<img src="docs/images/wechat_sponsor.jpg" alt="wechat" width="200" />
