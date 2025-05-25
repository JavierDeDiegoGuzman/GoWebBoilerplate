#!/bin/bash

# Generar templates
templ generate

# Asegurar que el directorio existe
mkdir -p ./static/css

# Generar CSS y crear hash
npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/output.css
CSS_HASH=$(shasum -a 256 ./static/css/output.css | cut -d' ' -f1)
mv ./static/css/output.css "./static/css/output_${CSS_HASH}.css"

# Ejecutar la aplicación
go run cmd/main.go
