#!/bin/bash
# Download the latest Huntress OpenAPI/Swagger spec for model/schema compliance checking
set -euo pipefail

OUT="swagger_doc.json"
URL="https://api.huntress.io/swagger_doc.json"

curl -sSfL "$URL" -o "$OUT"
echo "Downloaded $URL to $OUT"
