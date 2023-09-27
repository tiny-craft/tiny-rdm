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

![](screenshots/light_en.png)

## Feature

* Built on Webview, no embedded browsers (Thanks to [Wails](https://github.com/wailsapp/wails)).
* More elegant UI and visualized layout (Thanks to [Naive UI](https://github.com/tusen-ai/naive-ui)
  and [IconPark](https://iconpark.oceanengine.com)).
* Multi-language support (Click here to contribute and support more languages).
* Convenient data viewing and editing.
* More features under continuous development...

## Installation

We publish binaries for Mac, Windows, and Linux.
Available to download for free from [here](https://github.com/tiny-craft/tiny-rdm/releases).

> If you can't open it after installation on macOS, exec the following command then reopen:
> ``` shell
>  sudo xattr -d com.apple.quarantine /Applications/Tiny\ RDM.app
> ```

## Build
### Prerequisites
* Go >= 1.21
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
