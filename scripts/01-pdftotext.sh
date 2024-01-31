#!/bin/bash
set -e
set -u

if [[ -z "$(command -v pdftotext)" ]]; then
    export DEBIAN_FRONTEND=noninteractive

    sudo apt-get -y update
    sleep 2
    sudo rm -f /var/lib/apt/lists/lock
    sudo rm -f /var/lib/dpkg/lock
    sudo killall apt-get || true

    sudo apt-get install -y poppler-utils
    sleep 2
    sudo rm -f /var/lib/apt/lists/lock
    sudo rm -f /var/lib/dpkg/lock
    sudo killall apt-get || true
fi
