<div align="center">
<a href="https://github.com/tiny-craft/tiny-rdm/"><img src="build/appicon.png" width="120"/></a>
</div>
<h1 align="center">Tiny RDM</h1>
<h4 align="center"><strong><a href="/">English</a></strong> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_zh.md">简体中文</a> | 日本語</h4>
<div align="center">

[![License](https://img.shields.io/github/license/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/blob/main/LICENSE)
[![GitHub release](https://img.shields.io/github/release/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/releases)
![GitHub All Releases](https://img.shields.io/github/downloads/tiny-craft/tiny-rdm/total)
[![GitHub stars](https://img.shields.io/github/stars/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/fork)
[![Discord](https://img.shields.io/discord/1170373259133456434?label=Discord&color=5865F2)](https://discord.gg/VTFbBMGjWh)
[![X](https://img.shields.io/badge/Twitter-black?logo=x&logoColor=white)](https://twitter.com/Lykin53448)

<strong>Tiny RDMは、Mac、Windows、Linuxで利用可能な、モダンで軽量なクロスプラットフォームのRedisデスクトップマネージャーです。</strong>
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

## 特徴

* 超軽量、Webview2をベースにしており、埋め込みブラウザなし（[Wails](https://github.com/wailsapp/wails)に感謝）。
* 視覚的でユーザーフレンドリーなUI、ライトとダークテーマを提供（[Naive UI](https://github.com/tusen-ai/naive-ui)と[IconPark](https://iconpark.oceanengine.com)に感謝）。
* 多言語サポート（[もっと多くの言語が必要ですか？ここをクリックして貢献してください](.github/CONTRIBUTING.md)）。
* より良い接続管理：SSHトンネル/SSL/センチネルモード/クラスターモード/HTTPプロキシ/SOCKS5プロキシをサポート。
* キー値操作の可視化、リスト、ハッシュ、文字列、セット、ソートセット、ストリームのCRUDサポート。
* 複数のデータ表示形式とデコード/解凍方法をサポート。
* SCANを使用してセグメント化された読み込みを行い、数百万のキーを簡単にリスト化。
* コマンド操作履歴のログリスト。
* コマンドラインモードを提供。
* スローログリストを提供。
* リスト/ハッシュ/セット/ソートセットのセグメント化された読み込みとクエリ。
* リスト/ハッシュ/セット/ソートセットの値のデコード/解凍を提供。
* Monaco Editorと統合。
* リアルタイムコマンド監視をサポート。
* データのインポート/エクスポートをサポート。
* パブリッシュ/サブスクライブをサポート。
* 接続プロファイルのインポート/エクスポートをサポート。
* 値表示のためのカスタムデータエンコーダーとデコーダーをサポート（[こちらが手順です](https://redis.tinycraft.cc/guide/custom-decoder/)）。

## インストール

[こちら](https://github.com/tiny-craft/tiny-rdm/releases)から無料でダウンロードできます。

> macOSにインストール後に開けない場合は、以下のコマンドを実行してから再度開いてください：
> ``` shell
>  sudo xattr -d com.apple.quarantine /Applications/Tiny\ RDM.app
> ```

## ビルドガイドライン

### 前提条件

* Go（最新バージョン）
* Node.js >= 20
* NPM >= 9

### Wailsのインストール

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### コードの取得

```bash
git clone https://github.com/tiny-craft/tiny-rdm --depth=1
```

### フロントエンドのビルド

```bash
npm install --prefix ./frontend
```

または

```bash
cd frontend
npm install
```

### コンパイルと実行

```bash
wails dev
```
## について

### Wechat公式アカウント

<img src="docs/images/wechat_official.png" alt="wechat" width="360" />

### スポンサー

このプロジェクトが役立つ場合は、コーヒーを一杯おごってください ☕️。

* Wechatスポンサー

<img src="docs/images/wechat_sponsor.jpg" alt="wechat" width="200" />
