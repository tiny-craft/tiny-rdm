<div align="center">
<a href="https://github.com/tiny-craft/tiny-rdm/"><img src="build/appicon.png" width="120"/></a>
</div>
<h1 align="center">Tiny RDM</h1>
<h4 align="center"><a href="/">English</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_zh.md">简体中文</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_tw.md">繁體中文</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_ja.md">日本語</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_ko.md">한국어</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_fr.md">Français</a> | <strong>Español</strong> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_pt.md">Português (BR)</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_ru.md">Русский</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_tr.md">Türkçe</a></h4>
<div align="center">

[![License](https://img.shields.io/github/license/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/blob/main/LICENSE)
[![GitHub release](https://img.shields.io/github/release/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/releases)
![GitHub All Releases](https://img.shields.io/github/downloads/tiny-craft/tiny-rdm/total)
[![GitHub stars](https://img.shields.io/github/stars/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/fork)

<strong>Tiny RDM es un gestor Redis moderno, ligero y multiplataforma, disponible para Mac, Windows y Linux. También ofrece una versión web que se puede desplegar mediante Docker.</strong>
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

## Características

* Ultra ligero, basado en Webview2, sin navegador integrado (Gracias a [Wails](https://github.com/wailsapp/wails))
* Interfaz visual y fácil de usar, temas claro y oscuro (Gracias a [Naive UI](https://github.com/tusen-ai/naive-ui) e [IconPark](https://iconpark.oceanengine.com))
* Soporte multilingüe ([¿Necesitas más idiomas? Haz clic aquí para contribuir](.github/CONTRIBUTING.md))
* Gestión mejorada de conexiones: túnel SSH/SSL/modo Sentinel/modo Cluster/proxy HTTP/proxy SOCKS5
* Visualización de operaciones clave-valor, soporte CRUD para List, Hash, String, Set, Sorted Set y Stream
* Soporte de múltiples formatos de visualización y métodos de decodificación/descompresión
* Carga segmentada con SCAN para listar fácilmente millones de claves
* Lista de registros del historial de comandos
* Modo línea de comandos
* Lista de registros lentos
* Carga segmentada y consultas para List/Hash/Set/Sorted Set
* Decodificación/descompresión de valores para List/Hash/Set/Sorted Set
* Integración con Monaco Editor
* Monitoreo de comandos en tiempo real
* Importación/exportación de datos
* Publicación/suscripción
* Importación/exportación de perfiles de conexión
* Codificador y decodificador de datos personalizados para la visualización de valores ([Instrucciones aquí](https://redis.tinycraft.cc/guide/custom-decoder/))

## Instalación

Disponible para descargar gratis [aquí](https://github.com/tiny-craft/tiny-rdm/releases).

> Si no puedes abrirlo después de la instalación en macOS, ejecuta el siguiente comando y vuelve a abrirlo:
> ``` shell
>  sudo xattr -d com.apple.quarantine /Applications/Tiny\ RDM.app
> ```

## Guía de compilación

### Requisitos previos

* Go (última versión)
* Node.js >= 20
* NPM >= 9

### Instalar Wails

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Obtener el código

```bash
git clone https://github.com/tiny-craft/tiny-rdm --depth=1
```

### Compilar el frontend

```bash
npm install --prefix ./frontend
```

o

```bash
cd frontend
npm install
```

### Compilar y ejecutar

```bash
wails dev
```

## Despliegue con Docker

Además del cliente de escritorio, Tiny RDM también ofrece una versión web que se puede desplegar rápidamente con Docker.

### Usando Docker Compose (recomendado)

Crea un archivo `docker-compose.yml`:

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

Inicia el servicio:

```bash
docker compose up -d
```

Una vez iniciado, visita `http://localhost:8086` e inicia sesión con las credenciales configuradas arriba.

### Usando el comando Docker

```bash
docker run -d --name tinyrdm \
  -p 8086:8086 \
  -e ADMIN_USERNAME=admin \
  -e ADMIN_PASSWORD=tinyrdm \
  -v ./data:/app/tinyrdm \
  ghcr.io/tiny-craft/tiny-rdm:latest
```

### Variables de entorno

| Variable | Descripción | Valor por defecto |
|----------|-------------|-------------------|
| `ADMIN_USERNAME` | Nombre de usuario | - |
| `ADMIN_PASSWORD` | Contraseña | - |

## Acerca de

### Patrocinar

Si este proyecto te resulta útil, no dudes en invitar al autor a un café ☕️

* Wechat Sponsor

<img src="docs/images/wechat_sponsor.jpg" alt="wechat" width="200" />
