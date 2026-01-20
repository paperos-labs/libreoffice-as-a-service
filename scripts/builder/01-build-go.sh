#!/bin/bash
set -e
set -u

if [[ -z "$(command -v go)" ]]; then
    curl https://webinstall.dev/go@stable | bash
fi
export PATH="$HOME/.local/opt/go/bin:$HOME/go/bin:${PATH}"

go build -o libreoffice-as-a-service .
