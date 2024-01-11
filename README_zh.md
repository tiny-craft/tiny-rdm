<div align="center">
<a href="https://github.com/tiny-craft/tiny-rdm/"><img src="build/appicon.png" width="120"/></a>
</div>
<h1 align="center">Tiny RDM</h1>
<h4 align="center"><strong><a href="/">English</a></strong> | 简体中文</h4>
<div align="center">

[![License](https://img.shields.io/github/license/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/blob/main/LICENSE)
[![GitHub release](https://img.shields.io/github/release/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/releases)
![GitHub All Releases](https://img.shields.io/github/downloads/tiny-craft/tiny-rdm/total)
[![GitHub stars](https://img.shields.io/github/stars/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/fork)

<strong>一个现代化轻量级的跨平台Redis桌面客户端，支持Mac、Windows和Linux</strong>
</div>

<picture>
 <source media="(prefers-color-scheme: dark)" srcset="screenshots/dark_zh.png">
 <source media="(prefers-color-scheme: light)" srcset="screenshots/light_zh.png">
 <img alt="screenshot" src="screenshots/dark_zh.png">
</picture>

## 功能特性

* 极度轻量，基于Webview2，无内嵌浏览器（感谢[Wails](https://github.com/wailsapp/wails)）
* 界面精美易用，提供浅色/深色主题（感谢[Naive UI](https://github.com/tusen-ai/naive-ui)
  和 [IconPark](https://iconpark.oceanengine.com)）
* 多国语言支持：英文/中文（[需要更多语言支持？点我贡献语言](.github/CONTRIBUTING_zh.md)）
* 更好用的连接管理：支持SSH隧道/SSL/哨兵模式/集群模式
* 可视化键值操作，增删查改一应俱全
* 支持多种数据查看格式以及转码/解压方式
* 采用SCAN分段加载，可轻松处理数百万键列表
* 操作命令执行日志展示
* 提供命令行操作
* 提供慢日志展示
* List/Hash/Set/Sorted Set的分段加载和查询
* List/Hash/Set/Sorted Set值的转码显示
* 内置高级编辑器Monaco Editor
* 支持命令实时监控
* 支持数据导入/导出

## 未来版本规划

- [ ] 发布/订阅支持
- [ ] 连接配置导入/导出

## 安装

提供Mac、Windows和Linux安装包，可[免费下载](https://github.com/tiny-craft/tiny-rdm/releases)。

> 如果在macOS上安装后无法打开，报错**不受信任**或者**移到垃圾箱**，执行下面命令后再启动即可：
> ``` shell
>  sudo xattr -d com.apple.quarantine /Applications/Tiny\ RDM.app
> ```

## 构建项目

### 运行环境要求

* Go（最新版本）
* Node.js >= 16
* NPM >= 9

### 安装wails

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 拉取代码

```bash
git clone https://github.com/tiny-craft/tiny-rdm --depth=1
```

### 构建前端代码

```bash
npm install --prefix ./frontend
```

### 编译运行开发版本

```bash
wails dev
```

## 关于

如果你也同为独立开发者（团队），喜欢开源，或者对Tiny Craft的相关产品感兴趣，可以关注微信公众号或者加入QQ群，探讨心得，反馈意见，交个朋友。

### 微信公众号（用户交流微信群）

我会不定期更新一些关于独立开发的思考和感悟，以及独立产品的介绍，欢迎扫码关注~👏

<img src="docs/images/wechat_official.png" alt="wechat" width="360" />

### 独立开发互助QQ群

```
831077639
```

### 赞助

该项目完全为爱发电，如果对你有所帮助，可以请作者喝杯咖啡 ☕️

* 微信赞赏

<img src="docs/images/wechat_sponsor.jpg" alt="wechat" width="200" />
