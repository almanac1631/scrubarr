#!/bin/bash
set -ex
openapi-generator generate -i ./openapi.yaml -g openapi -o .
oapi-codegen -config ./oapi-codegen-config.yaml ./openapi.json
openapi-generator generate -i ./openapi.yaml -g typescript-axios -o ../web/app/src/api/
