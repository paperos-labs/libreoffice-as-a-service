#!/bin/bash
set -e
set -u

BASE_URL="http://127.0.0.1:5227"

# Download sample file if it doesn't exist
if [[ ! -f fixtures/Writing1.docx ]]; then
    mkdir -p fixtures
    curl -fsSL 'https://github.com/paperos-labs/libreoffice-as-a-service/raw/refs/heads/main/fixtures/Writing1.docx' \
        -o fixtures/Writing1.docx
fi

# Convert docx to pdf
curl -fS "${BASE_URL}/api/convert/pdf?filename=Writing1.docx" \
    -H 'Content-Type: application/octet-stream' \
    --data-binary @'fixtures/Writing1.docx' \
    -o Writing1.pdf

echo "Converted Writing1.docx to Writing1.pdf"
ls -lh Writing1.pdf
