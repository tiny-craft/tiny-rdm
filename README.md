<h4 align="right"><strong>English</strong> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_zh.md">简体中文</a></h4>
<div align="center">
<a href="https://github.com/tiny-craft/tiny-rdm/"><img src="build/appicon.png" width="120"/></a>
</div>
<h1 align="center">Tiny RDM</h1>
<div align="center">

[![License](https://img.shields.io/github/license/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/blob/main/LICENSE)
[![GitHub release](https://img.shields.io/github/release/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/releases)
[![GitHub stars](https://img.shields.io/github/stars/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/fork)

<strong>Tiny RDM is a modern lightweight cross-platform Redis desktop manager available for Mac, Windows, and Linux.</strong>
</div>

![](screenshots/dark_en.png)

## Feature

* Super lightweight, built on Webview2, without embedded browsers (Thanks to [Wails](https://github.com/wailsapp/wails)).
* More elegant UI, frameless, offering light and dark themes (Thanks to [Naive UI](https://github.com/tusen-ai/naive-ui)
  and [IconPark](https://iconpark.oceanengine.com)).
* Multi-language support ([Need more languages ? Click here to contribute](.github/CONTRIBUTING.md)).
* Better connection management: supports SSH Tunnel/SSL/Sentinel Mode/Cluster Mode.
* Visualize key value operations, CRUD support for Lists, Hashes, Strings, Sets, Sorted Sets, and Streams.
* Support multiple data viewing format and decode/decompression methods.
* Use SCAN for segmented loading, making it easy to list millions of keys.
* Operation command execution logs.
* Provides command-line operations.
* Provides slow logs.

## Roadmap
- [ ] Pagination and querying for List/Hash/Set/Sorted Set
- [ ] Decode/decompression display for value of List/Hash/Set/Sorted Set
- [ ] Real-time commands monitoring
- [ ] Pub/Sub operations
- [ ] Embedding Monaco Editor

## Installation

Available to download for free from [here](https://github.com/tiny-craft/tiny-rdm/releases).

> If you can't open it after installation on macOS, exec the following command then reopen:
> ``` shell
>  sudo xattr -d com.apple.quarantine /Applications/Tiny\ RDM.app
> ```

## Build Guidelines
### Prerequisites
* Go (latest version)
* Node.js >= 16
* NPM >= 9

### Install wails
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Clone the code
```bash
git clone https://github.com/tiny-craft/tiny-rdm --depth=1
```

### Build frontend
```bash
npm install --prefix ./frontend
```

### Compile and run
```bash
wails dev
```

## License

Tiny RDM is licensed under [GNU General Public](/LICENSE) license.
