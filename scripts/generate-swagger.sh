#!/usr/bin/env bash
set -euo pipefail

echo "Installing swag CLI (if needed)..."
if ! command -v swag >/dev/null 2>&1; then
  go install github.com/swaggo/swag/cmd/swag@latest
fi

SWAG_DIRS=(
  serverbase
  shared-modules/audit
  shared-modules/booking
  shared-modules/calendar
  shared-modules/documents
  shared-modules/invoice
  shared-modules/invoice_number
  shared-modules/organization
  shared-modules/pdf
  shared-modules/saas-base
  shared-modules/settings
  shared-modules/static
)

for d in "${SWAG_DIRS[@]}"; do
  echo "=== Generating swagger for $d ==="
  if [ ! -d "$d" ]; then
    echo "SKIP: $d does not exist"
    continue
  fi
  # Skip directories without Go files
  if ! ls "$d"/*.go >/dev/null 2>&1; then
    echo "SKIP: no Go files in $d"
    continue
  fi

  # Each module must register under a unique swag instance name; the default
  # "swagger" name collides with the app-level instance once linked together.
  instance="$(basename "$d")"

  if [ -f "$d/module.go" ]; then
    (cd "$d" && swag init -g module.go --instanceName "$instance" -o docs)
  elif [ -f "$d/server.go" ]; then
    (cd "$d" && swag init -g server.go --instanceName "$instance" -o docs)
  else
    (cd "$d" && swag init --instanceName "$instance" -o docs)
  fi

  # swag prefixes output files/identifiers with the instance name; restore the
  # canonical docs.go/SwaggerInfo/docTemplate names so existing imports keep working.
  (cd "$d/docs" && \
    mv -f "${instance}_docs.go" docs.go 2>/dev/null || true; \
    mv -f "${instance}_swagger.json" swagger.json 2>/dev/null || true; \
    mv -f "${instance}_swagger.yaml" swagger.yaml 2>/dev/null || true; \
    sed -i '' "s/SwaggerInfo${instance}/SwaggerInfo/g; s/docTemplate${instance}/docTemplate/g" docs.go 2>/dev/null || true)
done

echo "ALL SWAGGER DOCS GENERATED"
