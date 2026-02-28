<div align="center">
<a href="https://github.com/tiny-craft/tiny-rdm/"><img src="build/appicon.png" width="120"/></a>
</div>
<h1 align="center">Tiny RDM</h1>
<h4 align="center"><a href="/">English</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_zh.md">简体中文</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_tw.md">繁體中文</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_ja.md">日本語</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_ko.md">한국어</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_fr.md">Français</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_es.md">Español</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_pt.md">Português (BR)</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_ru.md">Русский</a> | <strong>Türkçe</strong></h4>
<div align="center">

[![License](https://img.shields.io/github/license/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/blob/main/LICENSE)
[![GitHub release](https://img.shields.io/github/release/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/releases)
![GitHub All Releases](https://img.shields.io/github/downloads/tiny-craft/tiny-rdm/total)
[![GitHub stars](https://img.shields.io/github/stars/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/fork)

<strong>Tiny RDM, Mac, Windows ve Linux için kullanılabilen modern, hafif ve çapraz platform bir Redis masaüstü yöneticisidir. Docker ile dağıtılabilen bir web sürümü de sunmaktadır.</strong>
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

## Özellikler

* Ultra hafif, Webview2 tabanlı, gömülü tarayıcı yok ([Wails](https://github.com/wailsapp/wails)'e teşekkürler)
* Görsel ve kullanıcı dostu arayüz, açık ve koyu tema desteği ([Naive UI](https://github.com/tusen-ai/naive-ui) ve [IconPark](https://iconpark.oceanengine.com)'a teşekkürler)
* Çoklu dil desteği ([Daha fazla dil mi gerekiyor? Katkıda bulunmak için tıklayın](.github/CONTRIBUTING.md))
* Gelişmiş bağlantı yönetimi: SSH Tüneli/SSL/Sentinel Modu/Cluster Modu/HTTP proxy/SOCKS5 proxy desteği
* Anahtar-değer işlemlerinin görselleştirilmesi, List, Hash, String, Set, Sorted Set ve Stream için CRUD desteği
* Çoklu veri görüntüleme formatı ve çözme/sıkıştırma açma yöntemleri desteği
* SCAN ile segmentli yükleme, milyonlarca anahtarı kolayca listeleme
* Komut işlem geçmişi günlük listesi
* Komut satırı modu
* Yavaş günlük listesi
* List/Hash/Set/Sorted Set için segmentli yükleme ve sorgulama
* List/Hash/Set/Sorted Set değerleri için çözme/sıkıştırma açma
* Monaco Editor entegrasyonu
* Gerçek zamanlı komut izleme desteği
* Veri içe/dışa aktarma desteği
* Yayınla/abone ol desteği
* Bağlantı profili içe/dışa aktarma desteği
* Değer görüntüleme için özel veri kodlayıcı ve çözücü ([Talimatlar burada](https://redis.tinycraft.cc/guide/custom-decoder/))

## Kurulum

[Buradan](https://github.com/tiny-craft/tiny-rdm/releases) ücretsiz olarak indirilebilir.

> macOS'ta kurulumdan sonra açamıyorsanız, aşağıdaki komutu çalıştırıp tekrar açın:
> ``` shell
>  sudo xattr -d com.apple.quarantine /Applications/Tiny\ RDM.app
> ```

## Derleme Kılavuzu

### Gereksinimler

* Go (en son sürüm)
* Node.js >= 20
* NPM >= 9

### Wails Kurulumu

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Kodu Çekme

```bash
git clone https://github.com/tiny-craft/tiny-rdm --depth=1
```

### Frontend Derleme

```bash
npm install --prefix ./frontend
```

veya

```bash
cd frontend
npm install
```

### Derleme ve Çalıştırma

```bash
wails dev
```

## Docker ile Dağıtım

Masaüstü istemcisinin yanı sıra, Tiny RDM Docker ile hızlıca dağıtılabilen bir web sürümü de sunmaktadır.

### Docker Compose Kullanımı (önerilen)

Bir `docker-compose.yml` dosyası oluşturun:

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

Servisi başlatın:

```bash
docker compose up -d
```

Başlatıldıktan sonra `http://localhost:8086` adresini ziyaret edin ve yukarıda yapılandırılan kimlik bilgileriyle giriş yapın.

### Docker Komutu Kullanımı

```bash
docker run -d --name tinyrdm \
  -p 8086:8086 \
  -e ADMIN_USERNAME=admin \
  -e ADMIN_PASSWORD=tinyrdm \
  -v ./data:/app/tinyrdm \
  ghcr.io/tiny-craft/tiny-rdm:latest
```

### Ortam Değişkenleri

| Değişken | Açıklama | Varsayılan |
|----------|----------|------------|
| `ADMIN_USERNAME` | Giriş kullanıcı adı | - |
| `ADMIN_PASSWORD` | Giriş şifresi | - |

## Hakkında

### Sponsor

Bu proje işinize yaradıysa, bir kahve ısmarlayabilirsiniz ☕️

* Wechat Sponsor

<img src="docs/images/wechat_sponsor.jpg" alt="wechat" width="200" />
