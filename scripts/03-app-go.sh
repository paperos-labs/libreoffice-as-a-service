#!/bin/bash
set -e
set -u

#my_domain="${1}"
#my_port="${2}"
my_servicename="${3}"

if [[ -z "$(command -v watchexec)" ]]; then
    webi watchexec
fi

if [[ -z "$(command -v go)" ]]; then
    webi go@stable
fi
export PATH="$HOME/.local/opt/go/bin:$HOME/go/bin:${PATH}"

# Build the Go binary
go build -o libreoffice-as-a-service .

if [[ "development" == "${NODE_ENV:-}" ]]; then
    # stop watchexec first in development environments
    if ! sudo systemctl status "${my_servicename}" | grep 'could not'; then
        sudo systemctl stop "${my_servicename}" || true
    fi
    sudo env PATH="${PATH}" \
        serviceman add --name "${my_servicename}" --system \
        --username "$(whoami)" --path "${PATH}" -- \
        watchexec -r -e go -- -- \
        ./libreoffice-as-a-service # -- --port "${my_port}"
else
    sudo env PATH="${PATH}" \
        serviceman add --name "${my_servicename}" --system \
        --username "$(whoami)" --path "${PATH}" -- \
        ./libreoffice-as-a-service # -- --port "${my_port}"
fi
