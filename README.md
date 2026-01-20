# libreoffice-as-a-service (LaaS)

Convert documents through libreoffice (soffice) as a service

## Table of Contents

- Install, Configure, Run
- System Requirements for Linux
- How to test locally on macOS
- How to Deploy to Digital Ocean

## Run LaaS

### Download

```bash
git clone git@github.com:savvi-legal/libreoffice-as-a-service.git
pushd ./libreoffice-as-a-service/
```

### Configure

```bash
echo 'PORT=5227' >> .env
```

<!--
```bash
rsync -avHP example.env .env
echo "API_TOKEN=$(openssl rand -hex 8)" >> .env
```
-->

### Build

```bash
# Install Go if needed (using Webi)
curl https://webinstall.dev/go@stable | bash
export PATH="$HOME/.local/opt/go/bin:$HOME/go/bin:${PATH}"

# Build the application
go build -o libreoffice-as-a-service .
```

### Run

```bash
./libreoffice-as-a-service
```

Or with custom port and bind address:

```bash
./libreoffice-as-a-service --port 5227 --bind 0.0.0.0
```

### Install

```bash
bash scripts/install.sh
```

## Demo Page

```bash
# LaaS => 5227
open http://127.0.0.1:5227/
```

## Demo cURL

```sh
mkdir -p ./fixtures/
curl -fsSL 'https://github.com/paperos-labs/libreoffice-as-a-service/raw/refs/heads/main/fixtures/Writing1.docx' \
    -o ./fixtures/Writing1.docx
```

```sh
export LAAS_BASE_URL="http://127.0.0.1:5227"
export LAAS_API_TOKEN="xxxx-xxxx-xxxx-xxxx"
```

### docx to pdf

```sh
curl -fS "${LAAS_BASE_URL}"'/api/convert/pdf?filename=Writing1.docx' \
    -H "Authorization: Bearer ${LAAS_API_TOKEN}" \
    -H 'Content-Type: application/octet-stream' \
    --data-binary @'fixtures/Writing1.docx' \
    -o ./Writing1.pdf
```

### docx to txt

```sh
curl -fS "${LAAS_BASE_URL}"'/api/convert/txt?filename=Writing1.docx' \
    -H "Authorization: Bearer ${LAAS_API_TOKEN}" \
    --data-binary @'fixtures/Writing1.docx' \
    -o ./Writing1.docx.txt
```

### pdf to txt

```sh
curl -fS "${LAAS_BASE_URL}"'/api/convert/txt?filename=Writing1.pdf' \
    -H "Authorization: Bearer ${LAAS_API_TOKEN}" \
    -H 'Content-Type: application/octet-stream' \
    --data-binary @'./Writing1.pdf' \
    -o ./Writing1.pdf.txt
```

**Important**: `-d` is NOT the same as `--data-binary` (the former may strip whitespace).

## Example Scripts

Example scripts are provided in multiple languages:

### Bash (example.sh)

```bash
bash example.sh
```

### JavaScript/Node.js (example.js)

```bash
node example.js
```

### Go (example/main.go)

```bash
cd example
go run .
```

All examples convert `fixtures/Writing1.docx` to `Writing1.pdf`. Set `LAAS_BASE_URL` and `LAAS_API_TOKEN` environment variables to customize.

## System Requirements for Linux

- Go 1.22+
- LibreOffice v6.4+
- poppler-utils (for pdftotext)

```bash
# Install Go
if [[ -z "$(command -v go)" ]]; then
    curl -fsSL https://webinstall.dev/go@stable | bash
    export PATH="$HOME/.local/opt/go/bin:$HOME/go/bin:${PATH}"
fi

# Install LibreOffice
if [[ -z "$(command -v libreoffice)" ]]; then
    sudo add-apt-repository -y ppa:libreoffice/ppa
    sudo add-apt-repository -y ppa:libreoffice/libreoffice-6-4
    sudo apt-get -y update
    sudo apt-get install -y libreoffice
fi

# Install poppler-utils (for pdftotext)
if [[ -z "$(command -v pdftotext)" ]]; then
    sudo apt-get -y update
    sudo apt-get install -y poppler-utils
fi
```

### How to test locally on macOS

- [Download LibreOffice for macOS](https://www.libreoffice.org/download/download/)

If you install LibreOffice to `~/Applications`, you can add `soffice` to your PATH, like so:

1. Temporarily add `soffice` to your `PATH`
   ```bash
   export PATH="/Applications/LibreOffice.app/Contents/MacOS:$PATH"
   ```
2. Install `pathman`
   ```bash
   curl https://webinstall.dev/pathman | bash
   export PATH="${HOME}/.local/bin:${PATH}"
   ```
3. Permanently add `soffice` to your `PATH`
   ```bash
   pathman add /Applications/LibreOffice.app/Contents/MacOS
   ```

Now Go will be able to find `soffice` and run it the same as on Linux.

**Note**: You can also use Webi to install Go:

```bash
curl -L https://webinstall.dev/go@stable | bash
```

### How to deploy to Digital Ocean

Create a `.env` with at least your Digital Ocean API Token and a DNS provider's token:

(as you can see, currently Cloudflare, DuckDNS, and Godaddy are supported)

```bash
DIGITALOCEAN_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
# the default tag will be 'delete-me' if you don't set one
DIGITALOCEAN_TAG=delete-me
# optional
DIGITALOCEAN_PROJECT=00000000-0000-4000-8000-000000000000

DNS_DEVELOPMENT_API=scripts/builder/00-duckdns-api.sh
DUCKDNS_API_TOKEN=00000000-0000-4000-8000-000000000000

#DNS_DEVELOPMENT_API=scripts/builder/00-godaddy-api.sh
#GODADDY_API_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
#GODADDY_API_SECRET=xxxxxxxxxxxxxxxxxxxxxx

DNS_PRODUCTION_API=scripts/builder/00-cloudflare-api.sh
CLOUDFLARE_API_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

Then you can run `script/provision.sh` like this:

```bash
SERVICE_NAME=libreoffice-as-a-service
GIT_REF_NAME='development'
DNS_RECORD_NAME=dev-laas
DNS_ZONE_NAME=example.net

bash scripts/provision.sh \
    "${SERVICE_NAME}" "${GIT_REF_NAME}" "${DNS_RECORD_NAME}" "${DNS_ZONE_NAME}"
```

Note: Everything other than the git branch `production` is considered to be `development` as far as
the deployment is concerned. You can add more branches and config at the bottom of
`scripts/provision.sh`.
