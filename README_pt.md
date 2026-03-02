<div align="center">
<a href="https://github.com/tiny-craft/tiny-rdm/"><img src="build/appicon.png" width="120"/></a>
</div>
<h1 align="center">Tiny RDM</h1>
<h4 align="center"><a href="/">English</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_zh.md">简体中文</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_tw.md">繁體中文</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_ja.md">日本語</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_ko.md">한국어</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_fr.md">Français</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_es.md">Español</a> | <strong>Português (BR)</strong> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_ru.md">Русский</a> | <a href="https://github.com/tiny-craft/tiny-rdm/blob/main/README_tr.md">Türkçe</a></h4>
<div align="center">

[![License](https://img.shields.io/github/license/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/blob/main/LICENSE)
[![GitHub release](https://img.shields.io/github/release/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/releases)
![GitHub All Releases](https://img.shields.io/github/downloads/tiny-craft/tiny-rdm/total)
[![GitHub stars](https://img.shields.io/github/stars/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/tiny-craft/tiny-rdm)](https://github.com/tiny-craft/tiny-rdm/fork)

<strong>Tiny RDM é um gerenciador Redis moderno, leve e multiplataforma, disponível para Mac, Windows e Linux. Também oferece uma versão web que pode ser implantada via Docker.</strong>
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

## Funcionalidades

* Ultra leve, baseado em Webview2, sem navegador embutido (Graças ao [Wails](https://github.com/wailsapp/wails))
* Interface visual e amigável, temas claro e escuro (Graças ao [Naive UI](https://github.com/tusen-ai/naive-ui) e [IconPark](https://iconpark.oceanengine.com))
* Suporte multilíngue ([Precisa de mais idiomas? Clique aqui para contribuir](.github/CONTRIBUTING.md))
* Gerenciamento aprimorado de conexões: túnel SSH/SSL/modo Sentinel/modo Cluster/proxy HTTP/proxy SOCKS5
* Visualização de operações chave-valor, suporte CRUD para List, Hash, String, Set, Sorted Set e Stream
* Suporte a múltiplos formatos de visualização e métodos de decodificação/descompressão
* Carregamento segmentado com SCAN para listar facilmente milhões de chaves
* Lista de logs do histórico de comandos
* Modo linha de comando
* Lista de logs lentos
* Carregamento segmentado e consultas para List/Hash/Set/Sorted Set
* Decodificação/descompressão de valores para List/Hash/Set/Sorted Set
* Integração com Monaco Editor
* Monitoramento de comandos em tempo real
* Importação/exportação de dados
* Publicação/assinatura
* Importação/exportação de perfis de conexão
* Codificador e decodificador de dados personalizados para exibição de valores ([Instruções aqui](https://tinyrdm.com/guide/custom-decoder/))

## Instalação

Disponível para download gratuito [aqui](https://github.com/tiny-craft/tiny-rdm/releases).

> Se não conseguir abrir após a instalação no macOS, execute o seguinte comando e reabra:
> ``` shell
>  sudo xattr -d com.apple.quarantine /Applications/Tiny\ RDM.app
> ```

## Guia de compilação

### Pré-requisitos

* Go (versão mais recente)
* Node.js >= 20
* NPM >= 9

### Instalar Wails

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Obter o código

```bash
git clone https://github.com/tiny-craft/tiny-rdm --depth=1
```

### Compilar o frontend

```bash
npm install --prefix ./frontend
```

ou

```bash
cd frontend
npm install
```

### Compilar e executar

```bash
wails dev
```

## Implantação com Docker

Além do cliente desktop, o Tiny RDM também oferece uma versão web que pode ser implantada rapidamente via Docker.

### Usando Docker Compose (recomendado)

Crie um arquivo `docker-compose.yml`:

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

Inicie o serviço:

```bash
docker compose up -d
```

Após iniciar, acesse `http://localhost:8086` e faça login com as credenciais configuradas acima.

### Usando o comando Docker

```bash
docker run -d --name tinyrdm \
  -p 8086:8086 \
  -e ADMIN_USERNAME=admin \
  -e ADMIN_PASSWORD=tinyrdm \
  -v ./data:/app/tinyrdm \
  ghcr.io/tiny-craft/tiny-rdm:latest
```

### Variáveis de ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `ADMIN_USERNAME` | Nome de usuário | - |
| `ADMIN_PASSWORD` | Senha | - |

## Sobre

### Patrocinar

Se este projeto foi útil para você, sinta-se à vontade para pagar um café ☕️

* Wechat Sponsor

<img src="docs/images/wechat_sponsor.jpg" alt="wechat" width="200" />
